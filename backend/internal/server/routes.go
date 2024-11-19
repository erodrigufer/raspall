package server

import (
	"net/http"

	"github.com/erodrigufer/raspall/internal/scraper"
	m "github.com/erodrigufer/raspall/internal/server/middlewares"
)

func (app *Application) routes() http.Handler {
	mws := m.NewMiddlewares(app.InfoLog, app.ErrorLog, app.sessionManager)

	globalMiddlewares := m.MiddlewareChain(app.sessionManager.LoadAndSave, mws.AddBuildHashCommitHeader, mws.LogRequest, mws.Cors, mws.RecoverPanic)

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
	protectedMux.Handle("POST /articles/nacio", app.undeliveredTemplate("Nació Digital", scraper.GetNacioArticles))
	protectedMux.Handle("POST /articles/hn", app.undeliveredTemplate("Hacker News", scraper.GetHackerNewsArticles))
	protectedMux.Handle("POST /articles/lobsters", app.undeliveredTemplate("Lobsters", scraper.GetLobstersArticles))
	protectedMux.Handle("POST /articles/theguardian", app.undeliveredTemplate("The Guardian", scraper.GetTheGuardianArticles))

	protectedMux.Handle("GET /articles/nacio/new", app.statusTemplate("Nació Digital", scraper.GetNacioArticles))
	protectedMux.Handle("GET /articles/hn/new", app.statusTemplate("Hacker News", scraper.GetHackerNewsArticles))
	protectedMux.Handle("GET /articles/lobsters/new", app.statusTemplate("Lobsters", scraper.GetLobstersArticles))
	protectedMux.Handle("GET /articles/theguardian/new", app.statusTemplate("The Guardian", scraper.GetTheGuardianArticles))

	return globalMiddlewares(mux)
}
