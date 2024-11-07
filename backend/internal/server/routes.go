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
	mux.Handle("GET /static/", fileServer)
	mux.Handle("GET /api/health", app.health())

	protectedMux := http.NewServeMux()
	mux.Handle("/", mws.Authenticate(protectedMux))

	protectedMux.Handle("GET /", app.index())
	protectedMux.Handle("POST /nacio", mws.CacheControl(app.nacio()))
	protectedMux.Handle("POST /hn", mws.CacheControl(app.hn()))
	protectedMux.Handle("POST /lobsters", mws.CacheControl(app.lobsters()))

	return globalMiddlewares(mux)
}
