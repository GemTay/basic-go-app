package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type application struct {
	infoLog  *log.Logger
	errorLog *log.Logger
}

func main() {

	ConfigureDynamoDB()

	// This takes three parameters: the destination to write the logs to (os.Stdout), a string
	// prefix for message (INFO followed by a tab), and flags to indicate what
	// additional information to include (local date and time). Note that the flags
	// are joined using the bitwise OR operator |.
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	// Create a logger for writing error messages in the same way, but use stderr as
	// the destination and use the log.Lshortfile flag to include the relevant
	// file name and line number.
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Initialize a new instance of application containing the dependencies.
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	// app.seedDrinks()

	// setting up the http server
	s := &http.Server{
		Addr:         ":8080",
		Handler:      app.routes(),
		ErrorLog:     errorLog,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 15 * time.Minute,
	}

	// starting the server
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			errorLog.Fatal(err)
		}
		infoLog.Println("Server is up and running!")
	}()

	// catching interrupt and kill signals to gracefully shutdown the server
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	infoLog.Println("Recieve terminate, gracefully shutting down", sig)
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
