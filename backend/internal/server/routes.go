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

	fileServer := http.StripPrefix("/static", http.FileServer(http.Dir("./static")))

	mux := http.NewServeMux()
	mux.Handle("GET /", app.index())
	mux.Handle("GET /static/", fileServer)
	mux.Handle("GET /login", app.getLogin())
	mux.Handle("POST /login", app.postLogin())
	mux.Handle("POST /nacio", mws.CacheControl(app.nacio()))
	mux.Handle("POST /hn", mws.CacheControl(app.hn()))
	mux.Handle("POST /lobsters", mws.CacheControl(app.lobsters()))

	mux.Handle("GET /api/health", app.health())

	return globalMiddlewares(mux)
}

func (app *Application) getLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := hypermedia.RenderComponent(r.Context(), w, views.Login())
		if err != nil {
			utils.HandleServerError(w, fmt.Errorf("unable to render templ component: %w", err), app.ErrorLog)
		}
	}
}

func (app *Application) postLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			// TODO: Create function to handle invalid inputs, bad request!
			return
		}

		username := r.PostForm.Get("username")
		password := r.PostForm.Get("password")

		fmt.Println(username, password)

		// TODO: continue here
		// err := hypermedia.RenderComponent(r.Context(), w, views.Login())
		// if err != nil {
		// 	utils.HandleServerError(w, fmt.Errorf("unable to render templ component: %w", err), app.ErrorLog)
		// }
	}
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
