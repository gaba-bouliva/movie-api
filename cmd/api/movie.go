package main

import (
	"database/sql"
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

	err = app.writeJSON(w,jsonResponse{"movie": movie}, http.StatusOK, nil)
	if err != nil {
		app.logErr(err)
		app.serverError(w,r)
	}

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

	err = app.writeJSON(w,jsonResponse{"movie": movie}, http.StatusOK, nil)
	if err != nil {
		app.logErr(err)
		app.serverError(w,r)
	}
}

func (app *application) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestError(w,r,err)
		return
	}

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		app.logErr(err)
		switch {
			case errors.Is(err, sql.ErrNoRows):
				app.notFoundError(w,r)
			default:
				app.serverError(w,r)
		}
		return
	}

	var input struct{
		Title			string				`json:"title"`
		Year 			int32					`json:"year"`
		Runtime 	data.Runtime	`json:"runtime"`
		Genres 		[]string			`json:"genres"`
	}

	err = app.readJSON(w,r,&input)
	if err != nil {
		app.logErr(err)
		app.badRequestError(w,r,err)
		return 
	}

	movie.Title = input.Title
	movie.Year = input.Year
	movie.Runtime = input.Runtime
	movie.Genres = input.Genres

	v := validator.New()

	data.ValidateMovie(v, movie)

	if ok := v.Valid(); !ok{
		app.failedValidationError(w,r,v.Errors)	
		return
	}

	err = app.models.Movies.Update(movie)
	if err != nil {
		app.logErr(err)
		app.serverError(w,r)
		return
	}

	err = app.writeJSON(w,jsonResponse{"message": "movie updated successfully"}, http.StatusOK, nil)
	if err != nil {
		app.logErr(err)
		app.serverError(w,r)
	}
}

func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestError(w,r,err)
		return
	}

	err = app.models.Movies.Delete(id)
	if err != nil {
		app.notFoundError(w,r)
		return
	}

	err = app.writeJSON(w,jsonResponse{"message": "movie deleted successfully"}, http.StatusOK, nil)
	if err != nil {
		app.logErr(err)
		app.serverError(w,r)
	}
}