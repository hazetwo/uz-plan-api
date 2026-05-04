package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"uz-plan-api/internal/database"
	"uz-plan-api/internal/schedule"

	"github.com/go-chi/chi/v5"
	"golang.org/x/time/rate"
)

//var supported = []string{"401"}

func main() {
	ctx := context.Background()

	r := chi.NewRouter()

	rdb, err := database.Connect(ctx)
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	defer func() {
		err := rdb.Close()
		if err != nil {
			log.Printf("Failed to close Redis: %v", err)
		}
	}()

	fmt.Println("Connected to Redis")

	var port = "8080"

	limiter := rate.NewLimiter(rate.Limit(10), 20)

	scr := schedule.NewScraper()
	repo, rs := schedule.NewRedisRepository(rdb)
	svc := schedule.NewService(scr, repo, rs)
	handler := schedule.NewHandler(svc, limiter)

	r.Get("/api/fields", handler.GetFields)
	r.Get("/api/groups/{id}", handler.GetGroupsFromID)
	r.Get("/api/schedule/{id}", handler.GetScheduleFromID)

	addr := ":" + port
	fmt.Printf("Server started at http://localhost:%s\n", port)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
