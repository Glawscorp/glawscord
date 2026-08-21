package server

import (
	"github.com/go-chi/chi/v5"
)

func InitServer() *chi.Mux {
	r := chi.NewRouter()

	userRoutes(r)
	userMessageRoutes(r)

	return r
}
