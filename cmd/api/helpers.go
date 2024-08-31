package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type jsonResponse map[string]any

func (app *application) readIDParam(r *http.Request) (int, error) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id < 1{
		return 0, fmt.Errorf("invalid id provided")
	}

	return id, nil
}

func (app *application) writeJSON(
	w http.ResponseWriter,
	payload jsonResponse, 
	status int, 
	headers http.Header) error {

	jsonRes, err := json.MarshalIndent(payload, "", "\t")
	if err != nil {
		return err
	}

	jsonRes = append(jsonRes, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	w.Write(jsonRes)

	return nil
} 