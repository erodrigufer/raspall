package middlewares

import (
	"fmt"
	"net/http"

	"github.com/erodrigufer/raspall/internal/utils"
)

// RecoverPanic sends an Internal Server Error message code to a client, when
// the server has to close an HTTP connection with a client due to a panic
// inside the goroutine handling the client.
// This should be the first middleware defined in the chain of middlewares, so
// that it can capture the panics of all other further middlewares.
func (m *Middlewares) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function which will always be executed in the event
		// of a panic as Go unwinds the stack.
		defer func() {
			// Use the built-in recover function to check if there has been a
			// panic or not.
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				utils.HandleServerError(w, fmt.Errorf("%s", err), m.errorLog)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
