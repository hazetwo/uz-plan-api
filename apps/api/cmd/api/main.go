package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"uz-plan-api/internal/database"

	"github.com/go-chi/chi/v5"
)

var supported = []string{"401"}
var uFields = "https://plan.uz.zgora.pl/grupy_lista_kierunkow.php"
var uGroups = "https://plan.uz.zgora.pl/grupy_lista_grup_kierunku.php"
var uSchedule = "https://plan.uz.zgora.pl/grupy_plan.php"

func main() {
	ctx := context.Background()

	r := chi.NewRouter()

	db, err := database.Connect(ctx)
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Printf("Failed to close Redis: %v", err)
		}
	}()

	fmt.Println("Connected to Redis")

	var port = "8080"

	//scr := schedule.NewScraper("plan.uz.zgora.pl")
	//svc := schedule.NewService(scr)

	addr := ":" + port
	fmt.Printf("Server started at http://localhost:%s\n", port)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
