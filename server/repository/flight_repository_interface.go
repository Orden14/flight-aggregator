package repository

import (
	"context"

	"github.com/Orden14/flight-aggregator/domain"
)

type FlightRepositoryInterface interface {
	Fetch(ctx context.Context) ([]domain.Flight, error)
}
