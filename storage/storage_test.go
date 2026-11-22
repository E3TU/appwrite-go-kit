package storage

import (
	"errors"
	"testing"
)

func newTestService(public, bucket, project string) *Service {
	return &Service{
		defaultBucketID: bucket,
		publicEndpoint:  public,
		projectID:       project,
	}
}

func TestPublicFileURLUsesPublicEndpoint(t *testing.T) {
	s := newTestService("https://appwrite.example.com/v1/", "files", "proj")

	got, err := s.PublicFileURL("files", "abc")
	if err != nil {
		t.Fatalf("PublicFileURL() error = %v", err)
	}
	want := "https://appwrite.example.com/v1/storage/buckets/files/files/abc/view?project=proj"
	if got != want {
		t.Fatalf("url = %q, want %q", got, want)
	}
}

func TestPublicPreviewURLFallsBackToDefaultBucket(t *testing.T) {
	s := newTestService("https://appwrite.example.com/v1", "files", "proj")

	got, err := s.PublicPreviewURL("", "abc")
	if err != nil {
		t.Fatalf("PublicPreviewURL() error = %v", err)
	}
	want := "https://appwrite.example.com/v1/storage/buckets/files/files/abc/preview?project=proj"
	if got != want {
		t.Fatalf("url = %q, want %q", got, want)
	}
}

func TestPublicURLRequiresEndpoint(t *testing.T) {
	s := newTestService("", "files", "proj")

	if _, err := s.PublicFileURL("files", "abc"); !errors.Is(err, ErrNoPublicEndpoint) {
		t.Fatalf("error = %v, want ErrNoPublicEndpoint", err)
	}
}

func TestPublicURLRequiresBucket(t *testing.T) {
	s := newTestService("https://example.com/v1", "", "proj")

	if _, err := s.PublicFileURL("", "abc"); !errors.Is(err, ErrNoPublicEndpoint) {
		t.Fatalf("error = %v, want ErrNoPublicEndpoint", err)
	}
}

func TestPublicURLRequiresFileAndProject(t *testing.T) {
	s := newTestService("https://example.com/v1", "files", "proj")
	if _, err := s.PublicFileURL("files", ""); err == nil {
		t.Fatal("expected error for empty file id")
	}

	s = newTestService("https://example.com/v1", "files", "")
	if _, err := s.PublicFileURL("files", "abc"); err == nil {
		t.Fatal("expected error for empty project id")
	}
}
