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

func NewServer2FlightRepository(c config.JSONServerConfig) *Server2FlightRepository {
	return &Server2FlightRepository{
		baseURL: c.BaseURL(),
		client:  &http.Client{Timeout: 0},
	}
}

func (r *Server2FlightRepository) Fetch(ctx context.Context) ([]domain.Flight, error) {
	url := fmt.Sprintf("%s/flight_to_book", r.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return nil, fmt.Errorf("flight build request: %w", err)
	}

	resp, err := r.client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("flight GET %s: %w", url, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<14))
		return nil, fmt.Errorf("flight status %d: %s", resp.StatusCode, string(body))
	}

	var items []model.Server2FlightItem

	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, fmt.Errorf("flight decode array: %w", err)
	}

	out := make([]domain.Flight, 0, len(items))

	for _, it := range items {
		if len(it.Segments) == 0 {
			continue
		}

		first := it.Segments[0].Flight
		last := it.Segments[len(it.Segments)-1].Flight

		dep, err := time.Parse(time.RFC3339, first.Depart)

		if err != nil {
			return nil, fmt.Errorf("flight bad depart %q: %w", first.Depart, err)
		}

		arr, err := time.Parse(time.RFC3339, last.Arrive)

		if err != nil {
			return nil, fmt.Errorf("flight bad arrive %q: %w", last.Arrive, err)
		}

		out = append(out, domain.Flight{
			Reference:     it.Reference,
			FlightNumber:  first.Number,
			From:          first.From,
			To:            last.To,
			DepartureTime: dep,
			ArrivalTime:   arr,
			Price:         it.Total.Amount,
			Currency:      it.Total.Currency,
		})
	}

	return out, nil
}
