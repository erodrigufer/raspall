package server

import (
	"net/http"

	"github.com/erodrigufer/raspall/internal/scraper"
)

func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("GET /v1/news", app.news())
	mux.Handle("GET /v1/health", app.health())

	return mux
}

func (app *Application) news() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		articles := scraper.ScrapeZeit(r.Context(), app.InfoLog, app.ErrorLog)
		SendJSONResponse(w, 200, articles)
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
