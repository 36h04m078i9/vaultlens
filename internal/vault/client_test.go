package vault_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/vaultlens/internal/vault"
)

func newMockVaultServer(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()

	// KV v2 list endpoint
	mux.HandleFunc("/v1/secret/metadata/myapp", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "LIST" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{ //nolint:errcheck
			"data": map[string]interface{}{
				"keys": []string{"db-password", "api-key"},
			},
		})
	})

	// KV v2 read endpoint
	mux.HandleFunc("/v1/secret/data/myapp/db-password", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{ //nolint:errcheck
			"data": map[string]interface{}{
				"data": map[string]interface{}{
					"password": "s3cr3t",
				},
			},
		})
	})

	return httptest.NewServer(mux)
}

func TestListSecrets(t *testing.T) {
	srv := newMockVaultServer(t)
	defer srv.Close()

	client, err := vault.NewClient(vault.Config{
		Address: srv.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	keys, err := client.ListSecrets(context.Background(), "secret", "myapp")
	if err != nil {
		t.Fatalf("ListSecrets: %v", err)
	}

	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
	if keys[0] != "db-password" || keys[1] != "api-key" {
		t.Errorf("unexpected keys: %v", keys)
	}
}

func TestReadSecret(t *testing.T) {
	srv := newMockVaultServer(t)
	defer srv.Close()

	client, err := vault.NewClient(vault.Config{
		Address: srv.URL,
		Token:   "test-token",
	})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	data, err := client.ReadSecret(context.Background(), "secret", "myapp/db-password")
	if err != nil {
		t.Fatalf("ReadSecret: %v", err)
	}

	if pw, ok := data["password"]; !ok || pw != "s3cr3t" {
		t.Errorf("expected password=s3cr3t, got %v", data)
	}
}
