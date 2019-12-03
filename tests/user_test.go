package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/controllers"
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/services"
)

type TestData struct {
	Body []byte
	Method string
	URL string
	Status int
}

func NewTestData(body []byte, method string, url string, status int) TestData {
	return TestData{Body: body, Method: method, URL: url, Status: status}
}

var g_Data = make([]TestData, 0)

func init() {
	// Case 1. Correct request, code 200
	body1 := []byte(`{
		"email":"d_kokin@inbox.ru",
		"login":"dan_kokin",
		"password":"password"}`)

	// Case 2. Invalid body. Code 400
	body2 := []byte(`{
		"email":d_kokin2@inbox.ru",
		"login":"dan_kokin2",
		"password":"password2"}`)

	// Case 3. Invalid login
	body3 := []byte(`{
		"email":"d_kokin3@inbox.ru",
		"login":"",
		"password":"password"}`)

	// Case 4. Invalid password (quantity < 6)
	body4 := []byte(`{
		"email":"d_kokin3@inbox.ru",
		"login":"testuser4",
		"password":"pass"}`)

	// Case 5. Reinsert by mail
	body5 := []byte(`{
		"email":"d_kokin@inbox.ru",
		"login":"dan_kokin",
		"password":"password"}`)

	// Case 6. Reinsert by login
	body6 := []byte(`{
		"email":"d_kok1n@inbox.ru",
		"login":"dan_kokin",
		"password":"password"}`)

	user1 := NewTestData(body1, "POST", "/register", http.StatusOK)
	user2 := NewTestData(body2, "POST", "/register", http.StatusBadRequest)
	user3 := NewTestData(body3, "POST", "/register", http.StatusBadRequest)
	user4 := NewTestData(body4, "POST", "/register", http.StatusBadRequest)
	user5 := NewTestData(body5, "POST", "/register", http.StatusInternalServerError)
	user6 := NewTestData(body6, "POST", "/register", http.StatusInternalServerError)

	g_Data = append(g_Data, user1, user2, user3, user4, user5, user6)

}

func CheckRequestStatus(data TestData, handlerFunc http.HandlerFunc, t *testing.T) {
	request, err:= http.NewRequest(data.Method, data.URL, bytes.NewBuffer(data.Body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlerFunc)
	handler.ServeHTTP(rr, request)

	if status := rr.Code; status != data.Status {
		t.Log(string(data.Body))
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, data.Status)
	}
}

func TestUserHandlers(t *testing.T) {
	config := services.ReadConfig()

	db, err := services.NewDB(config)
	if err != nil {
		t.Fatal("Did not connected to database")
	}
	env := controllers.EnvironmentUser{db}

	for index, value := range g_Data {
		t.Log("Case № ", index)
		CheckRequestStatus(value, env.ResponseRegisterHandler, t)
		t.Log(" Success\n")
	}
}

