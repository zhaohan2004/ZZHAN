package service

import (
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"strings"
	"testing"
)

type coverTestStorage func(context.Context, string, io.Reader, int64, string) (string, error)

func (f coverTestStorage) Upload(ctx context.Context, key string, r io.Reader, size int64, mime string) (string, error) {
	return f(ctx, key, r, size, mime)
}

type coverTestTransport func(*http.Request) (*http.Response, error)

func (f coverTestTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestUploadLegacyCoverRequiresReadableImage(t *testing.T) {
	const encoded = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/x8AAwMCAO+jBXsAAAAASUVORK5CYII="
	const link = "https://images.example.test/cover.png"
	for _, status := range []int{http.StatusPartialContent, http.StatusForbidden} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			store := coverTestStorage(func(_ context.Context, key string, r io.Reader, size int64, mime string) (string, error) {
				data, err := io.ReadAll(r)
				if err != nil {
					t.Fatal(err)
				}
				if !strings.HasPrefix(key, "images/covers/15-") || !strings.HasSuffix(key, ".png") ||
					mime != "image/png" || size != int64(len(data)) || base64.StdEncoding.EncodeToString(data) != encoded {
					t.Fatalf("unexpected uploaded cover: %s %s (%d bytes)", key, mime, size)
				}
				return link, nil
			})
			client := &http.Client{Transport: coverTestTransport(func(req *http.Request) (*http.Response, error) {
				if req.Method != http.MethodGet || req.URL.String() != link || req.Header.Get("Range") != "bytes=0-0" {
					t.Fatalf("unexpected cover verification request: %v", req)
				}
				return &http.Response{
					StatusCode: status,
					Header:     http.Header{"Content-Type": {"image/png"}},
					Body:       io.NopCloser(strings.NewReader("")),
				}, nil
			})}

			got, err := uploadLegacyCover(context.Background(), client, store, 15, "data:image/png;base64,"+encoded)
			if status == http.StatusForbidden {
				if err == nil || got != "" {
					t.Fatalf("unreadable image returned link %q and error %v", got, err)
				}
			} else if err != nil || got != link {
				t.Fatalf("got %q, %v; want %q", got, err, link)
			}
		})
	}
}
