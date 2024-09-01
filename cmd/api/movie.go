package main

import (
	"errors"
	"net/http"

	"github.com/gaba-bouliva/movie-api/internal/data"
	"github.com/gaba-bouliva/movie-api/internal/validator"
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

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct{
		Title			string				`json:"title"`
		Year 			int32					`json:"year"`
		Runtime 	data.Runtime	`json:"runtime"`
		Genres 		[]string			`json:"genres"`
	}

	err := app.readJSON(w,r,&input)
	if err != nil {
		app.logErr(err)
		app.badRequestError(w,r,err)
		return 
	}

	movie := &data.Movie{
		Title: input.Title,
		Year: input.Year,
		Runtime: input.Runtime,
		Genres: input.Genres,
	}

	v := validator.New()

	data.ValidateMovie(v, movie)

	if ok := v.Valid(); !ok{
		app.failedValidationError(w,r,v.Errors)	
		return
	}

	err = app.models.Movies.Create(movie)
	if err != nil {
		app.logErr(err)
		app.serverError(w,r)
		return
	}

	app.writeJSON(w,jsonResponse{"movie": movie}, http.StatusOK, nil)
	
}