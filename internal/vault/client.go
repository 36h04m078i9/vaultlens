package vault

import (
	"context"
	"fmt"
	"os"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client with read-only helpers.
type Client struct {
	vc *vaultapi.Client
}

// Config holds configuration for connecting to Vault.
type Config struct {
	Address string
	Token   string
	TLSSkipVerify bool
}

// NewClient creates a new read-only Vault client.
func NewClient(cfg Config) (*Client, error) {
	vcfg := vaultapi.DefaultConfig()

	addr := cfg.Address
	if addr == "" {
		addr = os.Getenv("VAULT_ADDR")
	}
	if addr == "" {
		addr = "http://127.0.0.1:8200"
	}
	vcfg.Address = addr

	if cfg.TLSSkipVerify {
		if err := vcfg.ConfigureTLS(&vaultapi.TLSConfig{Insecure: true}); err != nil {
			return nil, fmt.Errorf("configuring TLS: %w", err)
		}
	}

	vc, err := vaultapi.NewClient(vcfg)
	if err != nil {
		return nil, fmt.Errorf("creating vault client: %w", err)
	}

	token := cfg.Token
	if token == "" {
		token = os.Getenv("VAULT_TOKEN")
	}
	vc.SetToken(token)

	return &Client{vc: vc}, nil
}

// ListSecrets returns the keys at the given KV v2 path (list operation).
func (c *Client) ListSecrets(ctx context.Context, mountPath, secretPath string) ([]string, error) {
	logical := c.vc.Logical()
	listPath := fmt.Sprintf("%s/metadata/%s", mountPath, secretPath)

	secret, err := logical.ListWithContext(ctx, listPath)
	if err != nil {
		return nil, fmt.Errorf("listing secrets at %q: %w", listPath, err)
	}
	if secret == nil || secret.Data == nil {
		return []string{}, nil
	}

	rawKeys, ok := secret.Data["keys"]
	if !ok {
		return []string{}, nil
	}

	ifaces, ok := rawKeys.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected type for keys field")
	}

	keys := make([]string, 0, len(ifaces))
	for _, v := range ifaces {
		if s, ok := v.(string); ok {
			keys = append(keys, s)
		}
	}
	return keys, nil
}

// ReadSecret reads a single KV v2 secret at the given path.
func (c *Client) ReadSecret(ctx context.Context, mountPath, secretPath string) (map[string]interface{}, error) {
	readPath := fmt.Sprintf("%s/data/%s", mountPath, secretPath)
	secret, err := c.vc.Logical().ReadWithContext(ctx, readPath)
	if err != nil {
		return nil, fmt.Errorf("reading secret at %q: %w", readPath, err)
	}
	if secret == nil || secret.Data == nil {
		return map[string]interface{}{}, nil
	}
	data, _ := secret.Data["data"].(map[string]interface{})
	return data, nil
}
