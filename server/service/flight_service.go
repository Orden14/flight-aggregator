package service

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/Orden14/flight-aggregator/domain"
	"github.com/Orden14/flight-aggregator/repository"
	"github.com/Orden14/flight-aggregator/sorter"
)

type FlightService interface {
	GetFlights(ctx context.Context, from, to string, by sorter.SortBy, order sorter.Order) ([]domain.Flight, error)
}

type flightService struct {
	repos       []repository.FlightRepositoryInterface
	repoTimeout time.Duration
}

func NewFlightService(timeout time.Duration, repos ...repository.FlightRepositoryInterface) FlightService {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	return &flightService{
		repos:       repos,
		repoTimeout: timeout * time.Second,
	}
}

func (s *flightService) GetFlights(ctx context.Context, from, to string, by sorter.SortBy, order sorter.Order) ([]domain.Flight, error) {
	if len(s.repos) == 0 {
		return nil, errors.New("no repositories configured")
	}

	all, err := s.fetchAll(ctx)
	if err != nil {
		return nil, err
	}

	all = s.dedupeFlights(all)
	filtered := s.filterFlights(all, from, to)
	sorter.SortFlights(filtered, by, order)
	s.enrichFlights(&filtered)

	return filtered, nil
}

func (s *flightService) fetchAll(ctx context.Context) ([]domain.Flight, error) {
	var wg sync.WaitGroup

	results := make(chan []domain.Flight, len(s.repos))
	errs := make(chan error, len(s.repos))

	for _, repo := range s.repos {
		wg.Add(1)
		go func(r repository.FlightRepositoryInterface) {
			defer wg.Done()

			reqCtx, cancel := context.WithTimeout(ctx, s.repoTimeout)
			defer cancel()

			flights, err := r.Fetch(reqCtx)

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

func (s *flightService) dedupeFlights(in []domain.Flight) []domain.Flight {
	if len(in) <= 1 {
		out := make([]domain.Flight, len(in))
		copy(out, in)

		return out
	}

	best := make(map[string]domain.Flight, len(in))

	for _, f := range in {
		if cur, ok := best[f.Reference]; ok {
			if f.Price < cur.Price || (f.Price == cur.Price && f.DepartureTime.Before(cur.DepartureTime)) {
				best[f.Reference] = f
			}
		} else {
			best[f.Reference] = f
		}
	}

	out := make([]domain.Flight, 0, len(best))
	for _, f := range best {
		out = append(out, f)
	}

	return out
}

func (s *flightService) filterFlights(in []domain.Flight, from, to string) []domain.Flight {
	if from == "" && to == "" {
		out := make([]domain.Flight, len(in))
		copy(out, in)

		return out
	}

	out := make([]domain.Flight, 0, len(in))

	for _, flight := range in {
		if from != "" && flight.From != from {
			continue
		}

		if to != "" && flight.To != to {
			continue
		}

		out = append(out, flight)
	}

	return out
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
