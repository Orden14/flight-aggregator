package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Orden14/flight-aggregator/src/service"
	"github.com/Orden14/flight-aggregator/src/util/sorter"
)

type FlightHandler struct {
	flightService service.FlightService
}

func NewFlightHandler(flightService service.FlightService) *FlightHandler {
	return &FlightHandler{flightService: flightService}
}

func (flightHandler *FlightHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()

	departureAirport := query.Get("from")
	arrivalAirport := query.Get("to")

	sortBy := sorter.NormalizeSortBy(query.Get("sort"))
	sortOrder := sorter.NormalizeOrder(query.Get("order"))

	flights, err := flightHandler.flightService.GetFlights(request.Context(), departureAirport, arrivalAirport, sortBy, sortOrder)

	if err != nil {
		http.Error(writer, "failed to fetch flights: "+err.Error(), http.StatusBadGateway)

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
