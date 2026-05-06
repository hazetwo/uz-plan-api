package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"uz-plan-api/docs"
	"uz-plan-api/internal/database"
	"uz-plan-api/internal/schedule"

	"github.com/go-chi/chi/v5"
	"golang.org/x/time/rate"
)

func main() {
	ctx := context.Background()

	r := chi.NewRouter()

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

	limiter := rate.NewLimiter(rate.Limit(10), 20)

	scr := schedule.NewScraper()
	repo, rs := schedule.NewRedisRepository(rdb)
	svc := schedule.NewService(scr, repo, rs)
	handler := schedule.NewHandler(svc, limiter)

	env := os.Getenv("APP_ENV")
	if env != "production" {
		r.Get("/openapi.yaml", docs.Spec)
	}

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/fields", handler.GetFields)
		r.Get("/groups/{id}", handler.GetGroupsFromID)
		r.Get("/schedule/{id}", handler.GetScheduleFromID)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	slog.Info("Server listening", "addr", "http://0.0.0.0:"+port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		slog.Error("Server failed", "err", err)
		os.Exit(1)
	}
}
