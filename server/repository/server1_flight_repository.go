package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Orden14/flight-aggregator/config"
	"github.com/Orden14/flight-aggregator/domain"
	"github.com/Orden14/flight-aggregator/model"
)

type Server1FlightRepository struct {
	baseURL string
	client  *http.Client
}

func NewServer1FlightRepository(config config.JSONServerConfig) *Server1FlightRepository {
	return &Server1FlightRepository{
		baseURL: config.BaseURL(),
		client:  &http.Client{Timeout: 0},
	}
}

func (flightRepository *Server1FlightRepository) Fetch(context context.Context) ([]domain.Flight, error) {
	url := fmt.Sprintf("%s/flights", flightRepository.baseURL)

	request, errors := http.NewRequestWithContext(context, http.MethodGet, url, nil)

	if errors != nil {
		return nil, fmt.Errorf("flight build request: %w", errors)
	}

	response, errors := flightRepository.client.Do(request)

	if errors != nil {
		return nil, fmt.Errorf("flight GET %s: %w", url, errors)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(response.Body, 1<<14))

		return nil, fmt.Errorf("flight status %d: %s", response.StatusCode, string(body))
	}

	var flightItems []model.Server1FlightItem

	if errors := json.NewDecoder(response.Body).Decode(&flightItems); errors != nil {
		return nil, fmt.Errorf("flight decode array: %w", errors)
	}

	flightsResponse := make([]domain.Flight, 0, len(flightItems))
	for _, flight := range flightItems {
		departureTime, errors := time.Parse(time.RFC3339, flight.DepartureTime)

		if errors != nil {
			return nil, fmt.Errorf("flight bad departureTime %q: %w", flight.DepartureTime, errors)
		}

		arrivalTime, errors := time.Parse(time.RFC3339, flight.ArrivalTime)

		if errors != nil {
			return nil, fmt.Errorf("flight bad arrivalTime %q: %w", flight.ArrivalTime, errors)
		}

		flightsResponse = append(flightsResponse, domain.Flight{
			Reference:     flight.BookingID,
			FlightNumber:  flight.FlightNumber,
			From:          flight.DepartureAirport,
			To:            flight.ArrivalAirport,
			DepartureTime: departureTime,
			ArrivalTime:   arrivalTime,
			Price:         flight.Price,
			Currency:      flight.Currency,
		})
	}

	return flightsResponse, nil
}
