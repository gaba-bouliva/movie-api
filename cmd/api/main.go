package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

type config struct {
	port 				int
	env 				string
	db 					struct{
		dsn 						string
		maxIdleConns		int
		maxOpenConns    int
		maxIdleTime			string
	}
}

type application 	struct {
	config 		config
	logger 		*log.Logger
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 8081, "Server API port")
	flag.StringVar(&cfg.env, "env", "development", "[ development | staging | production ]")

	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("MOVIE-API-DB-DSN"), "database connection string")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "database max idle connections")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "database max open connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "database max open connections")

	logger := log.New(os.Stdout, "", log.Ldate | log.Ltime | log.Lshortfile)

	app := application{
		config: cfg,
		logger: logger,
	}

	if err := app.run(); err != nil {
		logger.Fatal(err)
	} 

} 

func (app *application) run() error {
	srv := http.Server{
		Addr: fmt.Sprintf(":%d", app.config.port),
		Handler: app.routes(),
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.logger.Printf("%s server starting on port: %d", app.config.env, app.config.port)
	return srv.ListenAndServe()
}