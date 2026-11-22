// Package auth wraps the self-hosted Appwrite auth surface used by Gin
// services: session extraction from requests and a RequireUser middleware.
//
// The service is a thin layer over the SDK account and users APIs. The only
// logic it owns is deciding which credential a request carries and how to turn
// it into a *models.User.
package auth

import (
	"errors"

	"github.com/appwrite/sdk-for-go/v7/client"
	"github.com/appwrite/sdk-for-go/v7/users"
	"github.com/gin-gonic/gin"
)

// Service exposes the auth operations the library supports.
type Service struct {
	sdk        client.Client
	users      *users.Users
	cookieName string

	// verify turns a session secret into a user. It is a field so tests can
	// substitute a fake without a server; production uses sdkVerify.
	verify Verifier
}

// Verifier resolves a session secret to the authenticated user.
type Verifier interface {
	Verify(secret string) (*User, error)
}

// User is the authenticated identity stored in the Gin context.
type User struct {
	ID    string
	Email string
	Name  string
}

// New builds an auth Service around an SDK client.
//
// cookieName is the session cookie to read when no custom extraction is
// supplied; pass an empty string to fall back to the Appwrite convention.
func New(c client.Client, cookieName string) *Service {
	s := &Service{
		sdk:        c,
		users:      users.New(c),
		cookieName: cookieName,
	}
	s.verify = sdkVerify{s: s}
	return s
}

// NewWithVerifier builds a Service that delegates session verification to v.
// Used by tests and by apps that resolve sessions themselves.
func NewWithVerifier(v Verifier, cookieName string) *Service {
	return &Service{verify: v, cookieName: cookieName}
}

// CookieName returns the session cookie name this service reads.
func (s *Service) CookieName() string {
	return s.cookieName
}

// Users returns the underlying SDK users service.
func (s *Service) Users() *users.Users {
	return s.users
}

var (
	// ErrNotImplemented is returned by stubs awaiting implementation.
	ErrNotImplemented = errors.New("auth: not implemented")

	// ErrNoSession reports a request with no usable session credential.
	ErrNoSession = errors.New("auth: no session")

	// ErrSessionInvalid reports a session secret that could not be verified.
	ErrSessionInvalid = errors.New("auth: session invalid")
)

// UserKey is the Gin context key holding the authenticated *User.
const UserKey = "appwrite:user"

// userFromContext reads the authenticated user from a Gin context.
func userFromContext(c *gin.Context) (*User, bool) {
	v, ok := c.Get(UserKey)
	if !ok {
		return nil, false
	}
	u, ok := v.(*User)
	return u, ok
}
