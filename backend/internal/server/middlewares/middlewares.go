package middlewares

import (
	"log/slog"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

type middleware func(http.Handler) http.Handler

type Middlewares struct {
	infoLog               *slog.Logger
	errorLog              *slog.Logger
	sessionManager        *scs.SessionManager
	disableAuthentication bool
}

// NewMiddlewares creates a struct that contains all middlewares of the application.
// It injects all required dependencies for the middlewares.
func NewMiddlewares(infoLog, errorLog *slog.Logger, sessionManager *scs.SessionManager, disableAuthentication bool) *Middlewares {
	middlewares := new(Middlewares)
	middlewares.infoLog = infoLog
	middlewares.errorLog = errorLog
	middlewares.sessionManager = sessionManager
	middlewares.disableAuthentication = disableAuthentication

	return middlewares
}

func MiddlewareChain(middlewares ...middleware) middleware {
	return func(handler http.Handler) http.Handler {
		for _, mw := range middlewares {
			handler = mw(handler)
		}
		return handler
	}
}
