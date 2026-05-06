package vault

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"
)

func newTreeMockServer(t *testing.T) *httptest.Server {
	t.Helper()
	mux := http.NewServeMux()

	// LIST secret/
	mux.HandleFunc("/v1/secret/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "LIST" {
			http.NotFound(w, r)
			return
		}
		switch r.URL.Path {
		case "/v1/secret/":
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"data":{"keys":["alpha","sub/"]}}`))
		case "/v1/secret/sub/":
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"data":{"keys":["beta","gamma"]}}`))
		default:
			http.NotFound(w, r)
		}
	})
	return httptest.NewServer(mux)
}

func TestWalkPaths(t *testing.T) {
	srv := newTreeMockServer(t)
	defer srv.Close()

	client, err := NewClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	paths, err := client.WalkPaths(context.Background(), "secret/")
	if err != nil {
		t.Fatalf("WalkPaths: %v", err)
	}

	sort.Strings(paths)
	want := []string{"secret/alpha", "secret/sub/beta", "secret/sub/gamma"}
	if len(paths) != len(want) {
		t.Fatalf("got %d paths, want %d: %v", len(paths), len(want), paths)
	}
	for i, p := range paths {
		if p != want[i] {
			t.Errorf("paths[%d] = %q, want %q", i, p, want[i])
		}
	}
}

func TestBuildTree(t *testing.T) {
	srv := newTreeMockServer(t)
	defer srv.Close()

	client, err := NewClient(srv.URL, "test-token")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	tree, err := client.BuildTree(context.Background(), "secret/")
	if err != nil {
		t.Fatalf("BuildTree: %v", err)
	}

	if len(tree.Children) != 2 {
		t.Fatalf("root children = %d, want 2", len(tree.Children))
	}

	var subNode *PathNode
	for _, c := range tree.Children {
		if c.Name == "sub/" {
			subNode = c
		}
	}
	if subNode == nil {
		t.Fatal("expected sub/ child node")
	}
	if subNode.IsLeaf {
		t.Error("sub/ should not be a leaf")
	}
	if len(subNode.Children) != 2 {
		t.Errorf("sub/ children = %d, want 2", len(subNode.Children))
	}
}
