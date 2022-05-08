package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"greenlight.alexedwards.net/internal/data"
	"log"
	"net/http"
	"os"
	"time"

	/*Import the pq driver so that it can register itself with the database/sql package.
	Note that we alisa this import to the blank identifier, to stop the Go compiler
	complaining that the package is not being used*/
	_ "github.com/lib/pq"
)

// application version number
const version = "1.0.0"

//config struct to hold all configuration settings for application
type config struct {
	port int
	env  string
	db   db
}
type db struct {
	dsn string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime string
}

// application struct to hold the dependencies of our http handlers,helpers and middleware
type application struct {
	config config
	logger *log.Logger
	models data.Models
}

func main() {
	//declare an instance of cfg struct
	var cfg config

	//read the value of port and env command-line flags into the cfg struct ,set the default for each flag
	//if not specified
	flag.IntVar(&cfg.port, "port", 8000, "API server port")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connection")
	flag.IntVar(&cfg.db.maxIdleConns,"db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	//Read the dsn value from the db-dsn command-line flag into the config struct. We
	//default to using our development DSN if no flag is provided
	flag.StringVar(&cfg.db.dsn, "db-dsn","postgres://greenlight:pa55word@localhost/greenlight?sslmode=disable" , "PostgreSQL DSN")
	flag.StringVar(&cfg.db.maxIdleTime,"db-max-idle-time", "15m", "PostgreSQL max connection idle time")
	flag.Parse()

	/*cfg.port = os.Getenv("port")
	cfg.db.dsn = os.Getenv("greenlight_db_dsn")
	cfg.env = os.Getenv("env") */

	//Initialize a new logger which writes messages to the standard out stream,
	//prefixed with the current date and time
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	//call the OpenDB() helper function to create connection pool,
	//passing the config struct. if this returns a error ,log and exit the application immediately

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	//Defer a call to db.Close() so that the connection pool is closed before the
	//main() function exits
	defer db.Close()

	//Also log a message  to say that the connection pool has been successfully established
	logger.Println("database connection pool established")

	//Declare an instance of the application struct, containing the cfg struct and logger
	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	//Declare an HTTP server with timeout settings, which listens on the port provided in the config
	//struct and uses the serve mux we created above as the handler
	server := &http.Server{
		Addr:        fmt.Sprintf(":%d",cfg.port) ,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	//Start the server
	logger.Printf("starting %s server on port %s", cfg.env, server.Addr)
	err = server.ListenAndServe()
	logger.Fatal(err)

}

//The openDB() function returns a sql.DB connection pool
func openDB(cfg config) (*sql.DB, error){
	//Use sql.Open() to create an empty connection pool, using the DSN from the config
	//struct
	db , err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	//Set the maximum number of open (in-use + idle)  connections in the pool. Note that
	//passing a value less than or equal to 0 will mean there is no limit
	db.SetMaxOpenConns(cfg.db.maxOpenConns)

	//Set the maximum number of idle connections in the pool. Note that
	//passing a value less than or equal to 0 will mean there is no limit
	db.SetMaxIdleConns(cfg.db.maxIdleConns)


	//Create a context with a 5 second timeout deadline
	ctx, cancel := context.WithTimeout(context.Background(),5 * time.Second)
	defer cancel()

	//Use PingContext() to establish a new connection to the database, , passing in the context we created above
	// as a parameter. if the connection couldn't be established within the 5 second deadline, then this will return an
	//error
	err = db.PingContext(ctx)
	if err != nil {
		return nil,err
	}

	//return the sql.DB connection pool
	return db,nil
}


