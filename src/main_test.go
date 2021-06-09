package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gorilla/mux"
	dynamock "github.com/gusaul/go-dynamock"
	"github.com/stretchr/testify/assert"
)

var mock *dynamock.DynaMock

func init() {
	Dyna = new(MyDynamo)
	Dyna.Db, mock = dynamock.New()
}

func TestGetName(t *testing.T) {
	expectKey := map[string]*dynamodb.AttributeValue{
		"id": {
			N: aws.String("1"),
		},
	}

	expectedResult := aws.String("jaka")
	result := dynamodb.GetItemOutput{
		Item: map[string]*dynamodb.AttributeValue{
			"name": {
				S: expectedResult,
			},
		},
	}

	//lets start dynamock in action
	mock.ExpectGetItem().ToTable("employee").WithKeys(expectKey).WillReturns(result)

	actualResult, _ := GetName("1")
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

func TestGetItem(t *testing.T) {
	r := GetItem(1)

	assert.IsType(t, Drink{}, r)
	assert.Equal(t, Drink{1, "Latte", 2.45}, r)
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
