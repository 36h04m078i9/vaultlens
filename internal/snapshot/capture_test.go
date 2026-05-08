package snapshot

import (
	"errors"
	"strings"
	"testing"
)

// mockLister simulates a flat Vault path tree.
type mockLister struct {
	data map[string][]string
	failOn string
}

func (m *mockLister) ListSecrets(path string) ([]string, error) {
	if m.failOn != "" && strings.HasPrefix(path, m.failOn) {
		return nil, errors.New("mock list error")
	}
	return m.data[path], nil
}

func TestCaptureStoresSnapshot(t *testing.T) {
	lister := &mockLister{
		data: map[string][]string{
			"secret/": {"foo", "bar"},
		},
	}
	store := NewStore()
	snap, err := Capture("test", "secret/", lister, store)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap.Label != "test" {
		t.Errorf("label: got %q", snap.Label)
	}
	if len(snap.Paths) != 2 {
		t.Errorf("expected 2 paths, got %d", len(snap.Paths))
	}
	_, ok := store.Get(snap.ID)
	if !ok {
		t.Fatal("snapshot not persisted in store")
	}
}

func TestCaptureWalksSubdirectories(t *testing.T) {
	lister := &mockLister{
		data: map[string][]string{
			"secret/":     {"sub/", "top"},
			"secret/sub/": {"nested"},
		},
	}
	store := NewStore()
	snap, err := Capture("walk", "secret/", lister, store)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(snap.Paths) != 2 {
		t.Errorf("expected 2 leaf paths, got %d: %v", len(snap.Paths), snap.Paths)
	}
}

func TestCaptureNilListerReturnsError(t *testing.T) {
	_, err := Capture("x", "secret/", nil, NewStore())
	if err == nil {
		t.Fatal("expected error for nil lister")
	}
}

func TestCaptureNilStoreReturnsError(t *testing.T) {
	_, err := Capture("x", "secret/", &mockLister{}, nil)
	if err == nil {
		t.Fatal("expected error for nil store")
	}
}

func TestCaptureListErrorPropagates(t *testing.T) {
	lister := &mockLister{failOn: "secret/"}
	_, err := Capture("x", "secret/", lister, NewStore())
	if err == nil {
		t.Fatal("expected error from lister")
	}
}
