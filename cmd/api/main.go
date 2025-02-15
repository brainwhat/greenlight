package main

import (
	"bufio"
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	config config
	logger *slog.Logger
}

func main() {

	var cfg config

	SetENV()

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "dev", "Current environment (dev/stage/prod")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("dsn"), "PostgreSQL DSN")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	logger.Info("database connection pool established")

	app := application{
		config: cfg,
		logger: logger,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "addr", srv.Addr, "env", cfg.env)

	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	// db conns are established lazily (only when they are first called)
	// so we create context with 5 second timeout and establish a connection
	// if it isn't established within 5 seconds, close connection and return err
	//TODO: learn about contexts
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

// Read .env and set config variables (so far only dsn)
func SetENV() {

	envFile, err := os.Open("./.env")
	if err != nil {
		log.Fatalln(err)
	}
	defer envFile.Close()

	scanner := bufio.NewScanner(envFile)

	for scanner.Scan() {
		name, value, _ := strings.Cut(scanner.Text(), "=")
		os.Setenv(name, value)
	}

	// Check if there any errors during scanning
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
