package server

import (
	"context"
	"fmt"
	"log"
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
	ErrorLog *log.Logger
	// InfoLog informative server logger.
	InfoLog *log.Logger
	// cache is a key-value store to manage
	// the news that have already been delivered
	// to the user.
	cache              *cache.Cache
	sessionManager     *scs.SessionManager
	authorizedUsername string
	authorizedPassword string
}

func NewAPI(ctx context.Context) (*Application, error) {
	app := new(Application)

	app.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.ErrorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	var ok bool
	app.authorizedUsername, ok = os.LookupEnv("AUTH_USERNAME")
	if !ok {
		return nil, fmt.Errorf("AUTH_USERNAME env var is missing")
	}
	app.authorizedPassword, ok = os.LookupEnv("AUTH_PASSWORD")
	if !ok {
		return nil, fmt.Errorf("AUTH_PASSWORD env var is missing")
	}

	app.sessionManager = scs.New()

	app.srv = &http.Server{
		Addr:     ":80",
		ErrorLog: app.ErrorLog,
		Handler:  app.routes(),
		// Time after which inactive keep-alive connections will be closed.
		IdleTimeout: time.Minute,
		// Max. time to read the header and body of a request in the server.
		ReadTimeout: 30 * time.Second,
		// Close connection if data is still being written after this time since
		// accepting the connection.
		WriteTimeout: 30 * time.Second,
	}

	c := cache.New(48*time.Hour, 12*time.Hour)
	app.cache = c

	return app, nil
}

// StartServerWithGracefulShutdown starts a server and gracefully handles shutdowns.
// If the server receives an os.Interrupt signal the backend knows that it should
// start the process of gracefully shutting down, i.e. closing DB connections and
// closing client connections.
func (app *Application) StartServerWithGracefulShutdown(ctx context.Context) {
	go func() {
		app.InfoLog.Printf("Starting raspall server at %s", app.srv.Addr)

		// Run server.
		if err := app.srv.ListenAndServe(); err != nil {
			// Error returned when server is closed, not actually an error, log to
			// info log.
			if err == http.ErrServerClosed {
				app.InfoLog.Print(err)
				// An actual error happened, log to error log.
			} else {
				app.ErrorLog.Printf("an error occured while executing ListenAndServe(): %v", err)
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
			app.ErrorLog.Printf("Server is not shutting down! Reason: %s", err.Error())
			// An error happened while gracefully shutting down, close abruptly.
			app.srv.Close()
		}
	}()

	// Wait on all goroutines performing asynchronous shutdowns before returning.
	wg.Wait()
}
