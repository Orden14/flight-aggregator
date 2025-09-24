package test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Orden14/flight-aggregator/src/domain"
	"github.com/Orden14/flight-aggregator/src/repository"
	"github.com/Orden14/flight-aggregator/src/service"
	"github.com/stretchr/testify/require"
)

var _ repository.FlightRepositoryInterface = (*MockRepo)(nil)

type MockRepo struct {
	FetchFunc func(ctx context.Context) ([]domain.Flight, error)
}

func (m *MockRepo) Fetch(ctx context.Context) ([]domain.Flight, error) {
	if m.FetchFunc == nil {
		return nil, nil
	}

	return m.FetchFunc(ctx)
}

func tTime(t *testing.T, s string) time.Time {
	tt, err := time.Parse(time.RFC3339, s)
	require.NoError(t, err)
	return tt
}

func TestFiltersAndSorts(t *testing.T) {
	repoA := &MockRepo{
		FetchFunc: func(ctx context.Context) ([]domain.Flight, error) {
			return []domain.Flight{
				{
					Reference:     "REF-1",
					From:          "CDG",
					To:            "HND",
					Price:         900,
					DepartureTime: tTime(t, "2026-01-01T10:00:00Z"),
					ArrivalTime:   tTime(t, "2026-01-01T20:00:00Z"),
				},
				{
					Reference:     "REF-2",
					From:          "CDG",
					To:            "NRT",
					Price:         700,
					DepartureTime: tTime(t, "2026-01-01T07:00:00Z"),
					ArrivalTime:   tTime(t, "2026-01-01T18:00:00Z"),
				},
			}, nil
		},
	}

	repoB := &MockRepo{
		FetchFunc: func(ctx context.Context) ([]domain.Flight, error) {
			return []domain.Flight{
				{
					Reference:     "REF-1",
					From:          "CDG",
					To:            "HND",
					Price:         800,
					DepartureTime: tTime(t, "2026-01-01T09:00:00Z"),
					ArrivalTime:   tTime(t, "2026-01-01T20:00:00Z"),
				},
				{
					Reference:     "REF-3",
					From:          "CDG",
					To:            "HND",
					Price:         850,
					DepartureTime: tTime(t, "2026-01-01T06:00:00Z"),
					ArrivalTime:   tTime(t, "2026-01-01T15:00:00Z"),
				},
			}, nil
		},
	}

	svc := service.NewFlightService(2, repoA, repoB)

	ctx := context.Background()

	flights, err := svc.GetFlights(ctx, "CDG", "HND", service.SortByPrice, service.OrderAsc)
	require.NoError(t, err)

	require.Len(t, flights, 2)
	require.Equal(t, "REF-1", flights[0].Reference)
	require.Equal(t, float64(800), flights[0].Price)
	require.Equal(t, "REF-3", flights[1].Reference)
	require.Equal(t, float64(850), flights[1].Price)

	require.Equal(t, 11*60, flights[0].TravelTimeMinutes)
	require.Equal(t, 9*60, flights[1].TravelTimeMinutes)
}

func TestDedupPolicyCheapestThenEarliest(t *testing.T) {
	repoA := &MockRepo{
		FetchFunc: func(ctx context.Context) ([]domain.Flight, error) {
			return []domain.Flight{
				{
					Reference:     "DUP",
					From:          "SFO",
					To:            "LAX",
					Price:         120,
					DepartureTime: tTime(t, "2026-01-01T10:00:00Z"),
					ArrivalTime:   tTime(t, "2026-01-01T11:30:00Z"),
				},
			}, nil
		},
	}

	repoB := &MockRepo{
		FetchFunc: func(ctx context.Context) ([]domain.Flight, error) {
			return []domain.Flight{
				{
					Reference:     "DUP",
					From:          "SFO",
					To:            "LAX",
					Price:         100,
					DepartureTime: tTime(t, "2026-01-01T13:00:00Z"),
					ArrivalTime:   tTime(t, "2026-01-01T14:30:00Z"),
				},
				{
					Reference:     "DUP",
					From:          "SFO",
					To:            "LAX",
					Price:         100,
					DepartureTime: tTime(t, "2026-01-01T09:00:00Z"),
					ArrivalTime:   tTime(t, "2026-01-01T10:30:00Z"),
				},
			}, nil
		},
	}

	svc := service.NewFlightService(3, repoA, repoB)

	out, err := svc.GetFlights(context.Background(), "", "", service.SortByPrice, service.OrderAsc)
	require.NoError(t, err)
	require.Len(t, out, 1)
	require.Equal(t, "DUP", out[0].Reference)
	require.Equal(t, float64(100), out[0].Price)
	require.Equal(t, tTime(t, "2026-01-01T09:00:00Z"), out[0].DepartureTime)
}

func TestTimeoutErrorPropagates(t *testing.T) {
	blockingRepo := &MockRepo{
		FetchFunc: func(ctx context.Context) ([]domain.Flight, error) {
			<-ctx.Done()
			return nil, ctx.Err()
		},
	}

	okRepo := &MockRepo{
		FetchFunc: func(ctx context.Context) ([]domain.Flight, error) {
			return []domain.Flight{}, nil
		},
	}

	flightService := service.NewFlightService(1, blockingRepo, okRepo)

	start := time.Now()
	_, err := flightService.GetFlights(context.Background(), "", "", service.SortByPrice, service.OrderAsc)
	elapsed := time.Since(start)

	require.Error(t, err)
	require.True(t, errors.Is(err, context.DeadlineExceeded), "expected context deadline exceeded, got: %v", err)

	require.Less(t, elapsed, 3*time.Second, "timeout test ran too long")
}
