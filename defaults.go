package base

import "time"

// DefaultTimeout bounds one SDK request when Config.Timeout is zero.
const DefaultTimeout = 30 * time.Second

// Defaults returns a Config with neutral, non-secret values.
//
// Endpoints, ProjectID and APIKey are deliberately empty: the library does not
// guess a deployment. Every capability block is nil, so a Config built from
// Defaults alone is a valid client-only Base.
func Defaults() Config {
	return Config{
		Timeout: DefaultTimeout,
	}
}
