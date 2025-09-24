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

func (flightRepository *Server2FlightRepository) Fetch(context context.Context) ([]domain.Flight, error) {
	url := fmt.Sprintf("%s/flight_to_book", flightRepository.baseURL)

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

	var flightItems []model.Server2FlightItem

	if errors := json.NewDecoder(response.Body).Decode(&flightItems); errors != nil {
		return nil, fmt.Errorf("flight decode array: %w", errors)
	}

	flightsResponse := make([]domain.Flight, 0, len(flightItems))

	for _, flight := range flightItems {
		if len(flight.Segments) == 0 {
			continue
		}

		firstSegment := flight.Segments[0].Flight
		lastSegment := flight.Segments[len(flight.Segments)-1].Flight

		departureTime, errors := time.Parse(time.RFC3339, firstSegment.Depart)

		if errors != nil {
			return nil, fmt.Errorf("flight bad depart %q: %w", firstSegment.Depart, errors)
		}

		arrivalTime, errors := time.Parse(time.RFC3339, lastSegment.Arrive)

		if errors != nil {
			return nil, fmt.Errorf("flight bad arrive %q: %w", lastSegment.Arrive, errors)
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
