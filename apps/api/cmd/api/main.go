package main

import (
	"fmt"
	"uz-plan-api/internal/scraper"
)

var supported = []string{"401"}
var uFields = "https://plan.uz.zgora.pl/grupy_lista_kierunkow.php"
var uGroups = "https://plan.uz.zgora.pl/grupy_lista_grup_kierunku.php"
var uSchedule = "https://plan.uz.zgora.pl/grupy_plan.php"

func main() {
	//r := chi.NewRouter()
	//
	//var port = "8080"
	//
	//addr := ":" + port
	//fmt.Printf("Server started at http://localhost:%s\n", port)
	//if err := http.ListenAndServe(addr, r); err != nil {
	//	log.Fatal(err)
	//}

	s := scraper.New("plan.uz.zgora.pl")
	fields := s.GetIdsOfFields(uFields)
	groups := s.GetIdsOfGroups(uGroups, fields, supported)
	fmt.Printf("%v", groups)
}
