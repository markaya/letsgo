package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/markaya/snippetbox/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

// NOTE: You can parse flag into pre existing var in memory

type config struct {
	addr      string
	staticDir string
}

type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	snippets       models.SnipetModelInterface
	users          models.UserModelInterface
	templateCache  map[string]*template.Template
	sessionManager *scs.SessionManager
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

	// NOTE: Template cahce
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	// NOTE: Session manager
	sessionManager := scs.New()
	sessionManager.Store = sqlite3store.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	// NOTE: Application
	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		sessionManager: sessionManager,
	}

	// NOTE: TLS
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		// You can set tls version here as well
		// NOTE: Enforce this so that TLS1.3 be used which enforces SameSite cookies.
		MinVersion: tls.VersionTLS13,
		// MaxVersion:       tls.VersionTLS13,
		// And also Cipher suits, to only use modern ones
		// CipherSuites: []uint16{
		// tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384
		// etc..
		// }
	}

	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		// Limit max header length
		// MaxHeaderBytes:
	}

	infoLog.Printf("Starting server on %s\n", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}
