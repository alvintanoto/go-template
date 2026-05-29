package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(
	logger *ZapMiddleware,
) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(logger.Execute)

	r.Get("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	return r
}
