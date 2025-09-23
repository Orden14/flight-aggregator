package domain

import "time"

type Flight struct {
	Reference         string    `json:"reference"`
	FlightNumber      string    `json:"flightNumber"`
	From              string    `json:"from"`
	To                string    `json:"to"`
	DepartureTime     time.Time `json:"departureTime"`
	ArrivalTime       time.Time `json:"arrivalTime"`
	Price             float64   `json:"price"`
	Currency          string    `json:"currency"`
	TravelTimeMinutes int       `json:"travelTimeMinutes"`
}

func (f Flight) Duration() time.Duration {
	return f.ArrivalTime.Sub(f.DepartureTime)
}
