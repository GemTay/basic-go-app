package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/gorilla/mux"
)

type MyDynamo struct {
	Db dynamodbiface.DynamoDBAPI
}

var Dyna *MyDynamo

func ConfigureDynamoDB() {
	Dyna = new(MyDynamo)
	awsSession, _ := session.NewSession(&aws.Config{Region: aws.String("eu-west-1")})
	svc := dynamodb.New(awsSession)
	Dyna.Db = dynamodbiface.DynamoDBAPI(svc)
}

func main() {
	ConfigureDynamoDB()
	seedDrinks()

	// logger
	l := log.New(os.Stdout, "basic-app", log.LstdFlags)

	// creating the serve mux router
	sm := mux.NewRouter()
	sm.HandleFunc("/all-drinks", GetAllDrinks).Methods("GET")
	sm.HandleFunc("/drinks/{id:[0-9]+}", GetDrink).Methods("GET")
	sm.HandleFunc("/add-drink", AddDrink).Methods("GET", "POST")

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
