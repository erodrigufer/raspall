package server

import (
	"net/http"

	"github.com/erodrigufer/raspall/internal/scraper"
	m "github.com/erodrigufer/raspall/internal/server/middlewares"
	"github.com/erodrigufer/raspall/internal/static"
)

func (app *Application) defineEndpoints() (http.Handler, error) {
	mws := m.NewMiddlewares(app.InfoLog, app.ErrorLog, app.sessionManager, app.disableAuthentication)

	globalMiddlewares := m.MiddlewareChain(app.sessionManager.LoadAndSave, mws.AddBuildHashCommitHeader, mws.LogRequest, mws.RecoverPanic)

	fileServer := http.FileServer(http.FS(static.STATIC_CONTENT))

	mux := http.NewServeMux()
	mux.Handle("GET /login", mws.Authenticate(app.getLogin()))
	mux.Handle("POST /login", app.postLogin())
	mux.Handle("POST /logout", app.postLogout())
	mux.Handle("GET /content/", mws.PublicCacheCacheControl(fileServer))
	mux.Handle("GET /api/health", app.health())

	protectedMux := http.NewServeMux()
	mux.Handle("/", mws.Authenticate(mws.PrivateCacheControl(protectedMux)))

	protectedMux.Handle("GET /", app.dailyVisitingFrequency(app.index()))
	protectedMux.Handle("POST /articles/nacio", app.undeliveredTemplate("Nació Digital", scraper.GetNacioArticles, "nacio_settled"))
	protectedMux.Handle("POST /articles/hn", app.undeliveredTemplate("Hacker News", scraper.GetHackerNewsArticles, "hn_settled"))
	protectedMux.Handle("POST /articles/lobsters", app.undeliveredTemplate("Lobsters", scraper.GetLobstersArticles, "lobsters_settled"))
	protectedMux.Handle("POST /articles/theguardian", app.undeliveredTemplate("The Guardian", scraper.GetTheGuardianArticles, "the_guardian_settled"))

	protectedMux.Handle("GET /articles/nacio/new", app.statusTemplate("Nació Digital", scraper.GetNacioArticles))
	protectedMux.Handle("GET /articles/hn/new", app.statusTemplate("Hacker News", scraper.GetHackerNewsArticles))
	protectedMux.Handle("GET /articles/lobsters/new", app.statusTemplate("Lobsters", scraper.GetLobstersArticles))
	protectedMux.Handle("GET /articles/theguardian/new", app.statusTemplate("The Guardian", scraper.GetTheGuardianArticles))

	return globalMiddlewares(mux), nil
}
