package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	
	r.NotFound(app.notFoundError)
	r.MethodNotAllowed(app.methodNotAllowed)

	apiRouter := chi.NewRouter()
	apiRouter.Get("/movies/{id}", app.showMovieHandler)
	apiRouter.Post("/movies", app.createMovieHandler)
	apiRouter.Put("/movies/{id}", app.updateMovieHandler)
	apiRouter.Delete("/movies/{id}", app.deleteMovieHandler)
	apiRouter.Get("/movies", app.listMoviesHandler)

	r.Mount("/api/v1", apiRouter)

	return r
}