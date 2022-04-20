package data

import "time"

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
