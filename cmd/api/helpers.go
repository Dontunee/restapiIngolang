package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"greenlight.alexedwards.net/internal/validator"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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

//The readJSON() provides a reusable method to read json and handler errors centrally
func (app *application) readJSON(writer http.ResponseWriter, request *http.Request, destination interface{}) error {

	//use the http.MaxBytesReader to limit the size of the request body to 1MB
	maxBytes := 1_048_576
	request.Body = http.MaxBytesReader(writer, request.Body, int64(maxBytes))

	//initialize the json.Decoder, and  call the DisallowUnknownFields() method on it before decoding.
	//This means that if the JSON from the client now includes any field which cannot be mapped to the
	//target destination , the decoder will return an error instead of just ignoring the field
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()

	//Decode the request body into the target destination
	err := decoder.Decode(destination)
	if err == nil {
		return nil
	}

	//if there is an error during decoding ,start the triage...
	var syntaxError *json.SyntaxError
	var unmarshalTypeError *json.UnmarshalTypeError
	var invalidUnmarshalError *json.InvalidUnmarshalError

	switch {
	//use the errors.As() function to check whether the error has the type
	//*json.SyntaxError. If it does, then return a plain-english error message
	//which includes the location of the problem
	case errors.As(err, &syntaxError):
		return fmt.Errorf("body contains badly-formed JSON (at charcater %d)", syntaxError.Offset)

	//In some circumstances Decode() may also return an io.ErrUnexpectedEOF error
	//for syntax errors in the JSON. so we check for this using errors.IS()  and
	//return a generic error message.
	case errors.Is(err, io.ErrUnexpectedEOF):
		return errors.New("body contains badly-formed JSON")

	//catch any *json.UnmarshalTypeError errors. These occur when the
	//JSON value is the wrong type for the target destination. if the error relates
	// to a specific field , then we include that in our error message to make it
	//easier for the client to debug
	case errors.As(err, &unmarshalTypeError):
		if unmarshalTypeError.Field != "" {
			return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
		}
		return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

	//An io.EOF error will be returned by Decode() if the request body is empty. we
	//check for this with errors.Is() and return a plain-english error message
	//instead
	case errors.Is(err, io.EOF):
		return errors.New("body must not be empty")

	//A json.InvalidUnmarshalError error will be returned if we pass a non-nil
	//pointer to Decode().We catch this and panic , rather than returning an error
	//to our handler. At the end of this chapter we will talk about panicking versus returning
	// errors, and discuss why it is an appropriate thing to do in this specific situation
	case errors.As(err, &invalidUnmarshalError):
		panic(err)

	default:
		return err
	}

}

//The writeJSON() method is to encode object as json to write as response, set header and response code
func (app *application) writeJSON(writer http.ResponseWriter, status int, data envelope, headers http.Header) error {
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


//The readString() helper returns a string value from the query string, or the provided
//default value if no matching key could be found
func (app *application) readString(queryValues url.Values, key string, defaultValue string) string {
	//Extract the value for a give key from the query string. if empty it returns an empty string
	value := queryValues.Get(key)

	//if no key exists (or the value is empty) then return the default value
	if value == ""{
		return defaultValue
	}

	//Otherwise return the string
	return value
}


//The readCSV() helper reads a string value from the query string and then splits it
//into a slice on the comma character. if no matching key could not be found, it returns
//the provided default value
func (app *application) readCSV(queryValues url.Values, key string, defaultValue []string) []string {
	//extract the value from the query string
	csv := queryValues.Get(key)

	//if no key exists (or the value is empty) then return the default value
	if csv == "" {
		return  defaultValue
	}

	//otherwise parse the value into a []string slice and then return it
	return strings.Split(csv, ",")
}


//The readInt() helper reads a string value from the query string and converts it to an
//integer before returning. if no matching key could be found it returns the provided
//default value. if the value could not be converted to an integer, then we record an
//error message in the provided validator instance
func (app *application) readInt(queryValues url.Values, key string, defaultValue int, validate *validator.Validator) int {
	//Extract the value from the query string
	value := queryValues.Get(key)

	//if no key exists (or the value is empty) then return the default value
	if value == ""{
		return defaultValue
	}

	//try to convert the value to an int. adds an error message to the validator
	//instance and return the default value
	i, err := strconv.Atoi(value)
	if err != nil {
		validate.AddError(key, "must be an integer value")
		return defaultValue
	}

	//otherwise return the converted integer value
	return  i
}


