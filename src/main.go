package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	// CreateTable()

	var drinksList = []*Drink{
		&Drink{
			Id:    1,
			Name:  "Latte",
			Price: 2.45,
		},
		&Drink{
			Id:    2,
			Name:  "Espresso",
			Price: 1.99,
		},
		&Drink{
			Id:    3,
			Name:  "Cappuccino",
			Price: 2.55,
		},
	}

	for _, drink := range drinksList {
		PutItem(Drink{
			Id:    drink.Id,
			Name:  drink.Name,
			Price: drink.Price,
		})
	}

	// logger
	l := log.New(os.Stdout, "basic-app", log.LstdFlags)

	// creating the serve mux router
	sm := mux.NewRouter()

	// parsing all the templates ready to be executed
	tpls := template.Must(template.ParseGlob("./web/templates/*"))

	sm.HandleFunc("/drinks/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		l.Println("Running the get drink handler")

		id, convErr := strconv.Atoi(mux.Vars(r)["id"])
		if convErr != nil {
			l.Println("Error converting id from string to int :", convErr)
			return
		}

		err := tpls.ExecuteTemplate(w, "get-drink.html", GetItem(id))
		if err != nil {
			l.Println("Error executing template :", err)
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
