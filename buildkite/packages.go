package buildkite

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/buildkite/go-buildkite/v3/internal/bkmultipart"
)

const fileFormKey = "file"

// PackagesService handles communication with packages Buildkite API endpoints
type PackagesService struct {
	client *Client
}

// Package represents a package which has been stored in a registry
type Package struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	URL          string          `json:"url"`
	WebURL       string          `json:"web_url"`
	Organization Organization    `json:"organization"`
	Registry     PackageRegistry `json:"registry"`
}

func (ps *PackagesService) Get(ctx context.Context, organizationSlug, registrySlug, packageID string) (Package, *Response, error) {
	url := fmt.Sprintf("v2/packages/organizations/%s/registries/%s/packages/%s", organizationSlug, registrySlug, packageID)
	req, err := ps.client.NewRequest(ctx, "GET", url, nil)
	if err != nil {
		return Package{}, nil, fmt.Errorf("creating GET package request: %v", err)
	}

	var p Package
	resp, err := ps.client.Do(req, &p)
	if err != nil {
		return Package{}, resp, fmt.Errorf("executing GET package request: %v", err)
	}

	return p, resp, err
}

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
	defer os.Remove(packageTempFile.Name())
	defer packageTempFile.Close()

	s := bkmultipart.NewStreamer()
	s.WriteFile(fileFormKey, packageTempFile, filename)

	url := fmt.Sprintf("v2/packages/organizations/%s/registries/%s/packages", organizationSlug, registrySlug)
	req, err := ps.client.NewRequest(ctx, "POST", url, s)
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

	// Rename the temporary file to the desired filename, which is important for Buildkite Package indexing
	tempFileDir := filepath.Dir(f.Name())
	err = os.Rename(f.Name(), filepath.Join(tempFileDir, basename))
	if err != nil {
		return nil, fmt.Errorf("renaming temporary file: %v", err)
	}

	return f, nil
}
