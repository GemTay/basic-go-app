package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	CreateTable()

	var drink Drink = Drink{
		Id:    1,
		Name:  "Americano",
		Price: 2.10,
	}

	PutItem(drink)

	drink.Name = "Latte"

	UpdateItem(drink)

	// logger
	l := log.New(os.Stdout, "basic-app", log.LstdFlags)

	// creating the serve mux router
	sm := mux.NewRouter()

	// creating some template data to be passed in
	tplData := GetItem(1)

	// parsing all the templates ready to be executed
	tpls := template.Must(template.ParseGlob("./web/templates/*"))

	sm.HandleFunc("/hello-world", func(w http.ResponseWriter, r *http.Request) {
		l.Println("Running the hello world handler")
		fmt.Println(GetItem(1))
		w.WriteHeader(http.StatusOK)

		err := tpls.ExecuteTemplate(w, "hello-world.html", tplData)
		if err != nil {
			log.Println("Error executing template :", err)
			return
		}
	})

	// setting up the http server
	s := &http.Server{
		Addr:         ":8080",
		Handler:      sm,
		ErrorLog:     l,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 15 * time.Minute,
	}

	// starting the server
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
		l.Println("Server is up and running!")
	}()

	// catching interrupt and kill signals to gracefully shutdown the server
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	l.Println("Recieve terminate, gracefully shutting down", sig)
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
