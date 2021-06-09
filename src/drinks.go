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

var tpls = template.Must(template.ParseGlob("./web/templates/*"))

func GetAllDrinks(w http.ResponseWriter, r *http.Request) {
	log.Println("Running the GET ALL drinks handler")

	drinks := GetAllItems()
	log.Println(drinks)

	err := tpls.ExecuteTemplate(w, "get-all-drinks.gohtml", drinks)
	if err != nil {
		log.Println("Error executing template :", err)
		return
	}
}

func GetDrink(w http.ResponseWriter, r *http.Request) {
	log.Println("Running the GET drink handler")

	id, convErr := strconv.Atoi(mux.Vars(r)["id"])
	if convErr != nil {
		log.Println("Error converting id from string to int :", convErr)
		return
	}

	err := tpls.ExecuteTemplate(w, "get-drink.gohtml", GetItem(id))
	if err != nil {
		log.Println("Error executing template :", err)
		return
	}
}

func AddDrink(w http.ResponseWriter, r *http.Request) {
	log.Println("Running the ADD drink handler")

	if r.Method != http.MethodPost {
		err := tpls.ExecuteTemplate(w, "add-drink.gohtml", nil)
		if err != nil {
			log.Println("Error executing template :", err)
			return
		}
		return
	}

	p, err := strconv.ParseFloat(r.FormValue("drink-price"), 64)
	if err != nil {
		log.Fatal(err)
	}

	drink := Drink{
		Id:    4, //needs fixing
		Name:  r.FormValue("drink-name"),
		Price: p,
	}

	PutItem(drink)

	e := tpls.ExecuteTemplate(w, "add-drink.gohtml", struct{ Success bool }{true})
	if e != nil {
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
