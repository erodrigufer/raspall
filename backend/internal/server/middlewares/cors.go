package middlewares

import (
	"net/http"
	"sync"

	"github.com/rs/cors"
)

func (m *Middlewares) Cors(next http.Handler) http.Handler {
	var (
		singleExecution sync.Once
		c               *cors.Cors
	)

	// By looking at the codebase of the cors.New() function, this function might
	// be a costly operation, so only execute New() once and store its value in
	// the variable within the closure.
	singleExecution.Do(func() {
		c = cors.New(cors.Options{
			AllowedOrigins:   []string{"http://localhost:*"},
			AllowCredentials: true,
		})
	})

	return c.Handler(next)
}
