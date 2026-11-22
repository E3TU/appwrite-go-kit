// Package functions wraps the self-hosted Appwrite Functions surface used by
// Gin services.
//
// The service is a thin layer over the SDK functions API.
package functions

import (
	"github.com/appwrite/sdk-for-go/v7/client"
	"github.com/appwrite/sdk-for-go/v7/functions"
)

// Service exposes the functions operations the library supports.
type Service struct {
	sdk *functions.Functions
}

// New builds a functions Service around an SDK client.
func New(c client.Client) *Service {
	return &Service{sdk: functions.New(c)}
}

// Functions returns the underlying SDK functions service.
func (s *Service) Functions() *functions.Functions {
	return s.sdk
}
