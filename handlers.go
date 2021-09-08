package main

import (
	"net/http"
	"strconv"

	"github.com/GemTay/basic-go-app/forms"
	"github.com/GemTay/basic-go-app/web/templates"
	"github.com/pborman/uuid"
)

var drinksList = []*Drink{
	&Drink{
		Id:    "ea2ba76e-381e-4c3d-991d-ba8b82bf57f2",
		Name:  "Latte",
		Price: 2.45,
	},
	&Drink{
		Id:    "6c2eaae7-c3d2-4762-be83-f91c26b915c8",
		Name:  "Espresso",
		Price: 1.99,
	},
	&Drink{
		Id:    "30e9ae96-7dd8-41dd-9c49-e13d73ab23ac",
		Name:  "Cappuccino",
		Price: 2.55,
	},
}

func (app *application) GetAllDrinksHandler(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("Running the GET ALL drinks handler")

	drinks := app.GetAllItems()

	templates.Render(w, "get-all-drinks.gohtml", drinks)
}

func (app *application) GetDrinkHandler(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("Running the GET drink handler")

	uId, ok := r.URL.Query()["uId"]

	if !ok || len(uId[0]) < 1 {
		app.notFound(w)
		return
	}

	templates.Render(w, "get-drink.gohtml", app.GetItem(uId[0]))
}

func (app *application) AddDrinkHandler(w http.ResponseWriter, r *http.Request) {
	app.infoLog.Println("Running the ADD drink handler")

	if r.Method != http.MethodPost {
		templates.Render(w, "add-drink.gohtml", nil)
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.serverError(w, err)
	}

	data := forms.DrinkData{
		Name:  r.PostFormValue("drink-name"),
		Price: r.PostFormValue("drink-price"),
	}

	err = data.ValidateAddDrink()

	if err != nil {
		templates.Render(w, "add-drink.gohtml", err)
	} else {
		p, err := strconv.ParseFloat(r.FormValue("drink-price"), 64)
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
		}

		id := uuid.New()

		drink := Drink{
			Id:    id,
			Name:  r.FormValue("drink-name"),
			Price: p,
		}

		app.PutItem(drink)

		templates.Render(w, "add-drink-success.gohtml", nil)
	}
}

func (app *application) seedDrinks() {
	for _, drink := range drinksList {
		app.PutItem(Drink{
			Id:    drink.Id,
			Name:  drink.Name,
			Price: drink.Price,
		})
	}
}
