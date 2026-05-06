package config

import (
	"fmt"
	"os"
)

const (
	// DefaultConfigPath is the conventional location for the config file.
	DefaultConfigPath = ".vaultlens.yaml"

	// EnvConfigPath is the environment variable that overrides the config path.
	EnvConfigPath = "VAULTLENS_CONFIG"
)

// Resolve returns the effective config file path by checking, in order:
//  1. The explicit path argument (if non-empty)
//  2. The VAULTLENS_CONFIG environment variable
//  3. The default path (.vaultlens.yaml) if it exists on disk
//
// Returns an empty string when no file is found, which causes Load to use
// built-in defaults.
func Resolve(explicit string) string {
	if explicit != "" {
		return explicit
	}
	if v := os.Getenv(EnvConfigPath); v != "" {
		return v
	}
	if _, err := os.Stat(DefaultConfigPath); err == nil {
		return DefaultConfigPath
	}
	return ""
}

// MustLoad calls Load and panics on error. Intended for use during
// application startup where a misconfigured environment is fatal.
func MustLoad(path string) *Config {
	cfg, err := Load(path)
	if err != nil {
		panic(fmt.Sprintf("vaultlens: failed to load config: %v", err))
	}
	return cfg
}
