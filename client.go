package base

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/appwrite/sdk-for-go/client"
)

// newClient creates the Appwrite SDK client used by enabled capabilities.
//
// The client is built from Config.Endpoint only, never PublicEndpoint, and the
// constructor performs no network I/O. The SDK at the pinned version exposes
// no option helpers, so fields are set directly.
func newClient(cfg Config) (client.Client, error) {
	if err := validateClientConfig(cfg); err != nil {
		return client.Client{}, err
	}

	c := client.New()
	c.Endpoint = strings.TrimRight(cfg.Endpoint, "/")
	c.SelfSigned = cfg.SelfSigned
	if cfg.Timeout > 0 {
		c.Timeout = cfg.Timeout
	} else {
		c.Timeout = DefaultTimeout
	}
	// The SDK keeps one timeout for the whole client; the http.Client built by
	// SDK New keeps its own, so keep them in step.
	c.Client.Timeout = c.Timeout

	c.AddHeader("X-Appwrite-Project", cfg.ProjectID)
	if cfg.APIKey != "" {
		c.AddHeader("X-Appwrite-Key", cfg.APIKey)
	}

	return c, nil
}

// validateClientConfig checks the fields needed to build any client at all.
//
// A Config with every capability disabled is still valid: it yields a
// client-only Base.
func validateClientConfig(cfg Config) error {
	if cfg.ProjectID == "" {
		return fmt.Errorf("%w: ProjectID is required", ErrInvalidConfig)
	}
	if cfg.Endpoint == "" {
		return fmt.Errorf("%w: Endpoint is required", ErrInvalidConfig)
	}

	u, err := url.Parse(cfg.Endpoint)
	if err != nil {
		return fmt.Errorf("%w: Endpoint: %v", ErrInvalidConfig, err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("%w: Endpoint must be http or https, got %q", ErrInvalidConfig, u.Scheme)
	}
	if u.Host == "" {
		return fmt.Errorf("%w: Endpoint has no host: %q", ErrInvalidConfig, cfg.Endpoint)
	}

	if cfg.Timeout < 0 {
		return fmt.Errorf("%w: Timeout must not be negative", ErrInvalidConfig)
	}

	return nil
}

// sessionCookieName resolves the session cookie name for cfg.
//
// Capability config wins over the top-level field, which wins over the
// Appwrite convention.
func sessionCookieName(cfg Config) string {
	name := cfg.CookieName
	if cfg.Auth != nil && cfg.Auth.CookieName != "" {
		name = cfg.Auth.CookieName
	}
	if name != "" {
		return name
	}
	return "a_session_" + cfg.ProjectID
}

// resolveTimeout returns the effective request timeout.
func resolveTimeout(cfg Config) time.Duration {
	if cfg.Timeout > 0 {
		return cfg.Timeout
	}
	return DefaultTimeout
}
