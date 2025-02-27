package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/markaya/snippetbox/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

// NOTE: You can parse flag into pre existing var in memory

type config struct {
	addr      string
	staticDir string
}

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	snippets *models.SnippetModel
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, err
}

func main() {

	var cfg config
	flag.StringVar(&cfg.addr, "address", ":4000", "HTTP network addr")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	addr := flag.String("addr", ":4000", "Http network addr")

	dsn := flag.String("dsn", "./snippetbox.db?_busy_timeout=5000&_journal_mode=WAL", "Sqlite db string")

	flag.Parse()

	// NOTE: Loggers
	infoLog := log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile)

	// NOTE: Database
	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// NOTE: Application
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		snippets: &models.SnippetModel{DB: db},
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s\n", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
