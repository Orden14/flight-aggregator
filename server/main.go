package main

import (
	"log"
	"net/http"

	"github.com/Orden14/flight-aggregator/config"
	"github.com/Orden14/flight-aggregator/handler"
	"github.com/Orden14/flight-aggregator/httpserver"
	"github.com/Orden14/flight-aggregator/repository"
	"github.com/Orden14/flight-aggregator/service"
	"github.com/spf13/viper"
)

func main() {
	cfg, err := config.Load()

	if err != nil {
		log.Fatal("config error: ", err)
	}

	r1 := repository.NewServer1FlightRepository(cfg.JServer1)
	r2 := repository.NewServer2FlightRepository(cfg.JServer2)

	svc := service.NewFlightService(5, r1, r2)

	health := handler.NewHealthHandler()
	flight := handler.NewFlightHandler(svc)
	router := httpserver.NewRouter(health, flight)

	var addr string

	if v := viper.GetString("SERVER_PORT"); v != "" {
		addr = ":" + v
	} else {
		addr = ":3001"
	}

	log.Println("Flight Aggregator listening on", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}
