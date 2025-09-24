package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Orden14/flight-aggregator/src/config"
	"github.com/Orden14/flight-aggregator/src/domain"
	"github.com/Orden14/flight-aggregator/src/model"
)

type Server2FlightRepository struct {
	baseURL string
	client  *http.Client
}

func NewServer2FlightRepository(config config.JSONServerConfig) *Server2FlightRepository {
	return &Server2FlightRepository{
		baseURL: config.BaseURL(),
		client:  &http.Client{Timeout: 0},
	}
}

func (flightRepository *Server2FlightRepository) Fetch(ctx context.Context) ([]domain.Flight, error) {
	url := fmt.Sprintf("%s/flight_to_book", flightRepository.baseURL)

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

	var flightItems []model.Server2FlightItem

	if err := json.NewDecoder(response.Body).Decode(&flightItems); err != nil {
		return nil, fmt.Errorf("flight decode array: %w", err)
	}

	flightsResponse := make([]domain.Flight, 0, len(flightItems))

	for _, flight := range flightItems {
		if len(flight.Segments) == 0 {
			continue
		}

		firstSegment := flight.Segments[0].Flight
		lastSegment := flight.Segments[len(flight.Segments)-1].Flight

		departureTime, err := time.Parse(time.RFC3339, firstSegment.Depart)

		if err != nil {
			return nil, fmt.Errorf("flight bad depart %q: %w", firstSegment.Depart, err)
		}

		arrivalTime, err := time.Parse(time.RFC3339, lastSegment.Arrive)

		if err != nil {
			return nil, fmt.Errorf("flight bad arrive %q: %w", lastSegment.Arrive, err)
		}

		flightsResponse = append(flightsResponse, domain.Flight{
			Reference:     flight.Reference,
			FlightNumber:  firstSegment.Number,
			From:          firstSegment.From,
			To:            lastSegment.To,
			DepartureTime: departureTime,
			ArrivalTime:   arrivalTime,
			Price:         flight.Total.Amount,
			Currency:      flight.Total.Currency,
		})
	}

	return flightsResponse, nil
}
