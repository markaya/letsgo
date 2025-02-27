package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

// NOTE: You can parse flag into pre existing var in memory

type config struct {
	addr      string
	staticDir string
}

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	var cfg config
	flag.StringVar(&cfg.addr, "address", ":4000", "HTTP network addr")
	flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")
	addr := flag.String("addr", ":4000", "Http network addr")
	flag.Parse()

	// NOTE: Loggers
	infoLog := log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile)

	// NOTE: Application
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s\n", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
