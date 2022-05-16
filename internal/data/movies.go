package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"greenlight.alexedwards.net/internal/validator"
	"time"
)

// ErrRecordNotFound Define a custom ErrRecordNotFound. We'll return this from our Get() method when
//looking aup a movie that doesnt exist in our database
var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict = errors.New("edit conflict")
)


//Define a new Metadata struct for holding the pagination metadata

// MovieModel Defines a MovieModel struct type which wraps a swl.DB connection pool
type MovieModel struct {
	DB *sql.DB
}

// Models Create a Models struct which wraps the MovieModel. We will add models to this,
//like a UserModel and PermissionModel, as our build continues
type Models struct {
	Movies MovieModel
}


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

// GetAll Create a new GetAll() method which returns a slice of movies. Although we are not
//using them right now, we have set this up to accept the various filter parameters as
// arguments
func (movieModel MovieModel) GetAll(title string, genres []string, filters Filters) ([]*Movie, Metadata, error){
	//Construct the SQL query to retrieve all movie records
	query := fmt.Sprintf(`SELECT count(*) OVER(), id, created_at, title, year, runtime, genres, version 
				FROM movies 
				WHERE (to_tsvector('simple',title)  @@ plainto_tsquery('simple',$1) OR $1 = '')
				AND (genres @> $2 OR $2 = '{}')
				ORDER BY %s %s, id ASC
				LIMIT $3 OFFSET $4`,filters.sortColumn(),filters.sortDirection())

	//Create a context with a 3-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	//Use QueryContext() to execute the query. This returns a sql.Rows result set
	//containing the result
	rows,err := movieModel.DB.QueryContext(ctx,query,title,pq.Array(genres),filters.limit(),filters.offset())
	if err != nil {
		return nil,Metadata{}, err
	}

	//importantly , defer a call to rows.Close() to ensure that the result set is cosed before GetAll() returns
	defer rows.Close()

	// declare variables
	totalRecords := 0
	movies := []*Movie{}

	//Use rows.Next() to iterate through the rows in the result set
	for rows.Next() {
		//initialize an empty Movie struct to hold the data for an individual movie
		var movie Movie

		//scan the values from the row into the movie struct. Again, note that we are
		//using the pq.Array() adapter on the genres field here.
		err := rows.Scan(
			&totalRecords,
			&movie.ID,
			&movie.CreatedAt,
			&movie.Title,
			&movie.Year,
			&movie.RunTime,
			pq.Array(&movie.Genres),
			&movie.Version,
			)

		if err != nil {
			return nil,Metadata{}, err
		}

		//Add the Movie struct to the slice
		movies = append(movies,&movie)
	}

	//when the rows.Next() loop has finished, call rows.Err() to retrieve any error
	//that was encountered during the iteration
	if err = rows.Err(); err != nil {
		return nil,Metadata{}, err
	}

	//Generate a metadata struct, passing in the total record count and pagination
	//parameters from the client
	metaData := calculateMetadata(totalRecords,filters.Page,filters.PageSize)

	//if everything went ok , then return the slice of movies
	return movies,metaData, nil
}

// Insert Placeholder method for inserting a new record in the movies table
// accepts a pointer to a movie struct, which should contain the data for the
// new record
func (movieModel MovieModel) Insert (movie *Movie) error{
	//Define the SQL query for inserting a new record in the movies
	//table and returning the system-generated data.
	query := `INSERT INTO movies (title,year, runtime ,genres)
	         VALUES ($1,$2,$3,$4)
	         RETURNING id, created_at, version`

	//Create an args slice containing the values for the placeholder parameters from
	//the movie struct. Declaring the slice immediately next to our SQL query makes it
	//readable and clear its usage
	args := []interface{}{movie.Title, movie.Year, movie.RunTime, pq.Array(movie.Genres)}

	ctx , cancel := context.WithTimeout(context.Background(),3*time.Second)
	defer  cancel()

	//Use the QueryRow() method to execute the query on our connection pool
	//passing in the args slice as a variadic parameter and scanning the system-generated
	//id, created_at amd version values into the movie struct
	return movieModel.DB.QueryRowContext(ctx, query,args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

// Get Placeholder method for fetching a specific record from the movies table
func (movieModel *MovieModel) Get(id int64) (*Movie, error){

	//PostgreSQL bigserial type starts auto incrementing at 1 by default, therefore no movies will ID values
	//less than that, its best avoid unnecessary database call and return an error
	if id < 1{
		return nil,ErrRecordNotFound
	}

	//Define the SQL query for retrieving the movie data
	query := `SELECT id,created_at,title,year,runtime,genres,version FROM movies WHERE id = $1`

	//Declare a Movie struct to hold the data returned by the query
	var movie Movie

	//use the context.WithTimeout() function to create a context.context which carries
	//a 3second deadline. Note that we are using the empty context.Background as the
	//parent context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	//defer to make sure to cancel the context before the Get() method returns
	defer cancel()

	//Execute the query using the QueryRow() method, passing in the provided id value as a placeholder parameter
	//and scan the response data into the fields of the Movie struct. Importantly, it is require dto scan target
	//for the genres column using the pq.Array() adapter function again.
	err := movieModel.DB.QueryRowContext(ctx, query,id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.RunTime,
		pq.Array(&movie.Genres),
		&movie.Version,
		)

	if err != nil {
		if errors.Is(err,sql.ErrNoRows){
			return nil,ErrRecordNotFound
		}else {
			return nil,err
		}
	}

	//otherwise return pointer to the movie struct
	return &movie, nil
}

//Update Placeholder method for updating a specific record in the movies table
func (movieModel MovieModel) Update (movie *Movie) error{
	//Declare the SQL query for updating the record and returning the new version number
	//Use version to apply an optimistic solution for data race condition
	query := `UPDATE movies SET title=$1, year= $2 , runtime = $3, genres = $4, version = version + 1
			  WHERE id = $5 and version = $6
	          RETURNING version`

	//Create an args slice containing the values for the placeholder parameters
	args := []interface{}{movie.Title, movie.Year, movie.RunTime,pq.Array(movie.Genres),movie.ID,movie.Version}

	ctx , cancel := context.WithTimeout(context.Background(),3*time.Second)
	defer  cancel()
	//Use the QueryROW() method to execute the query, passing in the args slice as a variadic parameter
	//and scanning the new version value into the movie struct
	return movieModel.DB.QueryRowContext(ctx,query,args...).Scan(&movie.Version)
}

//Delete Placeholder method for deleting a specific record from the movies table
func (movieModel *MovieModel) Delete(id int64) error {
	// return a ErrRecordNotFound error if the movie ID is less than 1
	if id < 1 {
		return ErrRecordNotFound
	}

	//construct the SQL query to delete the record
	query := `DELETE FROM movies WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	//Execute the SQL query using the Exec() method, passing in the id variable as the value
	//for the placeholder parameter. The Exec()  method returns a sql.Result object
	result , err := movieModel.DB.ExecContext(ctx,query,id)
	if err != nil {
		return err
	}

	//Call the RowsAffected() method on the sql.Result object to get the number of rows affected by the query
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	//if no rows were affected, we know that the movies didnt contain a record
	//with the provided ID at the moment we tried to delete it.in that case we
	//return an ErrRecordNotFound error.
	if rowsAffected == 0 {
		return  ErrRecordNotFound
	}

	return nil
}

// NewModels A New() method which returns a models struct containing initialized MovieModel.
func NewModels( db *sql.DB) Models {
	return Models{
		Movies: MovieModel{DB: db},
	}
}