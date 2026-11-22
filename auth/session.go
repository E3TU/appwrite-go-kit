package auth

import (
	"net/http"

	"github.com/appwrite/sdk-for-go/account"
)

// Session is the credential material found on an incoming request.
type Session struct {
	// Secret is the session secret (cookie value or header value).
	Secret string

	// Source records where the secret came from: SourceHeader or SourceCookie.
	Source string
}

const (
	// SourceHeader means the secret came from the configured session header.
	SourceHeader = "header"

	// SourceCookie means the secret came from the session cookie.
	SourceCookie = "cookie"
)

// ExtractSession reads the session secret from r.
//
// Lookup order is fixed for the whole module: the configured session header
// first (when header is non-empty), then the session cookie named cookieName.
// Missing credentials return ErrNoSession.
func ExtractSession(r *http.Request, header, cookieName string) (Session, error) {
	if r == nil {
		return Session{}, ErrNoSession
	}

	if header != "" {
		if v := r.Header.Get(header); v != "" {
			return Session{Secret: v, Source: SourceHeader}, nil
		}
	}

	if cookieName != "" {
		if ck, err := r.Cookie(cookieName); err == nil && ck.Value != "" {
			return Session{Secret: ck.Value, Source: SourceCookie}, nil
		}
	}

	return Session{}, ErrNoSession
}

// sdkVerify resolves a session secret through the SDK account API.
//
// It builds a per-request client so the caller's session is never stored on
// the shared admin client.
type sdkVerify struct {
	s *Service
}

// Verify implements Verifier.
func (v sdkVerify) Verify(secret string) (*User, error) {
	if secret == "" {
		return nil, ErrNoSession
	}

	sessionClient := v.s.sdk
	sessionClient.AddHeader("X-Appwrite-Session", secret)

	acc := account.New(sessionClient)
	user, err := acc.Get()
	if err != nil {
		return nil, ErrSessionInvalid
	}

	return &User{
		ID:    user.Id,
		Email: user.Email,
		Name:  user.Name,
	}, nil
}
