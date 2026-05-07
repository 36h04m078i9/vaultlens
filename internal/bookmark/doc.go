// Package bookmark provides persistent bookmarking of Vault secret paths.
//
// Bookmarks are stored as a JSON file on disk and support add, remove, list,
// and search operations. The Store type is safe for concurrent use.
//
// Example usage:
//
//	store, err := bookmark.NewStore("/home/user/.vaultlens/bookmarks.json")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Save a path
//	_ = store.Add("secret/prod/db", "production database credentials")
//
//	// Search saved bookmarks
//	results := bookmark.Search(store.List(), "prod")
//	for _, r := range results {
//		fmt.Println(r.Bookmark.Path, r.Score)
//	}
package bookmark
