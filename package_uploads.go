package buildkite

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/buildkite/go-buildkite/v4/internal/bkmultipart"
)

// CreatePackageInput specifies the input parameters for the Create method.
// All fields are required, but if PackageFile is an [os.File], Filename can be omitted.
type CreatePackageInput struct {
	Package  io.Reader // The package to upload. This can be an [os.File], or any other [io.Reader].
	Filename string    // The name of the file to upload. If PackageFile is an [os.File], this can be omitted, and the file's name will be used.
}

// Create creates a package in a registry for an organization
func (ps *PackagesService) Create(ctx context.Context, organizationSlug, registrySlug string, cpi CreatePackageInput) (Package, *Response, error) {
	var file *os.File
	switch f := cpi.Package.(type) {
	case *os.File:
		file = f

	default:
		var err error
		file, err = readIntoTempFile(cpi.Package, cpi.Filename)
		if err != nil {
			return Package{}, nil, fmt.Errorf("writing package to tempfile: %w", err)
		}

		defer func() {
			_ = file.Close()
			_ = os.Remove(file.Name())
		}()
	}

	ppu, _, err := ps.RequestPresignedUpload(ctx, organizationSlug, registrySlug)
	if err != nil {
		return Package{}, nil, fmt.Errorf("requesting presigned upload: %w", err)
	}

	s3URL, err := ppu.Perform(ctx, file)
	if err != nil {
		return Package{}, nil, fmt.Errorf("performing presigned upload: %w", err)
	}

	p, resp, err := ppu.Finalize(ctx, s3URL)
	if err != nil {
		return Package{}, nil, fmt.Errorf("finalizing package: %w", err)
	}

	return p, resp, nil
}

// readIntoTempFile takes an io.Reader and writes it to a temporary file, returning the file handle.
// The file is written to a temporary directory, and then renamed to the desired filename.
// We do this normalization to ensure that we can accurately calculate the Content-Length of the request body, which is
// required by S3. We write to disk (instead of buffering in memory) to avoid memory exhaustion for large files.
func readIntoTempFile(r io.Reader, filename string) (*os.File, error) {
	basename := filepath.Base(filename)
	f, err := os.CreateTemp("", basename)
	if err != nil {
		return nil, fmt.Errorf("creating temporary file: %w", err)
	}

	_, err = io.Copy(f, r)
	if err != nil {
		return nil, fmt.Errorf("writing to temporary file: %w", err)
	}

	err = f.Close()
	if err != nil {
		return nil, fmt.Errorf("closing temporary file: %w", err)
	}

	// Rename the temporary file to the desired filename, which is important for Buildkite Package indexing
	newFileName := filepath.Join(filepath.Dir(f.Name()), basename)
	err = os.Rename(f.Name(), newFileName)
	if err != nil {
		return nil, fmt.Errorf("renaming temporary file: %w", err)
	}

	f, err = os.Open(newFileName)
	if err != nil {
		return nil, fmt.Errorf("opening renamed file: %w", err)
	}

	return f, nil
}

// PackagePresignedUpload represents a presigned upload URL for a Buildkite package, returned by the Buildkite API
type PackagePresignedUpload struct {
	bkClient *Client

	OrganizationSlug string `json:"-"`
	RegistrySlug     string `json:"-"`

	URI  string                     `json:"uri"`
	Form PackagePresignedUploadForm `json:"form"`
}

type PackagePresignedUploadForm struct {
	FileInput string            `json:"file_input"`
	Method    string            `json:"method"`
	URL       string            `json:"url"`
	Data      map[string]string `json:"data"`
}

// RequestPresignedUpload requests a presigned upload URL for a Buildkite package from the buildkite API
func (ps *PackagesService) RequestPresignedUpload(ctx context.Context, organizationSlug, registrySlug string) (*PackagePresignedUpload, *Response, error) {
	url := fmt.Sprintf("v2/packages/organizations/%s/registries/%s/packages/upload", organizationSlug, registrySlug)
	req, err := ps.client.NewRequest(ctx, "POST", url, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("creating POST presigned upload request: %w", err)
	}

	var p *PackagePresignedUpload
	resp, err := ps.client.Do(req, &p)
	if err != nil {
		return nil, resp, fmt.Errorf("executing POST presigned upload request: %w", err)
	}

	p.bkClient = ps.client
	p.OrganizationSlug = organizationSlug
	p.RegistrySlug = registrySlug

	return p, resp, err
}

// Perform performs uploads the package file referred to by `file` to the presigned upload URL.
// It does not create the package in the registry, only uploads the file to the package host. The returned string is the URL of the
// uploaded file in S3, which can then be passed to [Finalize] to create the package in the registry.
func (ppu *PackagePresignedUpload) Perform(ctx context.Context, file *os.File) (string, error) {
	if _, ok := ppu.Form.Data["key"]; !ok {
		return "", fmt.Errorf("missing 'key' in presigned upload form data")
	}

	baseFilePath := filepath.Base(file.Name())

	s := bkmultipart.NewStreamer()
	err := s.WriteFields(ppu.Form.Data)
	if err != nil {
		return "", fmt.Errorf("writing form fields: %w", err)
	}

	err = s.WriteFile(ppu.Form.FileInput, file, baseFilePath)
	if err != nil {
		return "", fmt.Errorf("writing form file: %w", err)
	}

	// note NOT using client.NewRequest here, as it'll add buildkite-specific stuff that we don't want
	req, err := http.NewRequestWithContext(ctx, ppu.Form.Method, ppu.Form.URL, s.Reader())
	if err != nil {
		return "", fmt.Errorf("creating %s request: %w", ppu.Form.Method, err)
	}

	req.Header.Set("Content-Type", s.ContentType)

	// Don't set the Content-Length header here, you fool, you absolute buffoon
	// When passed an io.Reader, http.NewRequestWithContext will not set the Content-Length header, and will instead
	// stream the request body. This _would_ be exactly what we want, except that S3 uploads infuriatingly require a
	// Content-Length header. So we have to calculate the length of the request body ourselves and set it manually on the
	// request. Adding:
	// 	req.Header.Set("Content-Length", fmt.Sprintf("%d", s.Len()))
	// is not sufficient, as the Content-Length header is stripped by the http client when the request body is an io.Reader.
	req.ContentLength = s.Len()

	resp, err := ppu.bkClient.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("executing %s request: %w", ppu.Form.Method, err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("S3 rejected upload with unexpected status code %d. Error reading response body: %w", resp.StatusCode, err)
		}

		return "", fmt.Errorf("S3 rejected upload with unexpected status code %d. Response body %s", resp.StatusCode, string(body))
	}

	uploadPath, err := url.JoinPath(ppu.Form.URL, strings.ReplaceAll(ppu.Form.Data["key"], "${filename}", baseFilePath))
	if err != nil {
		return "", fmt.Errorf("joining URL path: %w", err)
	}

	return uploadPath, nil
}

// Finalize creates a package in the registry for the organization, using the S3 URL of the uploaded package file.
func (ppu *PackagePresignedUpload) Finalize(ctx context.Context, s3URL string) (Package, *Response, error) {
	s := bkmultipart.NewStreamer()
	err := s.WriteField("package_url", s3URL)
	if err != nil {
		return Package{}, nil, fmt.Errorf("writing package_url field: %w", err)
	}

	url := fmt.Sprintf("v2/packages/organizations/%s/registries/%s/packages", ppu.OrganizationSlug, ppu.RegistrySlug)
	req, err := ppu.bkClient.NewRequest(ctx, "POST", url, s.Reader())
	if err != nil {
		return Package{}, nil, fmt.Errorf("creating POST package request: %w", err)
	}

	req.Header.Set("Content-Type", s.ContentType)
	req.ContentLength = s.Len()

	var p Package
	resp, err := ppu.bkClient.Do(req, &p)
	if err != nil {
		return Package{}, resp, fmt.Errorf("executing POST package request: %w", err)
	}

	return p, resp, err
}
