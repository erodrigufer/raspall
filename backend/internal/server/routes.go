package server

import (
	"net/http"
	"net/url"

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

		values, err := url.ParseQuery(r.URL.RawQuery)
		// TODO: Handle URL query parsing error.
		if err != nil {
		}
		rp := values.Get("removePaywall")
		// Default is to remove paywalled articles.
		removePaywall := true
		if rp == "false" {
			removePaywall = false
		}

		switch site {
		case "nacio":
			{
				articles := scraper.ScrapeNacioDigital(r.Context(), app.InfoLog, app.ErrorLog)
				SendJSONResponse(w, 200, articles)

			}
		case "zeit":
			{
				articles := scraper.GetZeitArticles(r.Context(), app.InfoLog, app.ErrorLog, removePaywall)
				SendJSONResponse(w, 200, articles)
			}
		default:
			SendJSONResponse(w, 200, "all news sites")
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
