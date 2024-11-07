package server

import (
	"net/http"

	m "github.com/erodrigufer/raspall/internal/server/middlewares"
)

func (app *Application) routes() http.Handler {
	mws := m.NewMiddlewares(app.InfoLog, app.ErrorLog)

	globalMiddlewares := m.MiddlewareChain(mws.LogRequest, mws.Cors, mws.RecoverPanic)

	fileServer := http.StripPrefix("/static", http.FileServer(http.Dir("./static")))

	mux := http.NewServeMux()
	mux.Handle("GET /", app.index())
	mux.Handle("GET /static/", fileServer)
	mux.Handle("GET /login", app.getLogin())
	mux.Handle("POST /login", app.postLogin())
	mux.Handle("POST /nacio", mws.CacheControl(app.nacio()))
	mux.Handle("POST /hn", mws.CacheControl(app.hn()))
	mux.Handle("POST /lobsters", mws.CacheControl(app.lobsters()))

	mux.Handle("GET /api/health", app.health())

	return globalMiddlewares(mux)
}
