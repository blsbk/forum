package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	dbase "forum.bbilisbe/internal/db"
	"forum.bbilisbe/models"
	_ "github.com/mattn/go-sqlite3"
)

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	posts         *models.PostModel
	users         *models.UserModel
	templateCache map[string]*template.Template
}

func main() {
	// Parsing the runtime configuration settings for the application;
	addr := flag.String("addr", ":8080", "HTTP Network Address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ldate)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := dbase.OpenDB()
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	if err = dbase.Exec(db); err != nil {
		errorLog.Fatal(err)
	}

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		posts:         &models.PostModel{DB: db},
		users:         &models.UserModel{DB: db},
		templateCache: templateCache,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// Establishing the dependencies for the handlers;
	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Running the HTTP server
	infoLog.Printf("Starting server on... https://127.0.0.1%s \n", *addr)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}
