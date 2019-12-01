package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/controllers"
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/services"
)


func TestRegisterHandler(t *testing.T) {
	config := services.ReadConfig()

	db, err := services.NewDB(config)
	if db == nil {
		t.Fatal("Did not connected to database")
	}
	env := controllers.EnvironmentUser{db}

	body := []byte(`{"email":"d_kokin@inbox.ru", "login":"dan_kokin", "password":"password"}`)

	request, err := http.NewRequest("GET", "/register", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(env.ResponseRegisterHandler)

	handler.ServeHTTP(rr, request)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}