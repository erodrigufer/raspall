package server

import (
	"fmt"
	"net/http"

	"github.com/erodrigufer/raspall/internal/scraper"
	m "github.com/erodrigufer/raspall/internal/server/middlewares"
	"github.com/erodrigufer/raspall/internal/utils"
)

func (app *Application) routes() http.Handler {
	mws := m.NewMiddlewares(app.InfoLog, app.ErrorLog)

	globalMiddlewares := m.MiddlewareChain(mws.LogRequest, mws.Cors, mws.RecoverPanic)

	mux := http.NewServeMux()
	mux.Handle("GET /v1/articles/{site...}", app.news())
	mux.Handle("GET /v1/health", app.health())

	return globalMiddlewares(mux)
}

func (app *Application) news() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		site := r.PathValue("site")

		options, err := getQueryOptions(r.URL.RawQuery)
		if err != nil {
			utils.HandleServerError(w, fmt.Errorf("error parsing URL query parameters: %w", err), app.ErrorLog)
			return
		}

		switch site {
		case "nació":
			{
				articles := scraper.GetNacioArticles(r.Context(), app.InfoLog, app.ErrorLog)
				articles = limit(options.limit, articles)
				unreadArticles, err := getUndeliveredObjects(articles, app.cache)
				if err != nil {
					utils.HandleServerError(w, err, app.ErrorLog)
					return
				}
				utils.SendJSONResponse(w, 200, unreadArticles)
				return
			}
		case "zeit":
			{
				articles := scraper.GetZeitArticles(r.Context(), app.InfoLog, app.ErrorLog, options.removePaywall)
				articles = limit(options.limit, articles)
				unreadArticles, err := getUndeliveredObjects(articles, app.cache)
				if err != nil {
					utils.HandleServerError(w, err, app.ErrorLog)
					return
				}
				utils.SendJSONResponse(w, 200, unreadArticles)
				return
			}
		case "hn":
			{
				articles := scraper.GetHackerNewsArticles(r.Context(), app.InfoLog, app.ErrorLog)
				articles = limit(options.limit, articles)
				unreadArticles, err := getUndeliveredObjects(articles, app.cache)
				if err != nil {
					utils.HandleServerError(w, err, app.ErrorLog)
					return
				}
				utils.SendJSONResponse(w, 200, unreadArticles)
				return
			}
		case "":
			{
				utils.SendJSONResponse(w, 200, "all news sites")
				return
			}
		default:
			{
				utils.HandleNotFoundError(w)
				return
			}
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
