package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Orden14/flight-aggregator/handler"
	"github.com/Orden14/flight-aggregator/httpserver"
)

func main() {
	health := handler.NewHealthHandler()
	router := httpserver.NewRouter(health)

	addr := ":3001"

	if v := os.Getenv("PORT"); v != "" {
		addr = ":" + v
	}

	log.Println("Flight Aggregator listening on", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}
