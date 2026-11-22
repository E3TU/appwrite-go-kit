// Package storage wraps the self-hosted Appwrite storage surface used by Gin
// services.
//
// The service is a thin layer over the SDK storage API. The only value it adds
// is building browser-facing URLs from the public endpoint.
package storage

import (
	"errors"

	"github.com/appwrite/sdk-for-go/client"
	"github.com/appwrite/sdk-for-go/storage"
)

// Service exposes the storage operations the library supports.
type Service struct {
	sdk             *storage.Storage
	defaultBucketID string
	publicEndpoint  string
	projectID       string
}

// New builds a storage Service around an SDK client.
//
// publicEndpoint is the base for public URLs and is stored unchanged; it is
// never rewritten to the internal endpoint. projectID is needed because public
// file URLs are scoped per project.
func New(c client.Client, defaultBucketID, publicEndpoint, projectID string) *Service {
	return &Service{
		sdk:             storage.New(c),
		defaultBucketID: defaultBucketID,
		publicEndpoint:  publicEndpoint,
		projectID:       projectID,
	}
}

// Storage returns the underlying SDK storage service.
func (s *Service) Storage() *storage.Storage {
	return s.sdk
}

// DefaultBucketID returns the configured default bucket.
func (s *Service) DefaultBucketID() string {
	return s.defaultBucketID
}

// ErrNoPublicEndpoint reports URL building without a configured public
// endpoint.
var ErrNoPublicEndpoint = errors.New("storage: public endpoint not configured")
