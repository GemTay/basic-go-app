package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	l := log.New(os.Stdout, "basic-app", log.LstdFlags)

	sm := mux.NewRouter()

	type templateData struct {
		Title string
		Date  string
	}

	tplData := templateData{
		"Hello World",
		time.Now().Format("Mon Jan 2 15:04:05 MST 2006"),
	}

	tpls := template.Must(template.ParseGlob("./web/templates/*"))

	sm.HandleFunc("/hello-world", func(w http.ResponseWriter, r *http.Request) {
		l.Println("Running the hello world handler")
		w.WriteHeader(http.StatusOK)

		err := tpls.ExecuteTemplate(w, "hello-world.html", tplData)
		if err != nil {
			log.Println("Error executing template :", err)
			return
		}
	})

	s := &http.Server{
		Addr:         ":8080",
		Handler:      sm,
		ErrorLog:     l,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 15 * time.Minute,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
		l.Println("Server is up and running!")
	}()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	l.Println("Recieve terminate, gracefully shutting down", sig)
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)
}
