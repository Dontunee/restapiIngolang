package data

import "time"

type Movies struct {
	ID        int64     //unique integer for the movie
	CreatedAt time.Time //Time stamp for when the movie is added to our database
	Title     string    // Title of the movie
	Year      int32     //Movie release year
	RunTime   int32     //Move runtime in (minutes)
	Genres    []string  //Slice of genres for the movie (romance, comedy, etc.)
	Version   int32     //The version number starts at 1 and will be incremented each time the movie
	//information is updated
}
