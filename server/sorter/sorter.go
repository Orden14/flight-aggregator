package sorter

import (
	"sort"
	"strings"

	"github.com/Orden14/flight-aggregator/domain"
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

func NormalizeSortBy(v string) SortBy {
	switch strings.ToLower(v) {
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

func NormalizeOrder(v string) Order {
	switch strings.ToLower(v) {
	case "desc", "descending":
		return OrderDesc
	default:
		return OrderAsc
	}
}

func SortFlights(f []domain.Flight, by SortBy, order Order) {
	less := func(i, j int) bool { return false }

	switch by {
	case SortByTravelTime:
		less = func(i, j int) bool { return f[i].Duration() < f[j].Duration() }
	case SortByDepartureDate:
		less = func(i, j int) bool { return f[i].DepartureTime.Before(f[j].DepartureTime) }
	default: // SortByPrice
		less = func(i, j int) bool { return f[i].Price < f[j].Price }
	}

	sort.SliceStable(f, func(i, j int) bool {
		if order == OrderAsc {
			return less(i, j)
		}
		
		return less(j, i)
	})
}
