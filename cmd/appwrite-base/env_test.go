package main

import "testing"

func TestEnvBool(t *testing.T) {
	tests := []struct {
		value   string
		want    bool
		wantErr bool
	}{
		{"", false, false},
		{"true", true, false},
		{"1", true, false},
		{"false", false, false},
		{"0", false, false},
		{"yes please", false, true},
	}
	for _, tt := range tests {
		t.Setenv("TEST_BOOL", tt.value)
		got, err := envBool("TEST_BOOL")
		if (err != nil) != tt.wantErr {
			t.Errorf("%q: error = %v, wantErr %v", tt.value, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && got != tt.want {
			t.Errorf("%q: got %v, want %v", tt.value, got, tt.want)
		}
	}
}

func TestLoadEnvDefaults(t *testing.T) {
	t.Setenv("APPWRITE_PROJECT_ID", "proj")
	t.Setenv("APPWRITE_ENDPOINT", "http://appwrite.local/v1")

	cfg, err := loadEnv()
	if err != nil {
		t.Fatalf("loadEnv() error = %v", err)
	}
	if cfg.CookieName != "a_session_proj" {
		t.Errorf("CookieName = %q", cfg.CookieName)
	}
	if cfg.Auth != nil || cfg.Storage != nil || cfg.DB != nil || cfg.Functions != nil {
		t.Error("capabilities must stay disabled by default")
	}
}

func TestLoadEnvCapabilities(t *testing.T) {
	t.Setenv("APPWRITE_PROJECT_ID", "proj")
	t.Setenv("APPWRITE_ENDPOINT", "http://appwrite.local/v1")
	t.Setenv("APPWRITE_DB_ENABLED", "true")
	t.Setenv("APPWRITE_DATABASE_ID", "main")

	cfg, err := loadEnv()
	if err != nil {
		t.Fatalf("loadEnv() error = %v", err)
	}
	if cfg.DB == nil || cfg.DB.DatabaseID != "main" {
		t.Fatalf("DB = %+v, want database main", cfg.DB)
	}
}

func TestLoadEnvRejectsBadTimeout(t *testing.T) {
	t.Setenv("APPWRITE_TIMEOUT", "soon")

	if _, err := loadEnv(); err == nil {
		t.Fatal("expected error for malformed duration")
	}
}
