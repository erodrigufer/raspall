package middlewares

import (
	"net/http"
)

// LogRequest logs every client's request.
// Log IP address of client, protocol used, HTTP method and requested URL.
func (m *Middlewares) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL)

		next.ServeHTTP(w, r)
	})
}
