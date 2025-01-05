package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() *chi.Mux {
	router := chi.NewRouter()

	router.NotFound(http.HandlerFunc(app.notFoundResponse))

	router.Get("/v1/healthcheck", app.healthcheckHandler)

	router.Route("/v1/movies", func(r chi.Router) {
		r.Get("/{id}", app.showMovieHandler)
		r.Post("/", app.createMovieHandler)
		r.Put("/{id}", app.updateMovieHandler)
		r.Delete("/{id}", app.deleteMovieHandler)
	})

	return router
}
