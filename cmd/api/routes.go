package main

import "github.com/julienschmidt/httprouter"

// encapsulate all routing rules
func (app *application) routes() *httprouter.Router {

	//initialize an httprouter instance \
	router := httprouter.New()

	//register the relevant methods, URL patterns and handler functions for our
	//endpoints using the HandlerFunc() method.
	router.HandlerFunc("GET", "/v1/healthcheck", app.healthCheckHandler)
	router.HandlerFunc("POST", "/v1/movies", app.createMovieHandler)
	router.HandlerFunc("GET", "/v1/movies/:id", app.showMovieHandler)

	//return the http router instance
	return router
}
