package main

import (
	"errors"
	"net/http"

	"github.com/gaba-bouliva/movie-api/internal/data"
)

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestError(w,r,err)
		return
	}

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		app.logErr(err)
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundError(w,r)
		default:
			app.serverError(w,r)	
		}
		return
	}

	app.writeJSON(w,jsonResponse{"movie": movie}, http.StatusOK, nil)

}