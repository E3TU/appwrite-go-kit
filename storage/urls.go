package storage

import (
	"fmt"
	"net/url"
	"strings"
)

// Public file URLs are rooted at the configured public endpoint and include
// the project id. The internal endpoint must never appear here.

// PublicFileURL returns the public download URL for fileID in bucketID.
//
// An empty bucketID falls back to the configured default bucket.
func (s *Service) PublicFileURL(bucketID, fileID string) (string, error) {
	return s.publicURL(bucketID, fileID, "view")
}

// PublicPreviewURL returns the public preview URL for fileID in bucketID.
//
// An empty bucketID falls back to the configured default bucket.
func (s *Service) PublicPreviewURL(bucketID, fileID string) (string, error) {
	return s.publicURL(bucketID, fileID, "preview")
}

func (s *Service) publicURL(bucketID, fileID, mode string) (string, error) {
	if s.publicEndpoint == "" {
		return "", ErrNoPublicEndpoint
	}
	if bucketID == "" {
		bucketID = s.defaultBucketID
	}
	if bucketID == "" {
		return "", fmt.Errorf("%w: no bucket given and no default configured", ErrNoPublicEndpoint)
	}
	if fileID == "" {
		return "", fmt.Errorf("%w: file id is required", ErrNoPublicEndpoint)
	}
	if s.projectID == "" {
		return "", fmt.Errorf("%w: project id is required", ErrNoPublicEndpoint)
	}

	base := strings.TrimRight(s.publicEndpoint, "/")
	u := base + "/storage/buckets/" + url.PathEscape(bucketID) +
		"/files/" + url.PathEscape(fileID) + "/" + mode

	q := url.Values{}
	q.Set("project", s.projectID)

	return u + "?" + q.Encode(), nil
}
