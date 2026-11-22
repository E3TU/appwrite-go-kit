package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/e3tu/appwrite-base"
)

// verify runs connectivity and ID existence checks against the deployment.
//
// It checks the health endpoints, then asks the SDK whether the configured
// project, database, table and bucket ids exist. Every check reports one line;
// any failure makes the command exit non-zero.
func verify() error {
	cfg, err := loadEnv()
	if err != nil {
		return err
	}

	if _, err := base.New(cfg); err != nil {
		return err
	}

	v := &verifier{
		client:   &http.Client{Timeout: base.DefaultTimeout},
		endpoint: strings.TrimRight(cfg.Endpoint, "/"),
		stack:    cfg.ProjectID,
	}

	var failed int
	check := func(name string, err error) {
		if err != nil {
			failed++
			fmt.Printf("FAIL  %-12s %v\n", name, err)
			return
		}
		fmt.Printf("ok    %s\n", name)
	}

	check("health", v.getJSON("/health/status", nil))
	check("version", v.getJSON("/health/version", nil))
	check("project", v.getJSON("/projects/"+url.PathEscape(cfg.ProjectID), nil))

	if cfg.DB != nil && cfg.DB.DatabaseID != "" {
		check("database", v.getJSON("/tablesdb/"+url.PathEscape(cfg.DB.DatabaseID), nil))
	}
	if s := cfg.Storage; s != nil && s.DefaultBucketID != "" {
		check("bucket", v.getJSON("/storage/buckets/"+url.PathEscape(s.DefaultBucketID), nil))
	}

	if failed > 0 {
		return fmt.Errorf("%d check(s) failed", failed)
	}

	fmt.Println("all checks passed")
	return nil
}

type verifier struct {
	client   *http.Client
	endpoint string
	stack    string
}

func (v *verifier) getJSON(path string, dst any) error {
	ctx, cancel := context.WithTimeout(context.Background(), base.DefaultTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, v.endpoint+path, nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Appwrite-Project", v.stack)
	req.Header.Set("X-Appwrite-Response-Format", base.ResponseFormat)

	resp, err := v.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, shorten(body))
	}

	if dst != nil && len(body) > 0 {
		if err := json.Unmarshal(body, dst); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}
	return nil
}

func shorten(b []byte) string {
	s := strings.TrimSpace(string(b))
	if len(s) > 200 {
		return s[:200] + "..."
	}
	return s
}
