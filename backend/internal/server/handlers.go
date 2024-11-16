package server

import (
	"fmt"
	"net/http"

	"github.com/erodrigufer/raspall/internal/hypermedia"
	"github.com/erodrigufer/raspall/internal/scraper"
	"github.com/erodrigufer/raspall/internal/utils"
	"github.com/erodrigufer/raspall/internal/views"
)

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
			utils.SendErrorMessage(w, http.StatusBadRequest, app.ErrorLog, fmt.Sprintf("Bad request: %s", err.Error()))
			return
		}

		username := r.PostForm.Get("username")
		password := r.PostForm.Get("password")

		if username != app.authorizedUsername || password != app.authorizedPassword {
			utils.SendErrorMessage(w, http.StatusUnauthorized, app.ErrorLog, `<p class="error-response"><b>Username</b> and/or <b>password</b> are invalid.</p>`)
			return
		}

		// Renew the session token before making the privilege-level change.
		err = app.sessionManager.RenewToken(r.Context())
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		// Make the privilege-level change.
		app.sessionManager.Put(r.Context(), "userID", username)

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (app *Application) postLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := app.sessionManager.Destroy(r.Context())
		if err != nil {
			utils.HandleServerError(w, fmt.Errorf("unable to destroy session: %w", err), app.ErrorLog)
			return
		}
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func (app *Application) index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := hypermedia.RenderComponent(r.Context(), w, views.Index())
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
		err = hypermedia.RenderComponent(r.Context(), w, views.ArticleViewer(unreadArticles, "Nació Digital"))
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
		err = hypermedia.RenderComponent(r.Context(), w, views.ArticleViewer(unreadArticles, "Lobsters"))
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
		err = hypermedia.RenderComponent(r.Context(), w, views.ArticleViewer(unreadArticles, "Hacker News"))
		if err != nil {
			utils.HandleServerError(w, fmt.Errorf("unable to render templ component: %w", err), app.ErrorLog)
		}
	}
}

func (app *Application) theGuardian() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		options, err := getQueryOptions(r.URL.RawQuery)
		if err != nil {
			utils.HandleServerError(w, fmt.Errorf("error parsing URL query parameters: %w", err), app.ErrorLog)
			return
		}
		articles := scraper.GetTheGuardianArticles(r.Context(), app.InfoLog, app.ErrorLog)
		articles = limit(options.limit, articles)
		unreadArticles, err := getUndeliveredObjects(articles, app.cache)
		if err != nil {
			utils.HandleServerError(w, err, app.ErrorLog)
			return
		}
		err = hypermedia.RenderComponent(r.Context(), w, views.ArticleViewer(unreadArticles, "The Guardian"))
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
