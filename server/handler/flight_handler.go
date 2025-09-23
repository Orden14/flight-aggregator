package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Orden14/flight-aggregator/service"
)

type FlightHandler struct {
	svc service.FlightService
}

func NewFlightHandler(svc service.FlightService) *FlightHandler {
	return &FlightHandler{svc: svc}
}

func (h *FlightHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	flights, err := h.svc.GetFlights()

	if err != nil {
		http.Error(w, "failed to fetch flights: "+err.Error(), http.StatusBadGateway)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]any{
		"count": len(flights),
		"items": flights,
	})
}
