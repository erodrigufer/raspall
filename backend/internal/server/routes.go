package server

import (
	"fmt"
	"net/http"

	"github.com/erodrigufer/raspall/internal/scraper"
)

func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /v1/news/{site...}", app.news())
	mux.Handle("GET /v1/health", app.health())

	return mux
}

func (app *Application) news() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		site := r.PathValue("site")

		options, err := getQueryOptions(r.URL.RawQuery)
		if err != nil {
			HandleServerError(w, fmt.Errorf("error parsing URL query parameters: %w", err), app.ErrorLog)
			return
		}

		switch site {
		case "nacio":
			{
				articles := scraper.GetNacioArticles(r.Context(), app.InfoLog, app.ErrorLog)
				articles = limit(options.limit, articles)
				unreadArticles, err := getUndeliveredObjects(articles, app.cache)
				if err != nil {
					HandleServerError(w, err, app.ErrorLog)
					return
				}
				SendJSONResponse(w, 200, unreadArticles)
				return
			}
		case "zeit":
			{
				articles := scraper.GetZeitArticles(r.Context(), app.InfoLog, app.ErrorLog, options.removePaywall)
				articles = limit(options.limit, articles)
				unreadArticles, err := getUndeliveredObjects(articles, app.cache)
				if err != nil {
					HandleServerError(w, err, app.ErrorLog)
					return
				}
				SendJSONResponse(w, 200, unreadArticles)
				return
			}
		case "hn":
			{
				articles := scraper.GetHackerNewsArticles(r.Context(), app.InfoLog, app.ErrorLog)
				articles = limit(options.limit, articles)
				unreadArticles, err := getUndeliveredObjects(articles, app.cache)
				if err != nil {
					HandleServerError(w, err, app.ErrorLog)
					return
				}
				SendJSONResponse(w, 200, unreadArticles)
				return
			}
		case "":
			{
				SendJSONResponse(w, 200, "all news sites")
				return
			}
		default:
			{
				HandleNotFoundError(w)
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
		err := SendJSONResponse(w, http.StatusOK, response)
		if err != nil {
			HandleServerError(w, err, app.ErrorLog)
		}
	}
}
