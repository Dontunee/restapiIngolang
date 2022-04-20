package main

import (
	"net/http"
)

func (app *application) healthCheckHandler(writer http.ResponseWriter, request *http.Request) {
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	err := app.writeJSON(writer, 200, env, nil)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
	}
}
