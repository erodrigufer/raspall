package server

import (
	"fmt"
	"net/http"

	"github.com/erodrigufer/raspall/internal/scraper"
	"github.com/patrickmn/go-cache"
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
				SendJSONResponse(w, 200, articles)

			}
		case "zeit":
			{
				articles := scraper.GetZeitArticles(r.Context(), app.InfoLog, app.ErrorLog, options.removePaywall)
				articles = limit(options.limit, articles)
				SendJSONResponse(w, 200, articles)
			}
		case "hn":
			{
				articles := scraper.GetHackerNewsArticles(r.Context(), app.InfoLog, app.ErrorLog)
				articles = limit(options.limit, articles)
				unreadArticles := make([]scraper.Article, 0, 20)
				for _, article := range articles {
					// TODO: check error
					delivered, _ := checkDeliveryStatus(article, app.cache)
					if !delivered {
						unreadArticles = append(unreadArticles, article)
					}
				}
				SendJSONResponse(w, 200, unreadArticles)
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

func limit[O any](limit int, objects []O) []O {
	if limit < 1 {
		return objects
	}

	if len(objects) >= limit {
		return objects[:limit]
	}

	return objects
}

func checkDeliveryStatus(element Deliverable, c *cache.Cache) (bool, error) {
	hashId, err := element.CreateHash()
	if err != nil {
		return false, fmt.Errorf("unable to create a hash id: %w", err)
	}

	_, found := c.Get(hashId)
	if found {
		return true, nil
	}

	c.Set(hashId, true, cache.DefaultExpiration)
	return false, nil

}

type Deliverable interface {
	CreateHash() (string, error)
}
