package httpserver

import (
	"net/http"

	"github.com/Orden14/flight-aggregator/src/handler"
)

func NewRouter(healthHandler *handler.HealthHandler, flightHandler *handler.FlightHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			writer.Header().Set("Allow", http.MethodGet)
			http.Error(writer, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

			return
		}

		healthHandler.ServeHTTP(writer)
	})

	mux.HandleFunc("/flights", func(writer http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodGet {
			writer.Header().Set("Allow", http.MethodGet)
			http.Error(writer, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)

			return
		}

		flightHandler.ServeHTTP(writer, request)
	})

	return mux
}
