package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/your-org/vaultlens/internal/config"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("write temp config: %v", err)
	}
	return p
}

func TestLoadDefaults(t *testing.T) {
	cfg, err := config.Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Vault.Address != "http://127.0.0.1:8200" {
		t.Errorf("default address mismatch: %s", cfg.Vault.Address)
	}
	if cfg.Vault.CacheTTL != 5*time.Minute {
		t.Errorf("default cache_ttl mismatch: %v", cfg.Vault.CacheTTL)
	}
	if cfg.Search.MaxResults != 50 {
		t.Errorf("default max_results mismatch: %d", cfg.Search.MaxResults)
	}
}

func TestLoadFromFile(t *testing.T) {
	yaml := `
vault:
  address: "https://vault.example.com"
  token: "s.abc"
  mount: "kv"
  cache_ttl: 2m
search:
  max_results: 10
`
	p := writeTempConfig(t, yaml)
	cfg, err := config.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Vault.Address != "https://vault.example.com" {
		t.Errorf("address mismatch: %s", cfg.Vault.Address)
	}
	if cfg.Vault.Mount != "kv" {
		t.Errorf("mount mismatch: %s", cfg.Vault.Mount)
	}
	if cfg.Search.MaxResults != 10 {
		t.Errorf("max_results mismatch: %d", cfg.Search.MaxResults)
	}
}

func TestEnvOverridesFileValues(t *testing.T) {
	t.Setenv("VAULT_ADDR", "https://env.vault.io")
	t.Setenv("VAULT_TOKEN", "s.env_token")

	cfg, err := config.Load("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Vault.Address != "https://env.vault.io" {
		t.Errorf("env address not applied: %s", cfg.Vault.Address)
	}
	if cfg.Vault.Token != "s.env_token" {
		t.Errorf("env token not applied: %s", cfg.Vault.Token)
	}
}

func TestLoadInvalidYAML(t *testing.T) {
	p := writeTempConfig(t, "{{{not yaml")
	_, err := config.Load(p)
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}
}

func TestValidateMissingMount(t *testing.T) {
	yaml := `vault:\n  mount: ""`
	p := writeTempConfig(t, yaml)
	_, err := config.Load(p)
	if err == nil {
		t.Fatal("expected validation error for empty mount")
	}
}
