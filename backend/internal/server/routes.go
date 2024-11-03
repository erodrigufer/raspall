package server

import (
	"fmt"
	"net/http"

	"github.com/erodrigufer/raspall/internal/hypermedia"
	"github.com/erodrigufer/raspall/internal/scraper"
	m "github.com/erodrigufer/raspall/internal/server/middlewares"
	"github.com/erodrigufer/raspall/internal/utils"
	"github.com/erodrigufer/raspall/internal/views"
)

func (app *Application) routes() http.Handler {
	mws := m.NewMiddlewares(app.InfoLog, app.ErrorLog)

	globalMiddlewares := m.MiddlewareChain(mws.LogRequest, mws.Cors, mws.RecoverPanic)

	mux := http.NewServeMux()
	mux.Handle("GET /v1/health", app.health())
	mux.Handle("GET /", app.index())
	mux.Handle("POST /nacio", mws.CacheControl(app.nacio()))
	mux.Handle("POST /hn", mws.CacheControl(app.hn()))
	mux.Handle("POST /lobsters", mws.CacheControl(app.lobsters()))

	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	return globalMiddlewares(mux)
}

func (app *Application) index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := hypermedia.RenderComponent(r.Context(), w, views.Home())
		if err != nil {
			utils.HandleServerError(w, fmt.Errorf("unable to render templ component: %w", err), app.ErrorLog)
		}
	}
}

func (app *Application) nacio() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		options, err := getQueryOptions(r.URL.RawQuery)
		if err != nil {
			utils.HandleServerError(w, fmt.Errorf("error parsing URL query parameters: %w", err), app.ErrorLog)
			return
		}
		articles := scraper.GetNacioArticles(r.Context(), app.InfoLog, app.ErrorLog)
		articles = limit(options.limit, articles)
		unreadArticles, err := getUndeliveredObjects(articles, app.cache)
		if err != nil {
			utils.HandleServerError(w, err, app.ErrorLog)
			return
		}
		err = hypermedia.RenderComponent(r.Context(), w, views.ArticleViewer(unreadArticles))
		if err != nil {
			utils.HandleServerError(w, fmt.Errorf("unable to render templ component: %w", err), app.ErrorLog)
		}
	}
}

func (app *Application) lobsters() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		options, err := getQueryOptions(r.URL.RawQuery)
		if err != nil {
			utils.HandleServerError(w, fmt.Errorf("error parsing URL query parameters: %w", err), app.ErrorLog)
			return
		}
		articles := scraper.GetLobstersArticles(r.Context(), app.InfoLog, app.ErrorLog)
		articles = limit(options.limit, articles)
		unreadArticles, err := getUndeliveredObjects(articles, app.cache)
		if err != nil {
			utils.HandleServerError(w, err, app.ErrorLog)
			return
		}
		err = hypermedia.RenderComponent(r.Context(), w, views.ArticleViewer(unreadArticles))
		if err != nil {
			utils.HandleServerError(w, fmt.Errorf("unable to render templ component: %w", err), app.ErrorLog)
		}
	}
}

func (app *Application) hn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		options, err := getQueryOptions(r.URL.RawQuery)
		if err != nil {
			utils.HandleServerError(w, fmt.Errorf("error parsing URL query parameters: %w", err), app.ErrorLog)
			return
		}
		articles := scraper.GetHackerNewsArticles(r.Context(), app.InfoLog, app.ErrorLog)
		articles = limit(options.limit, articles)
		unreadArticles, err := getUndeliveredObjects(articles, app.cache)
		if err != nil {
			utils.HandleServerError(w, err, app.ErrorLog)
			return
		}
		err = hypermedia.RenderComponent(r.Context(), w, views.ArticleViewer(unreadArticles))
		if err != nil {
			utils.HandleServerError(w, fmt.Errorf("unable to render templ component: %w", err), app.ErrorLog)
		}
	}
}

func (app *Application) health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		response := map[string]string{
			"Status": "ok",
		}
		err := utils.SendJSONResponse(w, http.StatusOK, response)
		if err != nil {
			utils.HandleServerError(w, err, app.ErrorLog)
		}
	}
}
