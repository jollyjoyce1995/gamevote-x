package server

import (
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"
)

// staticHandler serves the built Vue app, falling back to index.html so that
// client-side routes (/parties/ABC123) survive a page reload.
func staticHandler(dir string) http.Handler {
	fsys := os.DirFS(dir)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
		if info, err := fs.Stat(fsys, name); name == "" || err != nil || info.IsDir() {
			name = "index.html"
		}
		if strings.HasPrefix(name, "assets/") {
			// Vite fingerprints these file names.
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		} else {
			w.Header().Set("Cache-Control", "no-cache")
		}
		http.ServeFileFS(w, r, fsys, name)
	})
}
