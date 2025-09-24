package model

type Server2FlightItem struct {
	Reference string `json:"reference"`
	Status    string `json:"status"`
	Traveler  struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	} `json:"traveler"`
	Segments []struct {
		Flight struct {
			Number string `json:"number"`
			From   string `json:"from"`
			To     string `json:"to"`
			Depart string `json:"depart"`
			Arrive string `json:"arrive"`
		} `json:"flight"`
	} `json:"segments"`
	Total struct {
		Amount   float64 `json:"amount"`
		Currency string  `json:"currency"`
	} `json:"total"`
}
