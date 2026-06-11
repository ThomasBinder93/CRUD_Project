package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"CRUD_Project/config"
	"CRUD_Project/database"
	"CRUD_Project/handlers"
	"CRUD_Project/middleware"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Setup logger
	logger := setupLogger(cfg.Logger.Level)

	logger.Info("starting application",
		slog.String("version", cfg.App.Version),
		slog.String("environment", cfg.App.Environment),
	)

	// Initialize database
	repo, err := database.New(cfg.GetDSN())
	if err != nil {
		logger.Error("failed to initialize database", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer func() {
		if err := repo.Close(); err != nil {
			logger.Error("failed to close database", slog.String("error", err.Error()))
		}
	}()

	logger.Info("database initialized successfully")

	// Setup HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      setupRouter(repo, logger, cfg),
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}

	// Graceful shutdown channel
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		logger.Info("server started", slog.String("address", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", slog.String("error", err.Error()))
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	logger.Info("shutdown signal received, gracefully shutting down...")

	// Graceful shutdown with 30 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server shutdown error", slog.String("error", err.Error()))
		os.Exit(1)
	}

	logger.Info("server shut down successfully")
}

// setupLogger creates and configures the structured logger
func setupLogger(level string) *slog.Logger {
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}
	handler := slog.NewTextHandler(os.Stdout, opts)
	return slog.New(handler)
}

// setupRouter configures all routes and middleware
func setupRouter(repo *database.Repository, logger *slog.Logger, cfg *config.Config) http.Handler {
	mux := http.NewServeMux()

	// Create handlers
	itemHandler := handlers.NewItemHandler(repo, logger)

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status":"healthy"}`)
	})

	// API routes - with JSON content-type and logging
	apiHandler := http.NewServeMux()
	apiHandler.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if r.Method == http.MethodGet {
			itemHandler.GetAll(w, r)
		} else if r.Method == http.MethodPost {
			itemHandler.Create(w, r)
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	apiHandler.HandleFunc("/items/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		id := r.URL.Path[len("/items/"):]
		if id == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Wrap with context middleware to pass ID
		wrappedHandler := handlers.SetIDInContext(id)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				itemHandler.GetByID(w, r)
			} else if r.Method == http.MethodPut {
				itemHandler.Update(w, r)
			} else if r.Method == http.MethodDelete {
				itemHandler.Delete(w, r)
			} else {
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		}))
		wrappedHandler.ServeHTTP(w, r)
	})

	// Wrap API routes with middleware (JSON content-type and logging)
	apiWithMiddleware := middleware.CORSMiddleware()(
		middleware.LoggingMiddleware(logger)(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				apiHandler.ServeHTTP(w, r)
			}),
		),
	)

	mux.Handle("/api/", http.StripPrefix("/api", apiWithMiddleware))

	// Static files (Frontend) - served WITHOUT middleware
	staticFS := http.FileServer(http.Dir("./static/"))
	mux.Handle("/", staticFS)

	// Wrap everything with recovery and CORS middleware
	return middleware.RecoveryMiddleware(logger)(
		middleware.CORSMiddleware()(mux),
	)
}
