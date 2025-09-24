package service

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/Orden14/flight-aggregator/src/domain"
	"github.com/Orden14/flight-aggregator/src/repository"
	"github.com/Orden14/flight-aggregator/src/util/errtools"
	"github.com/Orden14/flight-aggregator/src/util/sorter"
)

type FlightService interface {
	GetFlights(ctx context.Context, departureAirport string, arrivalAirport string, sortBy sorter.SortBy, sortOrder sorter.Order) ([]domain.Flight, error)
}

type flightService struct {
	repositories      []repository.FlightRepositoryInterface
	repositoryTimeout time.Duration
}

func NewFlightService(timeout time.Duration, repositories ...repository.FlightRepositoryInterface) FlightService {
	if timeout <= 0 {
		timeout = 10 * time.Second
	}

	return &flightService{
		repositories:      repositories,
		repositoryTimeout: timeout * time.Second,
	}
}

func (flightService *flightService) GetFlights(ctx context.Context, departureAirport string, arrivalAirport string, sortBy sorter.SortBy, sortOrder sorter.Order) ([]domain.Flight, error) {
	if len(flightService.repositories) == 0 {
		return nil, errors.New("no repositories configured")
	}

	flights, err := flightService.fetchAll(ctx)
	if err != nil {
		return nil, err
	}

	flights = flightService.dedupeFlights(flights)
	filteredFlights := flightService.filterFlights(flights, departureAirport, arrivalAirport)
	sorter.SortFlights(filteredFlights, sortBy, sortOrder)
	flightService.enrichFlights(&filteredFlights)

	return filteredFlights, nil
}

func (flightService *flightService) fetchAll(ctx context.Context) ([]domain.Flight, error) {
	var waitGroup sync.WaitGroup

	results := make(chan []domain.Flight, len(flightService.repositories))
	errs := make(chan error, len(flightService.repositories))

	for _, flightRepository := range flightService.repositories {
		waitGroup.Add(1)

		go func(r repository.FlightRepositoryInterface) {
			defer waitGroup.Done()

			requestContext, cancel := context.WithTimeout(ctx, flightService.repositoryTimeout)
			defer cancel()

			flights, err := r.Fetch(requestContext)

			if err != nil {
				errs <- err

				return
			}

			results <- flights
		}(flightRepository)
	}

	waitGroup.Wait()
	close(results)
	close(errs)

	if firstErr := errtools.GetFirstError(errs); firstErr != nil {
		return nil, firstErr
	}

	var flights []domain.Flight

	for flight := range results {
		flights = append(flights, flight...)
	}

	return flights, nil
}

func (flightService *flightService) dedupeFlights(flights []domain.Flight) []domain.Flight {
	if len(flights) <= 1 {
		out := make([]domain.Flight, len(flights))
		copy(out, flights)

		return out
	}

	optimalFlights := make(map[string]domain.Flight, len(flights))

	for _, flight := range flights {
		if currentFlight, isAlreadyExisting := optimalFlights[flight.Reference]; isAlreadyExisting {
			if flight.Price < currentFlight.Price || (flight.Price == currentFlight.Price && flight.DepartureTime.Before(currentFlight.DepartureTime)) {
				optimalFlights[flight.Reference] = flight
			}
		} else {
			optimalFlights[flight.Reference] = flight
		}
	}

	dedupedFlights := make([]domain.Flight, 0, len(optimalFlights))

	for _, flight := range optimalFlights {
		dedupedFlights = append(dedupedFlights, flight)
	}

	return dedupedFlights
}

func (flightService *flightService) filterFlights(flights []domain.Flight, departureAirport string, arrivalAirport string) []domain.Flight {
	if departureAirport == "" && arrivalAirport == "" {
		filteredFlights := make([]domain.Flight, len(flights))
		copy(filteredFlights, flights)

		return filteredFlights
	}

	filteredFlights := make([]domain.Flight, 0, len(flights))

	for _, flight := range flights {
		if departureAirport != "" && flight.From != departureAirport {
			continue
		}

		if arrivalAirport != "" && flight.To != arrivalAirport {
			continue
		}

		filteredFlights = append(filteredFlights, flight)
	}

	return filteredFlights
}

func (flightService *flightService) enrichFlights(flights *[]domain.Flight) {
	for i := range *flights {
		(*flights)[i].TravelTimeMinutes = int((*flights)[i].Duration().Minutes())
	}
}
