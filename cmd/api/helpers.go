package main

import (
	"encoding/json"
	"errors"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

type envelope map[string]interface{}

func (app *application) readIDParam(request *http.Request) (int64, error) {

	params := httprouter.ParamsFromContext(request.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

func (app *application) writeJson(writer http.ResponseWriter, status int, data envelope, headers http.Header) error {
	//Encode the data to JSON, returning the error if there was one
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	//Append a new line to make it easier to view in terminal applications
	js = append(js, '\n')

	//loop through header, does not throw error if it is empty
	for key, value := range headers {
		writer.Header()[key] = value
	}

	//Add the Content-Type : application/json header, then write the status code and JSON response
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	writer.Write(js)

	return nil
}
