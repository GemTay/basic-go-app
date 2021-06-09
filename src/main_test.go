package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gorilla/mux"
	dynamock "github.com/gusaul/go-dynamock"
)

var mock *dynamock.DynaMock

func init() {
	Dyna = new(MyDynamo)
	Dyna.Db, mock = dynamock.New()
}

func TestGetItem(t *testing.T) {
	expectedResult := Drink{
		Id:    1,
		Name:  "Latte",
		Price: 2.45,
	}

	expectKey := map[string]*dynamodb.AttributeValue{
		"Id": {
			N: aws.String(strconv.Itoa(expectedResult.Id)),
		},
	}

	result := dynamodb.GetItemOutput{
		Item: map[string]*dynamodb.AttributeValue{
			"Id": {
				N: aws.String(strconv.Itoa(expectedResult.Id)),
			},
			"Name": {
				S: &expectedResult.Name,
			},
			"Price": {
				N: aws.String(strconv.FormatFloat(expectedResult.Price, 'f', 2, 64)),
			},
		},
	}

	mock.ExpectGetItem().
		ToTable("drinks").
		WithKeys(expectKey).
		WillReturns(result)

	actualResult := GetItem(1)
	if actualResult != expectedResult {
		t.Errorf("Test Fail")
	}
}

func TestGetDrinkHandler(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/drinks/", nil)

	r = mux.SetURLVars(r, map[string]string{
		"id": "1",
	})

	GetDrink(w, r)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", res.StatusCode)
	}
}

func TestAddDrink(t *testing.T) {
	form := url.Values{}
	form.Add("drink-name", "Tea")
	form.Add("drink-price", "1.50")

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/add-drink", strings.NewReader(form.Encode()))

	AddDrink(w, r)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", res.StatusCode)
	}
}
