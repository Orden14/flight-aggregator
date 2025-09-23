package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Orden14/flight-aggregator/config"
	"github.com/Orden14/flight-aggregator/domain"
	"github.com/Orden14/flight-aggregator/model"
)

type BookingRepository struct {
	baseURL string
	client  *http.Client
}

func NewBookingRepository(c config.JSONServerConfig) *BookingRepository {
	return &BookingRepository{
		baseURL: c.BaseURL(),
		client:  &http.Client{Timeout: 10 * time.Second},
	}
}

func (r *BookingRepository) Fetch() ([]domain.Flight, error) {
	url := fmt.Sprintf("%s/flight_to_book", r.baseURL)
	resp, err := r.client.Get(url)

	if err != nil {
		return nil, fmt.Errorf("bookings GET %s: %w", url, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<14))
		return nil, fmt.Errorf("bookings status %d: %s", resp.StatusCode, string(body))
	}

	var items []model.BookingItem

	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, fmt.Errorf("bookings decode array: %w", err)
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
			return nil, fmt.Errorf("bookings bad depart %q: %w", first.Depart, err)
		}

		arr, err := time.Parse(time.RFC3339, last.Arrive)

		if err != nil {
			return nil, fmt.Errorf("bookings bad arrive %q: %w", last.Arrive, err)
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
