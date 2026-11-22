package base

import (
	"errors"
	"testing"
	"time"
)

func TestDefaultsIsClientOnly(t *testing.T) {
	cfg := Defaults()
	if cfg.Timeout != DefaultTimeout {
		t.Fatalf("Timeout = %v, want %v", cfg.Timeout, DefaultTimeout)
	}
	if cfg.Auth != nil || cfg.Storage != nil || cfg.DB != nil || cfg.Functions != nil {
		t.Fatal("Defaults must leave every capability disabled")
	}
}

func TestNewRejectsMissingProject(t *testing.T) {
	cfg := Defaults()
	cfg.Endpoint = "http://appwrite.local/v1"

	_, err := New(cfg)
	if !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("New() error = %v, want ErrInvalidConfig", err)
	}
}

func TestNewRejectsBadEndpoint(t *testing.T) {
	for _, endpoint := range []string{
		"appwrite.local/v1",
		"ftp://appwrite.local/v1",
		"http://",
	} {
		cfg := Defaults()
		cfg.ProjectID = "p"
		cfg.Endpoint = endpoint

		if _, err := New(cfg); !errors.Is(err, ErrInvalidConfig) {
			t.Errorf("endpoint %q: error = %v, want ErrInvalidConfig", endpoint, err)
		}
	}
}

func TestNewRejectsNegativeTimeout(t *testing.T) {
	cfg := Defaults()
	cfg.ProjectID = "p"
	cfg.Endpoint = "http://appwrite.local/v1"
	cfg.Timeout = -time.Second

	if _, err := New(cfg); !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("New() error = %v, want ErrInvalidConfig", err)
	}
}

func TestNewClientOnly(t *testing.T) {
	cfg := Defaults()
	cfg.ProjectID = "p"
	cfg.Endpoint = "http://appwrite.local/v1/"

	b, err := New(cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	if _, ok := b.Auth(); ok {
		t.Error("Auth() enabled on a client-only Base")
	}
	if _, ok := b.Storage(); ok {
		t.Error("Storage() enabled on a client-only Base")
	}
	if _, ok := b.DB(); ok {
		t.Error("DB() enabled on a client-only Base")
	}
	if _, ok := b.Functions(); ok {
		t.Error("Functions() enabled on a client-only Base")
	}
	if got := b.Client().Endpoint; got != "http://appwrite.local/v1" {
		t.Errorf("Endpoint = %q, want trailing slash trimmed", got)
	}
}

func TestNewEnablesOnlyConfiguredCapabilities(t *testing.T) {
	cfg := Defaults()
	cfg.ProjectID = "p"
	cfg.Endpoint = "http://appwrite.local/v1"
	cfg.DB = &DBConfig{DatabaseID: "main"}

	b, err := New(cfg)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	db, ok := b.DB()
	if !ok {
		t.Fatal("DB() disabled, want enabled")
	}
	if db.DatabaseID() != "main" {
		t.Errorf("DatabaseID = %q, want main", db.DatabaseID())
	}
	if _, ok := b.Auth(); ok {
		t.Error("Auth() enabled, want disabled")
	}
}

func TestSessionCookieName(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
		want string
	}{
		{"default", Config{ProjectID: "p"}, "a_session_p"},
		{"config overrides", Config{ProjectID: "p", CookieName: "sid"}, "sid"},
		{"capability wins", Config{ProjectID: "p", CookieName: "sid", Auth: &AuthConfig{CookieName: "auth"}}, "auth"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sessionCookieName(tt.cfg); got != tt.want {
				t.Errorf("sessionCookieName() = %q, want %q", got, tt.want)
			}
		})
	}
}
