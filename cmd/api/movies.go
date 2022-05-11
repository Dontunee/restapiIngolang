package main

import (
	"errors"
	"fmt"
	"greenlight.alexedwards.net/internal/data"
	"greenlight.alexedwards.net/internal/validator"
	"net/http"
	"strconv"
)

//add a listMoviesHandler for the GET "v1/movies" endpoint.
func (app *application) listMoviesHandler(writer http.ResponseWriter, request *http.Request){
	//To keep things consistent with our other handlers, we will define an input struct
	//to hold the expected values from the query string

	var input struct {
		Title string
		Genres []string
		data.Filters
	}

	//initialize  a new validator instance
	validate := validator.New()

	//call method to get the url.Values map containing the query string data
	queryValues := request.URL.Query()

	//Use our helpers to extract the title and genres query string values, falling back
	// to defaults of an empty string and empty slice respectively if they are not
	//provided by the client
	input.Title = app.readString(queryValues,"title", "")
	input.Genres = app.readCSV(queryValues, "genres", []string{})

	//Get the page and page_size query string values as integers. Notice that we
	//set the default page to 1 and default page_size to 20, and that we pass
	// the validator instance as the final argument here
	input.Filters.Page = app.readInt(queryValues, "page", 1, validate)
	input.Filters.PageSize = app.readInt(queryValues, "page_size", 20, validate)

	//Extract the sort query string value, falling back to "id" if it is not provided
	//by the client (which will imply a ascending sort on movie ID).
	input.Sort = app.readString(queryValues, "sort", "id")

	//Add the supported sort values for this endpoint to the sort safe list
	input.Filters.SortSafeList = []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"}

	//Execute the validation checks on the filters struct and send a response
	//containing the errors if necessary
	if data.ValidateFilters(validate, input.Filters); !validate.Valid(){
		app.failedValidationResponse(writer,request,validate.Errors)
		return
	}

	//Check the Validator instance for any errors and use the failedValidationResponse()
	//helper to send the client a response if necessary
	if !validate.Valid() {
		app.failedValidationResponse(writer,request,validate.Errors)
		return
	}

	movies, err := app.models.Movies.GetAll(input.Title, input.Genres, input.Filters)
	if err != nil {
		app.serverErrorResponse(writer, request,err)
		return
	}

	//send a JSON response containing the movie data
	err = app.writeJSON(writer,http.StatusOK, envelope{"movies":movies}, nil)
	if err != nil {
		app.serverErrorResponse(writer,request,err)
	}
}

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

	//Call the Insert() Method on our movies model, passing in a pointer to the
	//validated movie struct. This will create a record in the database and update the
	//movie struct with the generated information
	err = app.models.Movies.Insert(movie)
	if err != nil {
		app.serverErrorResponse(writer,request, err)
		return
	}

	//When sending a HTTP response, we want to include a Location header to let the
	//client know which URL they can find the newly created resource at. We make an empty
	// http.header map and then use the Set() method to add a new location header.
	//interpolating the generated ID for our new movie in the URL
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))


	//Write a JSON response with a 201 Created Status Coder, the movie data in the response body,
	//and the location header
	err = app.writeJSON(writer,http.StatusCreated, envelope{"movie":movie}, headers)
	if err != nil {
		app.serverErrorResponse(writer,request,err)
		return
	}
}

//Add a showMovieHandler for the GET "/v1/movies/:id" endpoint.
func (app *application) showMovieHandler(writer http.ResponseWriter, request *http.Request) {
	id, err := app.readIDParam(request)
	if err != nil {
		app.notFoundResponse(writer, request)
		return
	}

	//Call the Get() method to fetch the data for a specific movie. We also need to use the
	//errors.IS() function to check if it returns a data.ErrRecordNotFound error
	//in which case we send a 404 Not found response to the client
	movie, err := app.models.Movies.Get(id)
	if err != nil {
		if errors.Is(err,data.ErrRecordNotFound){
			app.notFoundResponse(writer,request)
		}else {
			app.serverErrorResponse(writer,request,err)
		}
		return
	}

	err = app.writeJSON(writer, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
		return
	}
}

//add an updateMovieHandler which updates the movie record, increases version of the item updated
func (app *application) updateMovieHandler(writer http.ResponseWriter, request *http.Request){
	//Extract the movie ID from the URL.
	id , err := app.readIDParam(request)
	if err != nil {
		app.notFoundResponse(writer,request)
		return
	}

	//Fetch the existing movie record from the database, send a 404 not found response to the client
	//if we couldn't find a matching record
	movie , err := app.models.Movies.Get(id)
	if err != nil {
		if errors.Is(err,data.ErrRecordNotFound){
			app.notFoundResponse(writer,request)
			return
		}else{
			app.serverErrorResponse(writer,request,err)
			return
		}
	}

	//if the request contains a X-Expected-Version header, verify that the movie
	//version in the database matches the expected version specified in the header
	if request.Header.Get("X-Expected-Version") !=  ""{
		if strconv.FormatInt(int64(movie.Version), 32) != request.Header.Get("X-Expected-Version"){
			app.editConflictResponse(writer,request)
			return
		}
	}

	//Declare an input struct to hold the expected data from the client
	var input struct{
		Title *string `json:"title"`
		Year   *int32  `json:"year"`
		Runtime *data.Runtime  `json:"runtime"`
		Genres []string  `json:"genres"`
	}

	//Read the JSON request body data into the input struct
	err = app.readJSON(writer,request,&input)
	if err != nil {
		app.badRequestResponse(writer,request,err)
		return
	}

	//Copy the values from the request body to the appropriate fields of the  movie record
	if input.Title != nil {
		movie.Title = *input.Title
	}
	if input.Year != nil {
		movie.Year = *input.Year
	}
	if input.Runtime != nil {
		movie.RunTime = *input.Runtime
	}
	if input.Genres != nil {
		movie.Genres = input.Genres
	}

	//Validate the updated movie record, sending the client a 422 unprocessable entity response
	//if any checks fail
	v := validator.New()

	if data.ValidateMovie(v,movie); !v.Valid() {
		app.failedValidationResponse(writer,request,v.Errors)
		return
	}

	//Pass the updated movie record to our new Update() method
	err = app.models.Movies.Update(movie)
	if err != nil {
		if errors.Is(err,data.ErrEditConflict){
			app.editConflictResponse(writer,request)
		}
		app.serverErrorResponse(writer,request,err)
		return
	}

	err = app.writeJSON(writer,http.StatusOK, envelope{"movie":movie},nil)
	if err != nil {
		app.serverErrorResponse(writer,request,err)
		return
	}
}

func (app *application) deleteMovieHandler(writer http.ResponseWriter, request *http.Request){
	//Extract the movie ID from the URL
	id , err := app.readIDParam(request)
	if err != nil {
		app.notFoundResponse(writer,request)
		return
	}

	//Delete the movie from the database, sending a 404 NotFound Response to the client if there isnt
	//a matching record
	err = app.models.Movies.Delete(id)
	if err != nil {
		if errors.Is(err,data.ErrRecordNotFound){
			app.notFoundResponse(writer,request)
		}else{
			app.serverErrorResponse(writer,request,err)
		}
		return
	}

	//Returns a 200 Ok status code along with a success message
	err = app.writeJSON(writer,http.StatusOK,envelope{"message": "movie successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(writer,request,err)
		return
	}
}
