package model

type FlightItem struct {
	BookingID        string  `json:"bookingId"`
	Status           string  `json:"status"`
	PassengerName    string  `json:"passengerName"`
	FlightNumber     string  `json:"flightNumber"`
	DepartureAirport string  `json:"departureAirport"`
	ArrivalAirport   string  `json:"arrivalAirport"`
	DepartureTime    string  `json:"departureTime"`
	ArrivalTime      string  `json:"arrivalTime"`
	Price            float64 `json:"price"`
	Currency         string  `json:"currency"`
}
