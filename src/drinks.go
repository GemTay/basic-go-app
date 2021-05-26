package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

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

func GetDrink(w http.ResponseWriter, r *http.Request) {
	log.Println("Running the GET drink handler")

	id, convErr := strconv.Atoi(mux.Vars(r)["id"])
	if convErr != nil {
		log.Println("Error converting id from string to int :", convErr)
		return
	}

	tpls := template.Must(template.ParseGlob("./web/templates/*"))

	err := tpls.ExecuteTemplate(w, "get-drink.html", GetItem(id))
	if err != nil {
		log.Println("Error executing template :", err)
		return
	}
}

func seedDrinks() {
	for _, drink := range drinksList {
		PutItem(Drink{
			Id:    drink.Id,
			Name:  drink.Name,
			Price: drink.Price,
		})
	}
}
