package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"lopa.to/sonimulus/api"
	"lopa.to/sonimulus/config"
	"lopa.to/sonimulus/controllers"
	"lopa.to/sonimulus/repository"
	"lopa.to/sonimulus/server"
)

func main() {
	// Load config struct from environment variables and program arguments
	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error("failed to initialize config: %v\n", "error", err)
		return
	}

	// Initialize database connection
	db, err := repository.NewDB(cfg)
	if err != nil {
		slog.Error("failed to initialize database: %v\n", "error", err)
		return
	}

	// Initialize database repositories
	usersRepository := repository.NewUsersRepository(db)

	// Initialize server
	authController := controllers.NewAuthController(usersRepository, cfg)
	usersController := controllers.NewUsersController(usersRepository)

	controller := server.NewServer(authController, usersController, cfg)
	handler := api.HandlerWithOptions(controller, api.StdHTTPServerOptions{
		Middlewares: []api.MiddlewareFunc{controller.AuthMiddleware},
	})
	server := http.Server{Addr: fmt.Sprintf(":%d", cfg.Port), Handler: handler}

	ctrlc := make(chan os.Signal, 1)
	signal.Notify(ctrlc, os.Interrupt, syscall.SIGTERM)
	go func() {
		// Wait for Ctrl-C signal
		<-ctrlc
		server.Close()
	}()

	// Start server
	slog.Info("Listening", "port", cfg.Port)
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		slog.Error("Server closed", "error", err)
	} else {
		slog.Info("Server closed", "error", err)
	}
}
