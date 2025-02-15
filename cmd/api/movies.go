package main

import (
	"fmt"
	"net/http"
	"time"

	"greenlight.brainwhat/internal/data"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Title   string   `json:"title"`
		Year    int32    `json:"year"`
		Runtime int32    `json:"runtime"`
		Genres  []string `json:"genres"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestReponse(w, r, err)
		return
	}

	fmt.Fprintf(w, "%+v", input)
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParams(r)
	if err != nil {
		app.notFoundError(w, r)
		return
	}

	// We didn't set the value of Runtime. It'll be set to zero by json.Marshal
	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Titanic",
		Runtime:   130,
		Year:      0,
		Genres:    []string{"sexy dicaprio", "drama", "tragedy"},
		Version:   1,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
