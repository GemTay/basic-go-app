package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/GemTay/basic-go-app/web/templates"
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

func GetAllDrinks(w http.ResponseWriter, r *http.Request) {
	log.Println("Running the GET ALL drinks handler")

	drinks := GetAllItems()

	templates.Render(w, "get-all-drinks.gohtml", drinks)
}

func GetDrink(w http.ResponseWriter, r *http.Request) {
	log.Println("Running the GET drink handler")

	id, convErr := strconv.Atoi(mux.Vars(r)["id"])
	if convErr != nil {
		log.Println("Error converting id from string to int :", convErr)
		return
	}

	templates.Render(w, "get-drink.gohtml", GetItem(id))
}

func AddDrink(w http.ResponseWriter, r *http.Request) {
	log.Println("Running the ADD drink handler")

	if r.Method != http.MethodPost {
		templates.Render(w, "add-drink.gohtml", nil)
		return
	}

	p, err := strconv.ParseFloat(r.FormValue("drink-price"), 64)
	if err != nil {
		log.Fatal(err)
	}

	drink := Drink{
		Id:    5, //needs fixing
		Name:  r.FormValue("drink-name"),
		Price: p,
	}

	PutItem(drink)

	templates.Render(w, "add-drink.gohtml", struct{ Success bool }{true})
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
