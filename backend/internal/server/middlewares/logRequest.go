package middlewares

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/urfave/negroni"
)

// LogRequest logs every client's request.
// Log IP address of client, protocol used, HTTP method and requested URL.
func (m *Middlewares) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Without negroni it is not possible to log the HTTP status code of the response
		// sent in another moment of the middleware chain.
		negroniRW := negroni.NewResponseWriter(w)

		next.ServeHTTP(negroniRW, r)

		var ipReqChain string
		forwardedHeader := r.Header.Values("X-Forwarded-For")
		if len(forwardedHeader) > 0 {
			ipReqChain = strings.Join(forwardedHeader, ", ")
		} else {
			ipReqChain = r.RemoteAddr
		}

		reqDuration := time.Since(startTime).Milliseconds()
		m.infoLog.Info("received HTTP request",
			slog.String("method", r.Method),
			slog.String("requested_url", r.URL.String()),
			slog.Int64("duration_ms", reqDuration),
			slog.Int("status", negroniRW.Status()),
			slog.String("ip_request_chain", ipReqChain))
	})
}
