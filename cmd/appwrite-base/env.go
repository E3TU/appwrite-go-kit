package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/e3tu/appwrite-base"
)

// Environment variables read by this binary (and nowhere else in the module).
//
//	APPWRITE_PROJECT_ID        required  base.Config.ProjectID
//	APPWRITE_API_KEY           required  base.Config.APIKey
//	APPWRITE_ENDPOINT          required  base.Config.Endpoint (internal)
//	APPWRITE_PUBLIC_ENDPOINT   required  base.Config.PublicEndpoint
//	APPWRITE_COOKIE_NAME       optional  base.Config.CookieName
//	APPWRITE_TIMEOUT           optional  base.Config.Timeout (Go duration)
//	APPWRITE_SELF_SIGNED       optional  base.Config.SelfSigned ("true"/"false")
//	APPWRITE_AUTH_ENABLED      optional  enables base.AuthConfig
//	APPWRITE_STORAGE_ENABLED   optional  enables base.StorageConfig
//	APPWRITE_STORAGE_BUCKET    optional  base.StorageConfig.DefaultBucketID
//	APPWRITE_FUNCTIONS_ENABLED optional  enables base.FunctionsConfig
//	APPWRITE_DB_ENABLED        optional  enables base.DBConfig
//	APPWRITE_DATABASE_ID       optional  base.DBConfig.DatabaseID
//
// Optional variables are read once; a malformed value is an error rather than
// a silent default.

// loadEnv builds a base.Config from the environment.
func loadEnv() (base.Config, error) {
	cfg := base.Defaults()

	cfg.ProjectID = strings.TrimSpace(os.Getenv("APPWRITE_PROJECT_ID"))
	cfg.APIKey = strings.TrimSpace(os.Getenv("APPWRITE_API_KEY"))
	cfg.Endpoint = strings.TrimSpace(os.Getenv("APPWRITE_ENDPOINT"))
	cfg.PublicEndpoint = strings.TrimSpace(os.Getenv("APPWRITE_PUBLIC_ENDPOINT"))
	if cfg.CookieName = strings.TrimSpace(os.Getenv("APPWRITE_COOKIE_NAME")); cfg.CookieName == "" {
		cfg.CookieName = "a_session_" + cfg.ProjectID
	}

	if raw := strings.TrimSpace(os.Getenv("APPWRITE_TIMEOUT")); raw != "" {
		d, err := time.ParseDuration(raw)
		if err != nil {
			return base.Config{}, fmt.Errorf("APPWRITE_TIMEOUT: %w", err)
		}
		if d <= 0 {
			return base.Config{}, fmt.Errorf("APPWRITE_TIMEOUT: must be positive, got %q", raw)
		}
		cfg.Timeout = d
	}

	selfSigned, err := envBool("APPWRITE_SELF_SIGNED")
	if err != nil {
		return base.Config{}, err
	}
	cfg.SelfSigned = selfSigned

	authOn, err := envBool("APPWRITE_AUTH_ENABLED")
	if err != nil {
		return base.Config{}, err
	}
	if authOn {
		cfg.Auth = &base.AuthConfig{}
	}

	storageOn, err := envBool("APPWRITE_STORAGE_ENABLED")
	if err != nil {
		return base.Config{}, err
	}
	if storageOn {
		cfg.Storage = &base.StorageConfig{
			DefaultBucketID: strings.TrimSpace(os.Getenv("APPWRITE_STORAGE_BUCKET")),
		}
	}

	dbOn, err := envBool("APPWRITE_DB_ENABLED")
	if err != nil {
		return base.Config{}, err
	}
	if dbOn {
		cfg.DB = &base.DBConfig{
			DatabaseID: strings.TrimSpace(os.Getenv("APPWRITE_DATABASE_ID")),
		}
	}

	functionsOn, err := envBool("APPWRITE_FUNCTIONS_ENABLED")
	if err != nil {
		return base.Config{}, err
	}
	if functionsOn {
		cfg.Functions = &base.FunctionsConfig{}
	}

	return cfg, nil
}

// listenAddr returns the address serve binds to.
//
// TODO: move APPWRITE_LISTEN into the documented variable list and into
// Config once the config shape allows it.
func listenAddr() string {
	if addr := strings.TrimSpace(os.Getenv("APPWRITE_LISTEN")); addr != "" {
		return addr
	}
	return ":8080"
}

// envBool reads an optional boolean. An unset variable is false.
func envBool(key string) (bool, error) {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return false, nil
	}
	v, err := strconv.ParseBool(raw)
	if err != nil {
		return false, fmt.Errorf("%s: %w", key, err)
	}
	return v, nil
}
