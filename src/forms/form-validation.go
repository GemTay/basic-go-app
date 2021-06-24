package forms

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type DrinkData struct {
	Name  string
	Price string
}

type Address struct {
	Street string
	City   string
	State  string
	Zip    string
}

func (d DrinkData) ValidateAddDrink() error {
	return validation.ValidateStruct(&d,
		// Name cannot be empty, and the length must between 2 and 50
		validation.Field(&d.Name,
			validation.Required.Error("Enter a name"),
			validation.Length(2, 15).Error("The name must be between 2 and 15 characters long")),
		// Price cannot be empty, and must be a number
		validation.Field(&d.Price,
			validation.Required.Error("Enter a price"),
			is.Float.Error("The price must be a number")),
	)
}

func (a Address) ValidateAddress() error {
	//val := reflect.Indirect(reflect.ValueOf(a))

	return validation.ValidateStruct(&a,
		// Street cannot be empty, and the length must between 5 and 50
		validation.Field(&a.Street,
			validation.Required,
			validation.Length(5, 50)),
		// validation.Length(5, 50).Error(ValidationMessages(val.Type().Field(3).Name, "LengthTooLong"))),
		// City cannot be empty, and the length must between 5 and 50
		validation.Field(&a.City,
			validation.Required,
			validation.Length(5, 50)),
		// State cannot be empty, and must be a string consisting of two letters in upper case
		validation.Field(&a.State, validation.Required, validation.Match(regexp.MustCompile("^[A-Z]{2}$"))),
		// State cannot be empty, and must be a string consisting of five digits
		validation.Field(&a.Zip, validation.Required, validation.Match(regexp.MustCompile("^[0-9]{5}$"))),
	)
}
