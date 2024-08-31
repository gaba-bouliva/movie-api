package main

import (
	"net/http"
)

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestError(w,r,err)
		return
	}
	app.logger.Printf("id: %d provided", id)

}