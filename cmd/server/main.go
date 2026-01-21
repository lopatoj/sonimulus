package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"lopa.to/sonimulus/api/v1"
	"lopa.to/sonimulus/env"
	"lopa.to/sonimulus/handlers"
	"lopa.to/sonimulus/internal/auth"
	"lopa.to/sonimulus/internal/data"
	"lopa.to/sonimulus/internal/repo"
)

func main() {
	// Load config struct from environment variables and program arguments
	e, err := env.NewEnv()
	if err != nil {
		slog.Error("failed to initialize config", "error", err)
		return
	}

	// Initialize PostgreSQL connection
	pgdb, err := data.NewPostgresDB(e.DB.PostgresURI)
	if err != nil {
		slog.Error("failed to initialize database", "error", err)
		return
	}

	// Initialize Redis connection
	rdb, err := data.NewRedisClient(e.DB.RedisURI)
	if err != nil {
		slog.Error("failed to initialize redis client", "error", err)
		return
	}

	scc, err := data.NewSoundCloudClient(e.Soundcloud.APIURL)

	// Initialize data repositories
	stateRepo := repo.NewStateRepository(rdb)
	sessionRepo := repo.NewSessionRepository(rdb)
	soundCloudRepo := repo.NewSoundCloudRepository(scc, e)
	usersRepo := repo.NewUsersRepository(pgdb)

	// Initialize server
	authController := auth.NewAuthController(
		e,
		stateRepo,
		sessionRepo,
		soundCloudRepo,
		usersRepo,
	)

	baseHandler := handlers.NewHandler(authController, e)
	apiHandler := api.HandlerWithOptions(baseHandler, api.StdHTTPServerOptions{
		BaseURL:     e.Server.Route,
		Middlewares: []api.MiddlewareFunc{handlers.CorsMiddleware},
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
