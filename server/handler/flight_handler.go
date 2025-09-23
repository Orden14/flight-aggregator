package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Orden14/flight-aggregator/service"
	"github.com/Orden14/flight-aggregator/sorter"
)

type FlightHandler struct {
	svc service.FlightService
}

func NewFlightHandler(svc service.FlightService) *FlightHandler {
	return &FlightHandler{svc: svc}
}

func (h *FlightHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	from := q.Get("from")
	to := q.Get("to")

	by := sorter.NormalizeSortBy(q.Get("sort"))
	order := sorter.NormalizeOrder(q.Get("order"))

	flights, err := h.svc.GetFlights(r.Context(), from, to, by, order)

	if err != nil {
		http.Error(w, "failed to fetch flights: "+err.Error(), http.StatusBadGateway)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]any{
		"count":   len(flights),
		"sort_by": by,
		"order":   order,
		"items":   flights,
	})
}
