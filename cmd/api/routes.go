package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// encapsulate all routing rules
func (app *application) routes() *httprouter.Router {

	//initialize a httprouter instance
	router := httprouter.New()

	//convert the notFoundResponse() helper to a http.handler
	//using the http.HandlerFunc() adapter, and then set as the custom error
	// handler for 404 Not Found Responses.
	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	//convert the methodNotAllowedResponse() helper to a http.Handler and set
	//it as the custom error handler for 405 method not allowed responses.
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	//register the relevant methods, URL patterns and handler functions for our
	//endpoints using the HandlerFunc() method.
	router.HandlerFunc("GET", "/v1/healthcheck", app.healthCheckHandler)
	router.HandlerFunc("POST", "/v1/movies", app.createMovieHandler)
	router.HandlerFunc("GET", "/v1/movies/:id", app.showMovieHandler)
	router.HandlerFunc("PATCH", "/v1/movies/:id", app.updateMovieHandler)
	router.HandlerFunc("DELETE", "/v1/movies/:id",app.deleteMovieHandler)


	//return the http router instance
	return router
}
