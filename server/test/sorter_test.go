package test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/Orden14/flight-aggregator/domain"
	"github.com/Orden14/flight-aggregator/sorter"
)

func mustRFC3339(t *testing.T, s string) time.Time {
	tm, err := time.Parse(time.RFC3339, s)
	require.NoError(t, err)

	return tm
}

func sampleFlights(t *testing.T) []domain.Flight {
	return []domain.Flight{
		{
			Reference:     "R1",
			Price:         900,
			DepartureTime: mustRFC3339(t, "2026-01-01T10:00:00Z"),
			ArrivalTime:   mustRFC3339(t, "2026-01-01T20:00:00Z"),
		},
		{
			Reference:     "R2",
			Price:         800,
			DepartureTime: mustRFC3339(t, "2026-01-01T09:00:00Z"),
			ArrivalTime:   mustRFC3339(t, "2026-01-01T18:00:00Z"),
		},
		{
			Reference:     "R3",
			Price:         950,
			DepartureTime: mustRFC3339(t, "2026-01-01T08:00:00Z"),
			ArrivalTime:   mustRFC3339(t, "2026-01-01T22:00:00Z"),
		},
	}
}

func TestSortByPriceAsc(t *testing.T) {
	flights := sampleFlights(t)
	sorter.SortFlights(flights, sorter.SortByPrice, sorter.OrderAsc)

	require.Equal(t, "R2", flights[0].Reference)
	require.Equal(t, "R1", flights[1].Reference)
	require.Equal(t, "R3", flights[2].Reference)
}

func TestSortByDepartureDesc(t *testing.T) {
	flights := sampleFlights(t)
	sorter.SortFlights(flights, sorter.SortByDepartureDate, sorter.OrderDesc)

	require.Equal(t, "R1", flights[0].Reference)
	require.Equal(t, "R2", flights[1].Reference)
	require.Equal(t, "R3", flights[2].Reference)
}

func TestSortByTravelTimeAsc(t *testing.T) {
	flights := sampleFlights(t)
	sorter.SortFlights(flights, sorter.SortByTravelTime, sorter.OrderAsc)

	require.Equal(t, "R2", flights[0].Reference)
	require.Equal(t, "R1", flights[1].Reference)
	require.Equal(t, "R3", flights[2].Reference)
}
