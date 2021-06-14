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

func TestPutItem(t *testing.T) {
	dr := Drink{
		Id:    1,
		Name:  "Latte",
		Price: 2.45,
	}

	result := dynamodb.PutItemOutput{
		Attributes: map[string]*dynamodb.AttributeValue{
			"Id": {
				N: aws.String(strconv.Itoa(dr.Id)),
			},
			"Name": {
				S: &dr.Name,
			},
			"Price": {
				N: aws.String(strconv.FormatFloat(dr.Price, 'f', 2, 64)),
			},
		},
	}

	mock.ExpectPutItem().
		ToTable("drinks").
		WithItems(result.Attributes).
		WillReturns(result)

	PutItem(dr)
	if actualResult != result {
		t.Errorf("Test Fail")
	}
}

func TestGetAllItems(t *testing.T) {
	expectedResult := []Drink{
		Drink{
			Id:    1,
			Name:  "Latte",
			Price: 2.45,
		},
		Drink{
			Id:    2,
			Name:  "Espresso",
			Price: 1.99,
		},
	}

	result := dynamodb.ScanOutput{
		Items: []map[string]*dynamodb.AttributeValue{
			{
				"Id": {
					N: aws.String(strconv.Itoa(expectedResult[0].Id)),
				},
				"Name": {
					S: aws.String(expectedResult[0].Name),
				},
				"Price": {
					N: aws.String(strconv.FormatFloat(expectedResult[0].Price, 'f', 2, 64)),
				},
			},
			{
				"Id": {
					N: aws.String(strconv.Itoa(expectedResult[1].Id)),
				},
				"Name": {
					S: &expectedResult[1].Name,
				},
				"Price": {
					N: aws.String(strconv.FormatFloat(expectedResult[1].Price, 'f', 2, 64)),
				},
			},
		},
	}

	mock.ExpectScan().Table("drinks").WillReturns(result)

	actualResult := GetAllItems()

	for i, _ := range actualResult {
		if actualResult[i] != expectedResult[i] {
			t.Errorf("Drink with index %d did not match expected", i)
		}
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
