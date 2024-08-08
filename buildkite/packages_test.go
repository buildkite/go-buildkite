package buildkite

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

var (
	pkg = Package{
		ID:       uuid.NewString(),
		URL:      "https://example.com/my-package",
		WebURL:   "https://buildkite.com/my-org/my-registry/my-package",
		Registry: registry,
		Organization: Organization{
			ID:   String("my-org"),
			Slug: String("my-org"),
			Name: String("My Org"),
		},
	}
)

func TestGetPackage(t *testing.T) {
	t.Parallel()

	mux, client, teardown := newMockServerAndClient(t)
	t.Cleanup(teardown)

	want := pkg
	endpoint := fmt.Sprintf("/v2/packages/organizations/my-org/registries/my-registry/packages/%s", pkg.ID)
	mux.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		err := json.NewEncoder(w).Encode(pkg)
		if err != nil {
			t.Fatalf("marshalling package to json: %v", err)
		}
	})

	p, _, err := client.PackagesService.Get("my-org", "my-registry", pkg.ID)
	if err != nil {
		t.Fatalf("Packages.Get returned error: %v", err)
	}

	if diff := cmp.Diff(p, want); diff != "" {
		t.Fatalf("client.PackagesService.Get(%q, %q, %q) diff: (-got +want)\n%s", "test-org", "my-cool-registry", pkg.ID, diff)
	}
}

func TestCreatePackage(t *testing.T) {
	t.Parallel()

	testPackage, err := os.CreateTemp("", "test-package")
	if err != nil {
		t.Fatalf("creating temporary package file: %v", err)
	}
	t.Cleanup(func() { os.Remove(testPackage.Name()) })

	packageContents := "this is totally a valid package! look, i'm a rubygem!"
	_, err = testPackage.WriteString(packageContents)
	if err != nil {
		t.Fatalf("writing to temporary package file: %v", err)
	}

	if _, err := testPackage.Seek(0, io.SeekStart); err != nil {
		t.Fatalf("seeking to start of temporary package file: %v", err)
	}

	cases := []struct {
		name         string
		in           CreatePackageInput
		wantContents string
		wantFileName string
	}{
		{
			name:         "file",
			in:           CreatePackageInput{Package: testPackage},
			wantContents: packageContents,
			wantFileName: testPackage.Name(),
		},
		{
			name: "io.Reader with filename",
			in: CreatePackageInput{
				Package:  bytes.NewBufferString(packageContents),
				Filename: "cool-package.gem",
			},
			wantContents: packageContents,
			wantFileName: "cool-package.gem",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mux, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)

			mux.HandleFunc("/v2/packages/organizations/my-org/registries/my-registry/packages", func(w http.ResponseWriter, r *http.Request) {
				defer r.Body.Close()

				testMethod(t, r, "POST")

				if !strings.HasPrefix(r.Header.Get("Content-Type"), "multipart/form-data") {
					t.Fatalf("unexpected Content-Type: %q", r.Header.Get("Content-Type"))
				}

				fi, header, err := r.FormFile(fileFormKey)
				if err != nil {
					t.Fatalf("getting file from request: %v", err)
				}
				defer fi.Close()

				// RFC 7578 says that the any path information should be stripped from the file name, which is what
				// r.FormFile does - see https://github.com/golang/go/blob/d9f9746/src/mime/multipart/multipart.go#L99-L100
				if header.Filename != filepath.Base(tc.wantFileName) {
					t.Fatalf("file name mismatch: got %q, want %q", header.Filename, tc.wantFileName)
				}

				fileContents, err := io.ReadAll(fi)
				if err != nil {
					t.Fatalf("reading file contents: %v", err)
				}

				if string(fileContents) != tc.wantContents {
					t.Fatalf("file contents mismatch: got %q, want %q", string(fileContents), tc.wantContents)
				}

				err = json.NewEncoder(w).Encode(pkg)
				if err != nil {
					t.Fatalf("encoding package to json: %v", err)
				}
			})

			p, _, err := client.PackagesService.Create("my-org", "my-registry", tc.in)
			if err != nil {
				t.Fatalf("Packages.Create returned error: %v", err)
			}

			if diff := cmp.Diff(p, pkg); diff != "" {
				t.Fatalf("client.PackagesService.Create(%q, %q, %v) diff: (-got +want)\n%s", "test-org", "my-cool-registry", tc.in, diff)
			}
		})
	}
}
