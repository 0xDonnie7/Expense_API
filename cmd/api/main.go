package main

import (
	"flag"
	"log"
	"os"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxIdleConns int
		maxOpenConns int
		maxIdleTime  string
	}
}

type application struct {
	config config
}

func main() {
	var cfg config
	flag.IntVar(&cfg.port, "port", 8080, "port number")
	flag.StringVar(&cfg.env, "env", "development", "development|staging|production")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("EXPENSE_DB_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open conns")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max idle time")

	flag.Parse()

	app := &application{
		config: cfg,
	}

	err := app.server()
	if err != nil {
		log.Fatal(err)
	}

}
