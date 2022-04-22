package main

import (
	"fmt"
	"greenlight.alexedwards.net/internal/data"
	"greenlight.alexedwards.net/internal/validator"
	"net/http"
	"time"
)

//add a createMovieHandler for the POST "v1/movies" endpoint.
func (app *application) createMovieHandler(writer http.ResponseWriter, request *http.Request) {
	//anonymous struct to hold information expected in request
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	//use readJSON() to decode request body into input struct
	err := app.readJSON(writer, request, &input)
	if err != nil {
		app.badRequestResponse(writer, request, err)
		return
	}

	//Copy the values from the input struct to a new movie struct
	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		RunTime: input.Runtime,
		Genres:  input.Genres,
	}
	//initialize a new validator instance
	v := validator.New()

	//Call the ValidateMovie() function and return a response containing the errors if any of the checks fail
	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(writer, request, v.Errors)
		return
	}

	fmt.Fprintf(writer, "%+v\n", input)
}

//Add a showMovieHandler for the GET "/v1/movies/:id" endpoint.
func (app *application) showMovieHandler(writer http.ResponseWriter, request *http.Request) {
	id, err := app.readIDParam(request)
	if err != nil {
		app.notFoundResponse(writer, request)
	}

	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Year:      2020,
		RunTime:   102,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}

	err = app.writeJSON(writer, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
	}
}
