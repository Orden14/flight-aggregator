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

func (flightRepository *Server1FlightRepository) Fetch(ctx context.Context) ([]domain.Flight, error) {
	url := fmt.Sprintf("%s/flights", flightRepository.baseURL)

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return nil, fmt.Errorf("flight build request: %w", err)
	}

	response, err := flightRepository.client.Do(request)

	if err != nil {
		return nil, fmt.Errorf("flight GET %s: %w", url, err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(response.Body, 1<<14))

		return nil, fmt.Errorf("flight status %d: %s", response.StatusCode, string(body))
	}

	var flightItems []model.Server1FlightItem

	if err := json.NewDecoder(response.Body).Decode(&flightItems); err != nil {
		return nil, fmt.Errorf("flight decode array: %w", err)
	}

	flightsResponse := make([]domain.Flight, 0, len(flightItems))
	for _, flight := range flightItems {
		departureTime, err := time.Parse(time.RFC3339, flight.DepartureTime)

		if err != nil {
			return nil, fmt.Errorf("flight bad departureTime %q: %w", flight.DepartureTime, err)
		}

		arrivalTime, err := time.Parse(time.RFC3339, flight.ArrivalTime)

		if err != nil {
			return nil, fmt.Errorf("flight bad arrivalTime %q: %w", flight.ArrivalTime, err)
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
