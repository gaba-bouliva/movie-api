package main

import (
	"fmt"
	"net/http"
)

func (app *application) logErr(err error) {
	app.logger.Println(err)
}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, message any, status int) error {
	err := app.writeJSON(w, jsonResponse{"message": message}, status, nil)
	if err != nil {
		app.logErr(err)
		w.WriteHeader(500)
	}

	return nil
}

func (app *application) serverError(w http.ResponseWriter, r *http.Request) {
	message := "server encountered a problem and could not process your request"
	app.errorResponse(w,r,message, http.StatusInternalServerError)
} 

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request) {
	message := "resource not found"
	app.errorResponse(w,r,message, http.StatusNotFound)
}

func (app *application) methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the method %s not supported for this resource", r.Method)
	app.errorResponse(w,r,message,http.StatusMethodNotAllowed)
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w,r,err.Error(),http.StatusBadRequest)
}