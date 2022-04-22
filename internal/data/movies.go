package data

import (
	"greenlight.alexedwards.net/internal/validator"
	"time"
)

type Movie struct {
	ID        int64     `json:"id"`                    //unique integer for the movie
	CreatedAt time.Time `-`                            //Time stamp for when the movie is added to our database
	Title     string    `json:"title"`                 // Title of the movie
	Year      int32     `json:"year,omitempty,string"` //Movie release year
	RunTime   Runtime   `json:"runTime"`               //Move runtime in (minutes)
	Genres    []string  `json:"genres"`                //Slice of genres for the movie (romance, comedy, etc.)
	Version   int32     `json:"version"`               //The version number starts at 1 and will be incremented each time the movie
	//information is updated
}

func ValidateMovie(v *validator.Validator, movie *Movie) {
	//Use the Check() method to execute our validation checks. This will add the provided key and error message to the map
	//if the check does not evaluate to be true
	v.Check(movie.Title != "", "title", "required")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(movie.Year != 0, "year", "required")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(movie.RunTime != 0, "runtime", "required")
	v.Check(movie.RunTime > 0, "runtime", "must be a positive integer")

	v.Check(movie.Genres != nil, "genres", "required")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")

	//use the unique helper in the validator to check that all values in the input.Genres are unique
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")
}
