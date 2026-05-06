package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultlens/internal/config"
)

func TestResolveExplicitPath(t *testing.T) {
	got := config.Resolve("/etc/vault.yaml")
	if got != "/etc/vault.yaml" {
		t.Errorf("expected explicit path, got %q", got)
	}
}

func TestResolveEnvVar(t *testing.T) {
	t.Setenv(config.EnvConfigPath, "/tmp/env-config.yaml")
	got := config.Resolve("")
	if got != "/tmp/env-config.yaml" {
		t.Errorf("expected env path, got %q", got)
	}
}

func TestResolveDefaultFileExists(t *testing.T) {
	// Change working directory to a temp dir containing the default file.
	dir := t.TempDir()
	defPath := filepath.Join(dir, config.DefaultConfigPath)
	if err := os.WriteFile(defPath, []byte("vault:\n  address: x\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	orig, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(orig) }) //nolint:errcheck

	got := config.Resolve("")
	if got != config.DefaultConfigPath {
		t.Errorf("expected default path, got %q", got)
	}
}

func TestResolveReturnsEmptyWhenNothingFound(t *testing.T) {
	dir := t.TempDir()
	orig, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chdir(orig) }) //nolint:errcheck

	got := config.Resolve("")
	if got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestMustLoadPanicsOnBadConfig(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic, got none")
		}
	}()
	p := writeTempConfig(t, "{{{invalid")
	config.MustLoad(p)
}
