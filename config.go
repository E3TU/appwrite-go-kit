package base

import "time"

// Config describes one self-hosted Appwrite deployment.
//
// Endpoint is used for server-to-server SDK calls and PublicEndpoint is used
// for browser-facing URLs. The two are never merged: an internal container
// address can never be handed to a browser.
//
// The capability blocks are optional and independent. A nil block means the
// capability is disabled: its accessor on Base returns (nil, false) and no
// SDK service is built for it.
type Config struct {
	// ProjectID is the Appwrite project identifier. Required.
	ProjectID string

	// APIKey is the server-side API key used for admin operations. Required
	// when any capability is enabled.
	APIKey string

	// Endpoint is the internal base URL used by the SDK, e.g.
	// "http://appwrite.local/v1". Required.
	Endpoint string

	// PublicEndpoint is the externally reachable base URL, used to build
	// browser-facing links. Required when storage is enabled.
	PublicEndpoint string

	// CookieName overrides the default session cookie name
	// ("a_session_<ProjectID>") for auth.
	CookieName string

	// Timeout bounds each SDK request. Zero means DefaultTimeout.
	Timeout time.Duration

	// SelfSigned allows TLS connections to endpoints presenting self-signed
	// certificates.
	SelfSigned bool

	// Capability blocks. Nil means disabled.
	Auth      *AuthConfig
	Storage   *StorageConfig
	DB        *DBConfig
	Functions *FunctionsConfig
}

// AuthConfig configures the optional auth capability.
type AuthConfig struct {
	// SessionHeader, when set, is checked for the session secret before the
	// cookie is. Useful for native clients that cannot hold cookies.
	SessionHeader string

	// CookieName overrides Config.CookieName for this capability.
	CookieName string
}

// StorageConfig configures the optional storage capability.
type StorageConfig struct {
	// DefaultBucketID is used by helpers when no bucket is given explicitly.
	DefaultBucketID string
}

// DBConfig configures the optional database capability.
type DBConfig struct {
	// DatabaseID is the database the service is bound to.
	DatabaseID string
}

// FunctionsConfig configures the optional functions capability.
type FunctionsConfig struct{}
