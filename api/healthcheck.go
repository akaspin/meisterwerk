package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func WithHealthcheck(m *mux.Router) {
	m.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}
