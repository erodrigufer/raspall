package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/erodrigufer/raspall/internal/server"
)

func main() {

	ctx := context.Background()

	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

// run encapsulates the web application.
func run(ctx context.Context) error {
	// Cancel context if application receives a SIGNINT signal, and use cancelled
	// context to start a graceful shutdown of the application.
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	app := server.NewAPI(ctx)
	app.StartServerWithGracefulShutdown(ctx)

	return nil

}
