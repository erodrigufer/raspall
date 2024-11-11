package server

import (
	"net/http"

	m "github.com/erodrigufer/raspall/internal/server/middlewares"
)

func (app *Application) routes() http.Handler {
	mws := m.NewMiddlewares(app.InfoLog, app.ErrorLog, app.sessionManager)

	globalMiddlewares := m.MiddlewareChain(app.sessionManager.LoadAndSave, mws.LogRequest, mws.Cors, mws.RecoverPanic)

	fileServer := http.StripPrefix("/static", http.FileServer(http.Dir("./static")))

	mux := http.NewServeMux()
	mux.Handle("GET /login", mws.AuthenticateLogin(app.getLogin()))
	mux.Handle("POST /login", app.postLogin())
	mux.Handle("POST /logout", app.postLogout())
	mux.Handle("GET /static/", mws.PublicCacheCacheControl(fileServer))
	mux.Handle("GET /api/health", app.health())

	protectedMux := http.NewServeMux()
	mux.Handle("/", mws.Authenticate(mws.PrivateCacheControl(protectedMux)))

	protectedMux.Handle("GET /", app.index())
	protectedMux.Handle("POST /nacio", app.nacio())
	protectedMux.Handle("POST /hn", app.hn())
	protectedMux.Handle("POST /lobsters", app.lobsters())

	return globalMiddlewares(mux)
}
