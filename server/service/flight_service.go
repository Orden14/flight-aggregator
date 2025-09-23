package service

import (
	"errors"
	"sync"

	"github.com/Orden14/flight-aggregator/domain"
	"github.com/Orden14/flight-aggregator/repository"
	"github.com/Orden14/flight-aggregator/sorter"
)

type FlightService interface {
	GetFlights(from string, to string, by sorter.SortBy, order sorter.Order) ([]domain.Flight, error)
}

type flightService struct {
	repos []repository.FlightRepositoryInterface
}

func NewFlightService(repos ...repository.FlightRepositoryInterface) FlightService {
	return &flightService{repos: repos}
}

func (s *flightService) GetFlights(from string, to string, by sorter.SortBy, order sorter.Order) ([]domain.Flight, error) {
	if len(s.repos) == 0 {
		return nil, errors.New("no repositories configured")
	}

	all, err := s.fetchAll()

	if err != nil {
		return nil, err
	}

	filtered := s.filterFlights(all, from, to)
	sorter.SortFlights(filtered, by, order)
	s.enrichFlights(&filtered)

	return filtered, nil
}

func (s *flightService) fetchAll() ([]domain.Flight, error) {
	var wg sync.WaitGroup

	results := make(chan []domain.Flight, len(s.repos))
	errs := make(chan error, len(s.repos))

	for _, repo := range s.repos {
		wg.Add(1)

		go func(r repository.FlightRepositoryInterface) {
			defer wg.Done()
			flights, err := r.Fetch()

			if err != nil {
				errs <- err
				return
			}

			results <- flights
		}(repo)
	}

	wg.Wait()
	close(results)
	close(errs)

	if firstErr := firstError(errs); firstErr != nil {
		return nil, firstErr
	}

	var all []domain.Flight

	for fs := range results {
		all = append(all, fs...)
	}

	return all, nil
}

func (s *flightService) filterFlights(in []domain.Flight, from, to string) []domain.Flight {
	if from == "" && to == "" {
		out := make([]domain.Flight, len(in))
		copy(out, in)

		return out
	}

	out := make([]domain.Flight, 0, len(in))

	for _, f := range in {
		if from != "" && f.From != from {
			continue
		}

		if to != "" && f.To != to {
			continue
		}

		out = append(out, f)
	}

	return out
}

func (s *flightService) sortFlights(flights []domain.Flight, by sorter.SortBy, order sorter.Order) {
	sorter.SortFlights(flights, by, order)
}

func (s *flightService) enrichFlights(f *[]domain.Flight) {
	for i := range *f {
		(*f)[i].TravelTimeMinutes = int((*f)[i].Duration().Minutes())
	}
}

func firstError(errs <-chan error) error {
	for err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}
