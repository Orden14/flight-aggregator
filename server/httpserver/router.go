package httpserver

import (
	"net/http"

	"github.com/Orden14/flight-aggregator/handler"
)

func NewRouter(health *handler.HealthHandler, flights *handler.FlightHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

			return
		}

		health.ServeHTTP(w, r)
	})

	mux.HandleFunc("/flights", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

			return
		}

		flights.ServeHTTP(w, r)
	})

	return mux
}
