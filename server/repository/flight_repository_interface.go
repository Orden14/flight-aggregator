package repository

import (
	"github.com/Orden14/flight-aggregator/domain"
)

type FlightRepositoryInterface interface {
	Fetch() ([]domain.Flight, error)
}
