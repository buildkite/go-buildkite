package buildkite

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
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

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, err := w.CreateFormFile(fileFormKey, filename)
	if err != nil {
		return Package{}, nil, fmt.Errorf("creating multipart form file: %v", err)
	}

	_, err = io.Copy(fw, cpi.Package)
	if err != nil {
		return Package{}, nil, fmt.Errorf("copying data into multipart payload: %v", err)
	}
	w.Close()

	url := fmt.Sprintf("v2/packages/organizations/%s/registries/%s/packages", organizationSlug, registrySlug)
	req, err := ps.client.NewRequest(ctx, "POST", url, &b)
	if err != nil {
		return Package{}, nil, fmt.Errorf("creating POST package request: %v", err)
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	var p Package
	resp, err := ps.client.Do(req, &p)
	if err != nil {
		return Package{}, resp, fmt.Errorf("executing POST package request: %v", err)
	}

	return p, resp, err
}
