package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

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
	// Cancel context if application receives a SIGNINT/SIGTERM signal, and
	// use cancelled context to start a graceful shutdown of the application.
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	app, err := server.NewAPI(ctx)
	if err != nil {
		return fmt.Errorf("unable to create a new app: %w", err)
	}
	app.StartServerWithGracefulShutdown(ctx)

	return nil
}
