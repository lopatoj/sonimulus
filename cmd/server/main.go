package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"lopa.to/sonimulus/api"
	"lopa.to/sonimulus/controllers"
	"lopa.to/sonimulus/env"
	"lopa.to/sonimulus/handlers"
	"lopa.to/sonimulus/repository"
)

func main() {
	// Load config struct from environment variables and program arguments
	e, err := env.NewEnv()
	if err != nil {
		slog.Error("failed to initialize config", "error", err)
		return
	}

	// Initialize database connection
	db, err := repository.NewDB(e)
	if err != nil {
		slog.Error("failed to initialize database: %v\n", "error", err)
		return
	}

	// Initialize database repositories
	usersRepository := repository.NewUsersRepository(db)

	// Initialize server
	authController := controllers.NewAuthController(usersRepository, e)
	usersController := controllers.NewUsersController(usersRepository)

	handler := handlers.NewHandler(authController, usersController, e)
	apiHandler := api.HandlerWithOptions(handler, api.StdHTTPServerOptions{
		Middlewares: []api.MiddlewareFunc{handler.AuthMiddleware},
	})
	server := http.Server{Addr: fmt.Sprintf(":%d", e.Server.Port), Handler: apiHandler}

	ctrlc := make(chan os.Signal, 1)
	signal.Notify(ctrlc, os.Interrupt, syscall.SIGTERM)
	go func() {
		// Wait for Ctrl-C signal
		<-ctrlc
		server.Close()
	}()

	// Start server
	slog.Info("Listening", "port", e.Server.Port)
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		slog.Error("Server closed", "error", err)
	} else {
		slog.Info("Server closed", "error", err)
	}
}
