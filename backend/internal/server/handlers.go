package server

import (
	"fmt"
	"net/http"
	"time"

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
		err := hypermedia.RenderComponent(r.Context(), w, views.Index(app.df.currentlyBlocked, app.formatUnblockTime()))
		if err != nil {
			utils.HandleServerError(w, fmt.Errorf("unable to render templ component: %w", err), app.ErrorLog)
		}
	}
}

func (app *Application) formatUnblockTime() string {
	unblockedTime := app.df.lastVisit.Add(app.df.allowedFrequency)
	unblockedTimeFormatted := unblockedTime.Format("15:04 MST -- 02 Jan 2006")
	return unblockedTimeFormatted
}

func (app *Application) undeliveredTemplate(sourceName string, scraperFunc scraper.ScraperFunc, eventName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		options, err := getQueryOptions(r.URL.RawQuery)
		if err != nil {
			utils.HandleServerError(w, fmt.Errorf("error parsing URL query parameters: %w", err), app.ErrorLog)
			return
		}
		var unreadArticles []scraper.Article
		cached := app.checkIfResponseCached(sourceName)
		if !cached {
			articles := scraperFunc(r.Context(), app.InfoLog, app.ErrorLog)
			articles = limit(options.limit, articles)
			unreadArticles, err = getUndeliveredObjects(articles, app.cache)
			if err != nil {
				utils.HandleServerError(w, err, app.ErrorLog)
				return
			}
			app.cacheResponse(sourceName)
		}
		// Triggers an event after the swaps of this response settled, so that the status of
		// the different sources gets triggered.
		w.Header().Add("HX-Trigger-After-Settle", eventName)
		err = hypermedia.RenderComponent(r.Context(), w, views.ArticleViewer(unreadArticles, sourceName))
		if err != nil {
			utils.HandleServerError(w, fmt.Errorf("unable to render templ component: %w", err), app.ErrorLog)
		}
	}
}

func (app *Application) statusTemplate(sourceName string, scraperFunc scraper.ScraperFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		options, err := getQueryOptions(r.URL.RawQuery)
		if err != nil {
			utils.HandleServerError(w, fmt.Errorf("error parsing URL query parameters: %w", err), app.ErrorLog)
			return
		}
		var unreadArticlesPresent bool
		cached := app.checkIfResponseCached(sourceName)
		if !cached {
			articles := scraperFunc(r.Context(), app.InfoLog, app.ErrorLog)
			articles = limit(options.limit, articles)
			unreadArticlesPresent, err = checkIfUndeliveredObjectsPresent(articles, app.cache)
			if err != nil {
				utils.HandleServerError(w, err, app.ErrorLog)
				return
			}
		}
		err = hypermedia.RenderComponent(r.Context(), w, views.UnreadArticlesNotifier(unreadArticlesPresent))
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

type dailyFrequency struct {
	lastVisit        time.Time
	currentlyBlocked bool
	allowedFrequency time.Duration
}

func (app *Application) dailyVisitingFrequency(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		lastVisit := app.df.lastVisit
		elapsedTime := lastVisit.Add(app.df.allowedFrequency)
		if elapsedTime.After(now) {
			app.df.currentlyBlocked = true
		} else {
			app.df.currentlyBlocked = false
			app.df.lastVisit = now
		}
		next.ServeHTTP(w, r)
	})
}
