package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/patrickmn/go-cache"
)

type Application struct {
	srv *http.Server
	// ErrorLog logs server errors.
	ErrorLog *slog.Logger
	// InfoLog informative server logger.
	InfoLog *slog.Logger
	// cache is a key-value store to manage
	// the news that have already been delivered
	// to the user.
	cache                 *cache.Cache
	sessionManager        *scs.SessionManager
	disableAuthentication bool
	authorizedUsername    string
	authorizedPassword    string
	df                    dailyFrequency
}

func NewAPI(ctx context.Context) (*Application, error) {
	app := new(Application)

	app.InfoLog = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	app.ErrorLog = slog.New(slog.NewJSONHandler(os.Stderr, nil))

	var ok bool
	app.authorizedUsername, ok = os.LookupEnv("AUTH_USERNAME")
	if !ok {
		return nil, fmt.Errorf("AUTH_USERNAME env var is missing")
	}
	app.authorizedPassword, ok = os.LookupEnv("AUTH_PASSWORD")
	if !ok {
		return nil, fmt.Errorf("AUTH_PASSWORD env var is missing")
	}

	port, ok := os.LookupEnv("PORT")
	if !ok {
		return nil, fmt.Errorf("PORT env var is missing")
	}

	disableAuthStr, ok := os.LookupEnv("DISABLE_AUTH")
	if !ok {
		app.disableAuthentication = false
	}
	if disableAuthStr == "true" {
		app.disableAuthentication = true
	} else {
		app.disableAuthentication = false
	}
	app.InfoLog.Info("configuring authentication environment", slog.Bool("DISABLE_AUTH", app.disableAuthentication))

	app.sessionManager = scs.New()
	app.sessionManager.Lifetime = 15 * 24 * time.Hour
	app.sessionManager.IdleTimeout = 15 * 24 * time.Hour

	allowedVisitFreq, err := time.ParseDuration("5h")
	if err != nil {
		return nil, fmt.Errorf("unable to parse allowedVisitFreq: %w", err)
	}
	app.df.allowedFrequency = allowedVisitFreq

	// http.Server can only handle loggers from the old log package.
	compatibleLogger := slog.NewLogLogger(slog.NewJSONHandler(os.Stderr, nil), slog.LevelError)

	endpoints, err := app.defineEndpoints()
	if err != nil {
		return nil, fmt.Errorf("unable to define endpoints: %w", err)
	}

	app.srv = &http.Server{
		Addr:     port,
		ErrorLog: compatibleLogger,
		Handler:  endpoints,
		// Time after which inactive keep-alive connections will be closed.
		IdleTimeout: time.Minute,
		// Max. time to read the header and body of a request in the server.
		ReadTimeout: 30 * time.Second,
		// Close connection if data is still being written after this time since
		// accepting the connection.
		WriteTimeout: 30 * time.Second,
	}

	c := cache.New(4*24*time.Hour, 12*time.Hour)
	app.cache = c

	return app, nil
}

// StartServerWithGracefulShutdown starts a server and gracefully handles shutdowns.
// If the server receives an os.Interrupt signal the backend knows that it should
// start the process of gracefully shutting down, i.e. closing DB connections and
// closing client connections.
func (app *Application) StartServerWithGracefulShutdown(ctx context.Context) {
	go func() {
		app.InfoLog.Info("starting raspall server", slog.String("server_address", app.srv.Addr))

		// Run server.
		if err := app.srv.ListenAndServe(); err != nil {
			// Error returned when server is closed, not actually an error, log to
			// info log.
			if err == http.ErrServerClosed {
				app.InfoLog.Info("server closed")
				// An actual error happened, log to error log.
			} else {
				app.ErrorLog.Error("an error occured while executing ListenAndServe()", slog.String("error_message", err.Error()))
			}
		}
	}()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		// Wait on each step of a gracious shutdown.
		defer wg.Done()
		// When ctx passed from main function gets cancelled with os.Interrupt signal
		// (ctx.Done() returns), this goroutine performs a shutdown.
		<-ctx.Done()

		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 15*time.Second)
		defer cancel()
		// Received an interrupt signal, shutdown.
		if err := app.srv.Shutdown(shutdownCtx); err != nil {
			// Error from closing listeners, or context timeout.
			app.ErrorLog.Error("server is not shutting down", slog.String("error_message", err.Error()))
			// An error happened while gracefully shutting down, close abruptly.
			app.srv.Close()
		}
	}()

	// Wait on all goroutines performing asynchronous shutdowns before returning.
	wg.Wait()
}
