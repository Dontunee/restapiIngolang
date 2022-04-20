package main

import (
	"fmt"
	"net/http"
)

//add a createMovieHandler for the POST "v1/movies" endpoint.
func (app *application) createMovieHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprintln(writer, "creating a new movie")
}

//Add a showMovieHandler for the GET "/v1/movies/:d" endpoint.
func (app *application) showMovieHandler(writer http.ResponseWriter, request *http.Request) {

	id, err := app.readIDParam(request)
	if err != nil {
		http.NotFound(writer, request)
		return
	}

	fmt.Fprintf(writer, "show the details of the movie %d\n", id)
}
