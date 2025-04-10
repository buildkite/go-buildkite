package buildkite

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCreatePackage(t *testing.T) {
	t.Parallel()

	testPackage, err := os.CreateTemp("", "test-package")
	if err != nil {
		t.Fatalf("creating temporary package file: %v", err)
	}
	t.Cleanup(func() { _ = os.Remove(testPackage.Name()) })

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

			server, client, teardown := newMockServerAndClient(t)
			t.Cleanup(teardown)

			s3Endpoint := "/s3"
			s3Path := "/fake/path"

			postData := map[string]string{
				"key":              s3Path + "/${filename}",
				"acl":              "private",
				"policy":           "bWFkZSB5b3UgbG9vayE=",
				"x-amz-credential": "AKIAS000000000000000/20241007/ap-southeast-2/s3/aws4_request",
				"x-amz-algorithm":  "AWS4-HMAC-SHA256",
				"x-amz-date":       "20241007T031838Z",
				"x-amz-signature":  "f6d24942026ffe7ec32b5f57beb46a2679b7a74a87673e1614b92c15ee2661f2",
			}

			// Signed Upload Request
			server.HandleFunc("/v2/packages/organizations/my-org/registries/my-registry/packages/upload", func(w http.ResponseWriter, r *http.Request) {
				defer func() { _ = r.Body.Close() }()

				testMethod(t, r, "POST")

				ppu := PackagePresignedUpload{
					URI: "s3://fake-s3-bucket/fake-s3-path", // URI is unused by go-buildkite, but here for completeness
					Form: PackagePresignedUploadForm{
						FileInput: "file",
						Method:    "POST",
						URL:       "http://" + r.Host + s3Endpoint,
						Data:      postData,
					},
				}

				w.Header().Set("Content-Type", "application/json")
				err := json.NewEncoder(w).Encode(ppu)
				if err != nil {
					t.Fatalf("encoding presigned upload to json: %v", err)
				}
			})

			// "S3" Upload
			server.HandleFunc(s3Endpoint, func(w http.ResponseWriter, r *http.Request) {
				defer func() { _ = r.Body.Close() }()

				testMethod(t, r, "POST")

				if r.Header.Get("Content-Length") == "" {
					t.Fatalf("missing Content-Length header - S3 requires it")
				}

				ct := r.Header.Get("Content-Type")
				mt, _, err := mime.ParseMediaType(ct)
				if err != nil {
					t.Fatalf("parsing Content-Type: %v", err)
				}

				if got, want := mt, "multipart/form-data"; got != want {
					t.Fatalf("unexpected media type: got %q, want %q", got, want)
				}

				fi, header, err := r.FormFile(fileFormKey)
				if err != nil {
					t.Fatalf("getting file from request: %v", err)
				}
				defer func() { _ = fi.Close() }()

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
			})

			// Create Package / Presigned upload finalization
			server.HandleFunc("/v2/packages/organizations/my-org/registries/my-registry/packages", func(w http.ResponseWriter, r *http.Request) {
				defer func() { _ = r.Body.Close() }()

				testMethod(t, r, "POST")

				err := r.ParseMultipartForm(2 << 10)
				if err != nil {
					t.Fatalf("parsing multipart form: %v", err)
				}

				wantPath, err := url.JoinPath(s3Endpoint, s3Path, filepath.Base(tc.wantFileName))
				if err != nil {
					t.Fatalf("joining URL path: %v", err)
				}

				wantURL := "http://" + r.Host + wantPath
				if got, want := r.Form["package_url"][0], wantURL; got != want {
					t.Fatalf("unexpected package URL: got %q, want %q", got, want)
				}

				err = json.NewEncoder(w).Encode(pkg)
				if err != nil {
					t.Fatalf("encoding package to json: %v", err)
				}
			})

			p, _, err := client.PackagesService.Create(context.Background(), "my-org", "my-registry", tc.in)
			if err != nil {
				t.Fatalf("Packages.Create returned error: %v", err)
			}

			wantHTTPCalls := []httpCall{
				{Method: "POST", Path: "/v2/packages/organizations/my-org/registries/my-registry/packages/upload"},
				{Method: "POST", Path: "/s3"},
				{Method: "POST", Path: "/v2/packages/organizations/my-org/registries/my-registry/packages"},
			}

			if diff := cmp.Diff(wantHTTPCalls, server.calls); diff != "" {
				t.Fatalf("unexpected HTTP calls (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(p, pkg); diff != "" {
				t.Fatalf("client.PackagesService.Create(%q, %q, %v) diff: (-got +want)\n%s", "test-org", "my-cool-registry", tc.in, diff)
			}

			// If we get passed in a file, we really don't want to delete it in the upload process. If this stat fails,
			// the file was deleted.
			if p, ok := tc.in.Package.(*os.File); ok {
				_, err := os.Stat(p.Name())
				if err != nil {
					t.Fatalf("expected stat file to have nil error, got: %v", err)
				}
			}
		})
	}
}
