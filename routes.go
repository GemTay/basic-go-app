package main

import "net/http"

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/all-drinks", app.GetAllDrinksHandler)
	mux.HandleFunc("/drinks/", app.GetDrinkHandler)
	mux.HandleFunc("/add-drink", app.AddDrinkHandler)

	return mux
}
