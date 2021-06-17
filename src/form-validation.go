package main

import (
	"strconv"
	"strings"
)

type DrinkData struct {
	Name   string
	Price  string
	Errors map[string]string
}

func (msg *DrinkData) Validate() bool {
	msg.Errors = make(map[string]string)

	if strings.TrimSpace(msg.Name) == "" {
		msg.Errors["Name"] = "Please enter a name"
	}

	if strings.TrimSpace(msg.Price) == "" {
		msg.Errors["Price"] = "Please enter a price"
	} else {
		_, err := strconv.ParseFloat(msg.Price, 64)

		if err != nil {
			msg.Errors["Price"] = "Price must be a decimal number"
		}
	}

	return len(msg.Errors) == 0
}
