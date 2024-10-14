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

	"github.com/buildkite/go-buildkite/v3/internal/bkmultipart"
)

// CreatePackageInput specifies the input parameters for the Create method.
// All fields are required, but if PackageFile is an [os.File], Filename can be omitted.
type CreatePackageInput struct {
	Package  io.Reader // The package to upload. This can be an [os.File], or any other [io.Reader].
	Filename string    // The name of the file to upload. If PackageFile is an [os.File], this can be omitted, and the file's name will be used.
}

// Create creates a package in a registry for an organization
func (ps *PackagesService) Create(ctx context.Context, organizationSlug, registrySlug string, cpi CreatePackageInput) (Package, *Response, error) {
	filename := cpi.Filename
	if f, ok := cpi.Package.(*os.File); ok && filename == "" {
		filename = f.Name()
	}

	packageTempFile, err := normalizeToFile(cpi.Package, filename)
	if err != nil {
		return Package{}, nil, fmt.Errorf("writing package to tempfile: %v", err)
	}
	defer func() {
		packageTempFile.Close()
		os.Remove(packageTempFile.Name())
	}()

	s := bkmultipart.NewStreamer()
	s.WriteFile(fileFormKey, packageTempFile, filename)

	url := fmt.Sprintf("v2/packages/organizations/%s/registries/%s/packages", organizationSlug, registrySlug)
	req, err := ps.client.NewRequest(ctx, "POST", url, s.Reader())
	if err != nil {
		return Package{}, nil, fmt.Errorf("creating POST package request: %v", err)
	}

	req.Header.Set("Content-Type", s.ContentType)
	req.Header.Set("Content-Length", fmt.Sprintf("%d", s.Len()))

	var p Package
	resp, err := ps.client.Do(req, &p)
	if err != nil {
		return Package{}, resp, fmt.Errorf("executing POST package request: %v", err)
	}

	return p, resp, err
}

// normalizeToFile takes and io.Reader (which might itself already be a file, but could be a stream or other source) and
// writes it to a temporary file, returning the file handle.
// The file is written to a temporary directory, and then renamed to the desired filename.
// We do this normalization to ensure that we can accurately calculate the Content-Length of the request body, which is
// required by S3. We write to disk (instead of buffering in memory) to avoid memory exhaustion for large files.
func normalizeToFile(r io.Reader, filename string) (*os.File, error) {
	basename := filepath.Base(filename)
	f, err := os.CreateTemp("", basename)
	if err != nil {
		return nil, fmt.Errorf("creating temporary file: %v", err)
	}

	_, err = io.Copy(f, r)
	if err != nil {
		return nil, fmt.Errorf("writing to temporary file: %v", err)
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		return nil, fmt.Errorf("seeking to beginning of temporary file: %v", err)
	}

	return f, nil
}

type PackagePresignedUpload struct {
	OrganizationSlug string `json:"-"`
	RegistrySlug     string `json:"-"`

	URI  string `json:"uri"`
	Form struct {
		FileInput string            `json:"file_input"`
		Method    string            `json:"method"`
		URL       string            `json:"url"`
		Data      map[string]string `json:"data"`
	} `json:"form"`
}

func (ps *PackagesService) RequestPresignedUpload(ctx context.Context, organizationSlug, registrySlug string) (PackagePresignedUpload, *Response, error) {
	url := fmt.Sprintf("v2/packages/organizations/%s/registries/%s/packages/upload", organizationSlug, registrySlug)
	req, err := ps.client.NewRequest(ctx, "POST", url, nil)
	if err != nil {
		return PackagePresignedUpload{}, nil, fmt.Errorf("creating POST presigned upload request: %v", err)
	}

	var p PackagePresignedUpload
	resp, err := ps.client.Do(req, &p)
	if err != nil {
		fmt.Println(string(err.(*ErrorResponse).RawBody))
		return PackagePresignedUpload{}, resp, fmt.Errorf("executing POST presigned upload request: %v", err)
	}

	p.OrganizationSlug = organizationSlug
	p.RegistrySlug = registrySlug

	return p, resp, err
}

func (ppu PackagePresignedUpload) Perform(ctx context.Context, ps *PackagesService, file *os.File) (string, error) {
	if _, ok := ppu.Form.Data["key"]; !ok {
		return "", fmt.Errorf("missing 'key' in presigned upload form data")
	}

	baseFilePath := filepath.Base(file.Name())

	s := bkmultipart.NewStreamer()
	s.WriteFields(ppu.Form.Data)
	s.WriteFile(ppu.Form.FileInput, file, baseFilePath)

	req, err := http.NewRequestWithContext(ctx, ppu.Form.Method, ppu.Form.URL, s.Reader())
	if err != nil {
		return "", fmt.Errorf("creating %s request: %v", ppu.Form.Method, err)
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

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("executing %s request: %v", ppu.Form.Method, err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("S3 rejected upload with unexpected status code %d. Error reading response body: %v", resp.StatusCode, err)
		}

		return "", fmt.Errorf("S3 rejected upload with unexpected status code %d. Response body %s", resp.StatusCode, string(body))
	}

	uploadPath, err := url.JoinPath(ppu.Form.URL, strings.ReplaceAll(ppu.Form.Data["key"], "${filename}", baseFilePath))
	if err != nil {
		return "", fmt.Errorf("joining URL path: %v", err)
	}

	return uploadPath, nil
}

func (ppu PackagePresignedUpload) Finalize(ctx context.Context, ps *PackagesService, s3URL string) (Package, *Response, error) {
	s := bkmultipart.NewStreamer()
	s.WriteField("package_url", s3URL)

	url := fmt.Sprintf("v2/packages/organizations/%s/registries/%s/packages", ppu.OrganizationSlug, ppu.RegistrySlug)
	req, err := ps.client.NewRequest(ctx, "POST", url, s)
	if err != nil {
		return Package{}, nil, fmt.Errorf("creating POST package request: %v", err)
	}

	req.Header.Set("Content-Type", s.ContentType)
	req.ContentLength = s.Len()

	var p Package
	resp, err := ps.client.Do(req, &p)
	if err != nil {
		return Package{}, resp, fmt.Errorf("executing POST package request: %v", err)
	}

	return p, resp, err
}
