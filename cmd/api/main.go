package main

import (
	"database/sql"
	"flag"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/0xdonnie7/Expense_API/internal/database"
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
	logger *slog.Logger
	db     *database.Queries
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

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := OpenDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	app := &application{
		config: cfg,
		logger: logger,
		db:     database.New(db),
	}

	err = app.server()
	if err != nil {
		log.Fatal(err)
	}

}

func OpenDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
