package main

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

func seedDrinks() {
	for _, drink := range drinksList {
		PutItem(Drink{
			Id:    drink.Id,
			Name:  drink.Name,
			Price: drink.Price,
		})
	}
}
