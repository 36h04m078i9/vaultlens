// Package history provides an in-memory, size-bounded list of recently
// accessed Vault secret paths for the current vaultlens session.
//
// Entries are stored most-recent-first and deduplicated: pushing a path
// that already exists moves it to the front rather than adding a duplicate.
//
// Usage:
//
//	h := history.New(50)
//	h.Push("secret/data/myapp/db")
//	entries := h.List() // []history.Entry{{Path: "secret/data/myapp/db", AccessedAt: ...}}
//
// History is safe for concurrent use.
package history
