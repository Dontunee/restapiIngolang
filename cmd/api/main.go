package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// application version number
const version = "1.0.0"

//config struct to hold all configuration settings for application
type config struct {
	port int
	env  string
}

// application struct to hold the dependencies of our http handlers,helpers and middleware
type application struct {
	config config
	logger *log.Logger
}

func main() {
	//declare an instance of cfg struct
	var cfg config

	// read the value of port and env command-line flags into the cfg struct ,set the default for each flag
	//if not specified
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	//Initialize a new logger which writes messages to the standard out stream,
	//prefixed with the current date and time
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	//Declare an instance of the application struct, containing the cfg struct and logger
	app := &application{
		config: cfg,
		logger: logger,
	}

	//Declare an HTTP server with timeout settings, which listens on the port provided in the config
	//struct and uses the servemux we created above as the handler
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	//Start the server
	logger.Printf("starting %s server on port %s", cfg.env, server.Addr)
	err := server.ListenAndServe()
	logger.Fatal(err)

}
