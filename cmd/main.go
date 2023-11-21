package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	delivery "forum.bbilisbe/pkg/delivery/http"
	"forum.bbilisbe/pkg/repository"
	"forum.bbilisbe/pkg/usecase"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Parsing the runtime configuration settings for the application;
	addr := flag.String("addr", ":7070", "HTTP Network Address")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ldate)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := repository.SetUpDB("sqlite3", "forum.db")
	if err != nil {
		errorLog.Fatal(err)
	}

	postRepo := repository.NewSqlPostsRepository(db)
	userRepo := repository.NewSqlUsersRepository(db)
	postUse := usecase.NewPostUsecase(postRepo, userRepo)
	userUse := usecase.NewUserUsecase(postRepo, userRepo)
	router := delivery.NewPostHandler(postUse, userUse, infoLog, errorLog)

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// Establishing the dependencies for the handlers;
	srv := &http.Server{
		Addr:         *addr,
		ErrorLog:     errorLog,
		Handler:      router,
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Running the HTTP server
	infoLog.Printf("Starting server on... http://127.0.0.1%s \n", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
