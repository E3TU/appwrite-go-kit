// Package base wires one self-hosted Appwrite project into Go (Gin)
// services: configuration, the SDK client, and accessors for each enabled
// capability.
//
// No package in the library reads environment variables and no constructor
// performs network I/O.
package base

import (
	"github.com/appwrite/sdk-for-go/v7/client"
	"github.com/e3tu/appwrite-base/auth"
	"github.com/e3tu/appwrite-base/database"
	"github.com/e3tu/appwrite-base/functions"
	"github.com/e3tu/appwrite-base/storage"
)

// Base holds one configured Appwrite client plus the capabilities enabled in
// Config.
//
// Accessors are the only supported way in: each returns (service, true) when
// the capability is enabled and (nil, false) when it is not, so callers never
// have to inspect struct fields.
type Base struct {
	cfg Config
	sdk client.Client

	auth      *auth.Service
	storage   *storage.Service
	db        *database.Service
	functions *functions.Service
}

// New validates cfg, builds the SDK client and wires every enabled
// capability. It performs no network I/O.
//
// A Config with all capability blocks nil is valid and yields a client-only
// Base.
func New(cfg Config) (*Base, error) {
	c, err := newClient(cfg)
	if err != nil {
		return nil, err
	}

	b := &Base{cfg: cfg, sdk: c}

	if cfg.Auth != nil {
		b.auth = auth.New(c, sessionCookieName(cfg))
	}
	if cfg.Storage != nil {
		b.storage = storage.New(c, cfg.Storage.DefaultBucketID, cfg.PublicEndpoint, cfg.ProjectID)
	}
	if cfg.DB != nil {
		b.db = database.New(c, cfg.DB.DatabaseID)
	}
	if cfg.Functions != nil {
		b.functions = functions.New(c)
	}

	return b, nil
}

// Auth returns the auth service, or (nil, false) when auth is disabled.
func (b *Base) Auth() (*auth.Service, bool) {
	return b.auth, b.auth != nil
}

// Storage returns the storage service, or (nil, false) when storage is
// disabled.
func (b *Base) Storage() (*storage.Service, bool) {
	return b.storage, b.storage != nil
}

// DB returns the database service, or (nil, false) when the database
// capability is disabled.
func (b *Base) DB() (*database.Service, bool) {
	return b.db, b.db != nil
}

// Functions returns the functions service, or (nil, false) when functions are
// disabled.
func (b *Base) Functions() (*functions.Service, bool) {
	return b.functions, b.functions != nil
}

// Config returns the configuration Base was built from.
func (b *Base) Config() Config {
	return b.cfg
}

// Client returns the configured SDK client.
func (b *Base) Client() client.Client {
	return b.sdk
}
