package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds all application configuration.
type Config struct {
	Vault  VaultConfig  `yaml:"vault"`
	Audit  AuditConfig  `yaml:"audit"`
	Search SearchConfig `yaml:"search"`
	UI     UIConfig     `yaml:"ui"`
}

// VaultConfig holds Vault connection settings.
type VaultConfig struct {
	Address   string        `yaml:"address"`
	Token     string        `yaml:"token"`
	Mount     string        `yaml:"mount"`
	CacheTTL  time.Duration `yaml:"cache_ttl"`
	Timeout   time.Duration `yaml:"timeout"`
}

// AuditConfig holds audit logging settings.
type AuditConfig struct {
	FilePath string `yaml:"file_path"`
	Enabled  bool   `yaml:"enabled"`
}

// SearchConfig holds fuzzy search tuning parameters.
type SearchConfig struct {
	MaxResults int     `yaml:"max_results"`
	MinScore   float64 `yaml:"min_score"`
}

// UIConfig holds terminal UI settings.
type UIConfig struct {
	PageSize    int  `yaml:"page_size"`
	ShowPreview bool `yaml:"show_preview"`
}

// Load reads a YAML config file from path. Environment variables override
// vault address and token when set.
func Load(path string) (*Config, error) {
	cfg := defaults()

	if path != "" {
		f, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("config: open %q: %w", path, err)
		}
		defer f.Close()

		if err := yaml.NewDecoder(f).Decode(cfg); err != nil {
			return nil, fmt.Errorf("config: decode: %w", err)
		}
	}

	if v := os.Getenv("VAULT_ADDR"); v != "" {
		cfg.Vault.Address = v
	}
	if v := os.Getenv("VAULT_TOKEN"); v != "" {
		cfg.Vault.Token = v
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func defaults() *Config {
	return &Config{
		Vault: VaultConfig{
			Address:  "http://127.0.0.1:8200",
			Mount:    "secret",
			CacheTTL: 5 * time.Minute,
			Timeout:  10 * time.Second,
		},
		Audit: AuditConfig{
			Enabled: true,
		},
		Search: SearchConfig{
			MaxResults: 50,
			MinScore:   0.1,
		},
		UI: UIConfig{
			PageSize:    20,
			ShowPreview: true,
		},
	}
}

func (c *Config) validate() error {
	if c.Vault.Address == "" {
		return fmt.Errorf("config: vault.address must not be empty")
	}
	if c.Vault.Mount == "" {
		return fmt.Errorf("config: vault.mount must not be empty")
	}
	if c.Search.MaxResults <= 0 {
		return fmt.Errorf("config: search.max_results must be > 0")
	}
	return nil
}
