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

type FlightRepository struct {
	baseURL string
	client  *http.Client
}

func NewFlightRepository(c config.JSONServerConfig) *FlightRepository {
	return &FlightRepository{
		baseURL: c.BaseURL(),
		client:  &http.Client{Timeout: 0},
	}
}

func (r *FlightRepository) Fetch(ctx context.Context) ([]domain.Flight, error) {
	url := fmt.Sprintf("%s/flights", r.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	
	if err != nil {
		return nil, fmt.Errorf("flights build request: %w", err)
	}

	resp, err := r.client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("flights GET %s: %w", url, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<14))
		return nil, fmt.Errorf("flights status %d: %s", resp.StatusCode, string(body))
	}

	var items []model.FlightItem

	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, fmt.Errorf("flights decode array: %w", err)
	}

	out := make([]domain.Flight, 0, len(items))
	for _, f := range items {
		dep, err := time.Parse(time.RFC3339, f.DepartureTime)

		if err != nil {
			return nil, fmt.Errorf("flights bad departureTime %q: %w", f.DepartureTime, err)
		}

		arr, err := time.Parse(time.RFC3339, f.ArrivalTime)

		if err != nil {
			return nil, fmt.Errorf("flights bad arrivalTime %q: %w", f.ArrivalTime, err)
		}

		out = append(out, domain.Flight{
			Reference:     f.BookingID,
			FlightNumber:  f.FlightNumber,
			From:          f.DepartureAirport,
			To:            f.ArrivalAirport,
			DepartureTime: dep,
			ArrivalTime:   arr,
			Price:         f.Price,
			Currency:      f.Currency,
		})
	}

	return out, nil
}
