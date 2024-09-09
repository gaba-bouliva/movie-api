package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gaba-bouliva/movie-api/internal/data"
	_ "github.com/lib/pq"
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
	models    data.Models
	version    string
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 8081, "Server API port")
	flag.StringVar(&cfg.env, "env", "development", "[ development | staging | production ]")

	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("MOVIE_API_DB_DSN"), "database connection string")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "database max idle connections")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "database max open connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "database max open connections")

	logger := log.New(os.Stdout, "", log.Ldate | log.Ltime | log.Lshortfile)

	db, err := OpenDB(cfg)
	if err != nil {
		logger.Fatal("error connecting to db\n",err)
	}
	defer db.Close()

	logger.Println("db connection successfully established...")

	app := application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
		version: version,
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

func OpenDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	maxIdleTime, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(maxIdleTime)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	db.SetMaxOpenConns(cfg.db.maxOpenConns)

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil

}