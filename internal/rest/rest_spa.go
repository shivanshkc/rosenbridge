package rest

import (
	"net/http"
)

// serveFrontend serves the SPA present in the given directory.
func serveFrontend(frontendDir string) http.Handler {
	// http.Dir protects against directory traversal ("../../secrets")
	dir := http.Dir(frontendDir)
	server := http.FileServer(dir)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, err := dir.Open(r.URL.Path)
		if err != nil {
			// Serve index.html, if file not found, or any error.
			r.URL.Path = "/index.html"
			server.ServeHTTP(w, r)
			return
		}

		// Serve file, if found.
		_ = file.Close()
		server.ServeHTTP(w, r)
	})
}
