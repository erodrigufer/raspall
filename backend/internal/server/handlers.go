package server

import (
	"fmt"
	"net/http"
	"os"

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
			// TODO: Create function to handle invalid inputs, bad request!
			return
		}

		username := r.PostForm.Get("username")
		password := r.PostForm.Get("password")

		// TODO: remove
		fmt.Println(username, password)

		// TODO: Check that envs are not ""
		AUTHORIZED_USERNAME := os.Getenv("AUTH_USERNAME")
		AUTHORIZED_PASSWORD := os.Getenv("AUTH_PASSWORD")
		if username != AUTHORIZED_USERNAME || password != AUTHORIZED_PASSWORD {
			// TODO: display error message
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Renew the session token before making the privilege-level change.
		err = app.sessionManager.RenewToken(r.Context())
		if err != nil {
			// TODO:
			http.Error(w, err.Error(), 500)
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
		userID := app.sessionManager.GetString(r.Context(), "userID")
		app.InfoLog.Printf("userID: %s", userID)

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
