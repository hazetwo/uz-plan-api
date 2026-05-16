package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"uz-plan-api/docs"
	"uz-plan-api/internal/database"
	"uz-plan-api/internal/ratelimiter"
	"uz-plan-api/internal/schedule"
	"uz-plan-api/internal/scraper"

	"github.com/go-chi/chi/v5"
)

func main() {
	ctx := context.Background()

	r := chi.NewRouter()

	origins := os.Getenv("ALLOWED_ORIGINS")
	if origins == "" {
		origins = "http://localhost:5173"
	}
	allowedOrigins := strings.Split(origins, ",")

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	}))

	rdb, err := database.Connect(ctx)
	if err != nil {
		slog.Error("Failed to connect to Redis", "err", err)
		os.Exit(1)
	}
	defer func() {
		err := rdb.Close()
		if err != nil {
			slog.Error("Failed to close Redis", "err", err)
		}
	}()

	slog.Info("Connected to Redis")

	rl := ratelimiter.New(
		ratelimiter.Options{
			RateLimit:          10,
			Bucket:             20,
			ViolationThreshold: 10,
			ViolationWindow:    time.Minute,
			BanDuration:        5 * time.Minute,
		},
	)

	scr := scraper.New()
	repo, rs := database.New(rdb)
	svc := schedule.NewService(scr, repo, rs)
	h := schedule.NewHandler(svc)

	r.Use(rl.LimitMiddleware)

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "dev"
	}

	switch env {
	case "prod", "production":
		// skip registering the docs route
	case "dev", "development":
		r.Get("/openapi.yaml", docs.Spec)
	default:
		slog.Warn("Unknown APP_ENV value, defaulting to development", "value", env)
		r.Get("/openapi.yaml", docs.Spec)
	}

	slog.Info("running environment", "env", env)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/fields", h.GetFields)
		r.Get("/groups/{id}", h.GetGroupsFromID)
		r.Get("/schedule/{id}", h.GetScheduleFromID)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	notify, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer stop()

	slog.Info("Server listening", "addr", "http://0.0.0.0:"+port)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server failed", "err", err)
			os.Exit(1)
		}
	}()
	<-notify.Done()

	shutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	slog.Info("Shutting down...")
	err = server.Shutdown(shutCtx)
	if err != nil {
		slog.Error("Shutdown failed", "err", err)
		os.Exit(1)
	}

}
