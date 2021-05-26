package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

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
