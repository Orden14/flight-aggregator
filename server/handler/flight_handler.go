package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Orden14/flight-aggregator/service"
	"github.com/Orden14/flight-aggregator/sorter"
)

type FlightHandler struct {
	flightService service.FlightService
}

func NewFlightHandler(flightService service.FlightService) *FlightHandler {
	return &FlightHandler{flightService: flightService}
}

func (flightHandler *FlightHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()

	from := query.Get("from")
	to := query.Get("to")

	sortBy := sorter.NormalizeSortBy(query.Get("sort"))
	sortOrder := sorter.NormalizeOrder(query.Get("sortOrder"))

	flights, errors := flightHandler.flightService.GetFlights(request.Context(), from, to, sortBy, sortOrder)

	if errors != nil {
		http.Error(writer, "failed to fetch flights: "+errors.Error(), http.StatusBadGateway)

		return
	}

	writer.Header().Set("Content-Type", "application/json")

	json.NewEncoder(writer).Encode(map[string]any{
		"flights_count": len(flights),
		"sort_by":       sortBy,
		"sort_order":    sortOrder,
		"items":         flights,
	})
}
