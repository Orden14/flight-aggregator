package service

import (
	"errors"
	"sync"

	"github.com/Orden14/flight-aggregator/domain"
	"github.com/Orden14/flight-aggregator/repository"
)

type FlightService interface {
	GetFlights() ([]domain.Flight, error)
}

type flightService struct {
	repos []repository.FlightRepositoryInterface
}

func NewFlightService(repos ...repository.FlightRepositoryInterface) FlightService {
	return &flightService{repos: repos}
}

func (s *flightService) GetFlights() ([]domain.Flight, error) {
	if len(s.repos) == 0 {
		return nil, errors.New("no repositories configured")
	}

	var wg sync.WaitGroup

	type res struct {
		flights []domain.Flight
		err     error
	}

	ch := make(chan res, len(s.repos))

	for _, r := range s.repos {
		wg.Add(1)

		go func(repo repository.FlightRepositoryInterface) {
			defer wg.Done()
			flights, err := repo.Fetch()
			ch <- res{flights: flights, err: err}
		}(r)
	}

	wg.Wait()
	close(ch)

	var all []domain.Flight

	for r := range ch {
		if r.err != nil {
			return nil, r.err
		}

		all = append(all, r.flights...)
	}

	return all, nil
}
