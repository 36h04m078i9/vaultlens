package vault

import (
	"context"
	"fmt"
	"strings"
)

// PathNode represents a node in the Vault secret path tree.
type PathNode struct {
	Name     string
	FullPath string
	IsLeaf   bool
	Children []*PathNode
}

// WalkPaths recursively walks all secret paths under the given root,
// returning a flat list of fully-qualified secret paths (leaves only).
func (c *Client) WalkPaths(ctx context.Context, root string) ([]string, error) {
	root = normalizePath(root)
	var paths []string
	if err := c.walk(ctx, root, &paths); err != nil {
		return nil, fmt.Errorf("walk %q: %w", root, err)
	}
	return paths, nil
}

func (c *Client) walk(ctx context.Context, path string, acc *[]string) error {
	keys, err := c.ListSecrets(ctx, path)
	if err != nil {
		return err
	}
	for _, key := range keys {
		full := strings.TrimRight(path, "/") + "/" + key
		if strings.HasSuffix(key, "/") {
			// directory node — recurse
			if err := c.walk(ctx, full, acc); err != nil {
				return err
			}
		} else {
			*acc = append(*acc, full)
		}
	}
	return nil
}

// BuildTree constructs a PathNode tree rooted at the given path.
func (c *Client) BuildTree(ctx context.Context, root string) (*PathNode, error) {
	root = normalizePath(root)
	node := &PathNode{Name: root, FullPath: root}
	if err := c.buildTree(ctx, root, node); err != nil {
		return nil, fmt.Errorf("build tree %q: %w", root, err)
	}
	return node, nil
}

func (c *Client) buildTree(ctx context.Context, path string, parent *PathNode) error {
	keys, err := c.ListSecrets(ctx, path)
	if err != nil {
		return err
	}
	for _, key := range keys {
		full := strings.TrimRight(path, "/") + "/" + key
		child := &PathNode{
			Name:     key,
			FullPath: full,
			IsLeaf:   !strings.HasSuffix(key, "/"),
		}
		if !child.IsLeaf {
			if err := c.buildTree(ctx, full, child); err != nil {
				return err
			}
		}
		parent.Children = append(parent.Children, child)
	}
	return nil
}

func normalizePath(p string) string {
	return strings.TrimRight(p, "/")
}
