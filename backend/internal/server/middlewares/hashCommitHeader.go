package middlewares

import (
	"net/http"
	"os"
)

// AddBuildHashCommitHeader adds a custom header to all responses with the hash commit used to
// build the Docker image being used by the backend.
// If the backend is run locally, no hash commit is sent in the response headers.
func (middleware *Middlewares) AddBuildHashCommitHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		hashCommit, ok := os.LookupEnv("BUILD_HASH_COMMIT")
		if !ok {
			hashCommit = "local"
		}
		w.Header().Set("x-raspall-hash-commit", hashCommit)

		next.ServeHTTP(w, req)
	})
}
