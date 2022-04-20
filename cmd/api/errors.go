package main

import (
	"fmt"
	"net/http"
)

//The logError() method generic helper for logging an error message
func (app *application) logError(request *http.Request, err error) {
	app.logger.Println(err)
}

//The errorResponse() method generic helper for sending JSON-formatted error messages to the client with a given status code.
//we are using an interface{} type for the message parameter, rather than just a string type, as this gives
// us more flexibility over the values that we can include in the response
func (app *application) errorResponse(writer http.ResponseWriter, request *http.Request, status int, message interface{}) {
	env := envelope{"error": message}

	// write the response using the writeJSON() helper. if this happens to return an error log it
	//and fall back to sending the client an empty response with 500 internal server error status code
	err := app.writeJson(writer, status, env, nil)
	if err != nil {
		app.logError(request, err)
		writer.WriteHeader(http.StatusInternalServerError)
	}
}

//The serverErrorResponse() method wil be used when our application encounters an unexpected problem at runtime.
//it details error message, then uses the errorResponse() helper to send a 500 internal server error status code
//JSON response (containing a generic error message) to the client
func (app *application) serverErrorResponse(writer http.ResponseWriter, request *http.Request, err error) {
	app.logError(request, err)

	message := "The server encountered a problem and could not process your request"
	app.errorResponse(writer, request, http.StatusInternalServerError, message)
}

//The notFoundResponse() method will be used to send a 404 Not found status code and JSON response to the client
func (app *application) notFoundResponse(writer http.ResponseWriter, request *http.Request) {
	message := "The requested resource could not be found"
	app.errorResponse(writer, request, http.StatusNotFound, message)
}

//The methodNotAllowedResponse() method will be used to send a 405 method not allowed status code
//and JSON response to the client.
func (app *application) methodNotAllowedResponse(writer http.ResponseWriter, request *http.Request) {
	message := fmt.Sprintf("The %s method is not supported for this resource", request.Method)
	app.errorResponse(writer, request, http.StatusMethodNotAllowed, message)
}
