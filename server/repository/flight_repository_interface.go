package repository

import (
	"context"

	"github.com/Orden14/flight-aggregator/domain"
)

type FlightRepositoryInterface interface {
	Fetch(context context.Context) ([]domain.Flight, error)
}
