package sorter

import (
	"sort"
	"strings"

	"github.com/Orden14/flight-aggregator/src/domain"
)

type SortBy string

const (
	SortByPrice         SortBy = "price"
	SortByDepartureDate SortBy = "departure_date"
	SortByTravelTime    SortBy = "travel_time"
)

type Order string

const (
	OrderAsc  Order = "asc"
	OrderDesc Order = "desc"
)

func NormalizeSortBy(inputValue string) SortBy {
	switch strings.ToLower(inputValue) {
	case "price":
		return SortByPrice
	case "travel_time", "duration":
		return SortByTravelTime
	case "departure_date", "departure":
		return SortByDepartureDate
	default:
		return SortByPrice
	}
}

func NormalizeOrder(inputValue string) Order {
	switch strings.ToLower(inputValue) {
	case "desc", "descending":
		return OrderDesc
	default:
		return OrderAsc
	}
}

func SortFlights(flights []domain.Flight, sortBy SortBy, sortOrder Order) {
	compareFlights := func(i, j int) bool { return false }

	switch sortBy {
	case SortByTravelTime:
		compareFlights = func(i, j int) bool { return flights[i].Duration() < flights[j].Duration() }
	case SortByDepartureDate:
		compareFlights = func(i, j int) bool { return flights[i].DepartureTime.Before(flights[j].DepartureTime) }
	default: // SortByPrice
		compareFlights = func(i, j int) bool { return flights[i].Price < flights[j].Price }
	}

	sort.SliceStable(flights, func(i, j int) bool {
		if sortOrder == OrderAsc {
			return compareFlights(i, j)
		}

		return compareFlights(j, i)
	})
}
