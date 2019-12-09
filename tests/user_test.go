package tests

import (
	"bytes"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/controllers"
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/services"
)



var gData = make([]RequestData, 0)

type RequestData struct {
	Id int
	Body []byte
	Method string
	URL string
	Status int
	Cookie []*http.Cookie
}

func init() {
	// Case 1. Correct request, code 200
	body1 := []byte(`{
		"email":"d_kokn@inbox.ru",
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
		"email":"d_kokn@inbox.ru",
		"login":"dan_kokin",
		"password":"password"}`)

	// Case 6. Reinsert by login
	body6 := []byte(`{
		"email":"d_kok1n@inbox.ru",
		"login":"dan_kokin",
		"password":"password"}`)

	user1 := RequestData{Id: 1, Body: body1, Method: "POST", URL: "/register",
		Status: http.StatusOK, Cookie: nil}
	user2 := RequestData{Id: 2, Body: body2, Method: "POST", URL: "/register",
		Status: http.StatusBadRequest, Cookie: nil}
	user3 := RequestData{Id: 3, Body: body3, Method: "POST", URL: "/register",
		Status: http.StatusBadRequest, Cookie: nil}
	user4 := RequestData{Id: 4, Body: body4, Method: "POST", URL: "/register",
		Status: http.StatusBadRequest, Cookie: nil}
	user5 := RequestData{Id: 5, Body: body5, Method: "POST", URL: "/register",
		Status: http.StatusInternalServerError, Cookie: nil}
	user6 := RequestData{Id: 6, Body: body6, Method: "POST", URL: "/register",
		Status: http.StatusInternalServerError, Cookie: nil}

	gData = append(gData, user1, user2, user3, user4, user5, user6)
}

func CheckStatus(data RequestData, handlerFunc http.HandlerFunc, t *testing.T) []*http.Cookie {
	request, err:= http.NewRequest(data.Method, data.URL, bytes.NewBuffer(data.Body))
	if err != nil {
		t.Fatal(err)
	}

	id := strconv.Itoa(data.Id)

	request = mux.SetURLVars(request, map[string]string{"id" : id})

	rr := httptest.NewRecorder()

	if len(data.Cookie) > 0 {
		for _, value := range data.Cookie {
			request.AddCookie(value)
			http.SetCookie(rr, value)
		}
	}

	handler := http.HandlerFunc(handlerFunc)
	handler.ServeHTTP(rr, request)


	if status := rr.Code; status != data.Status {
		t.Log(string(data.Body))
		t.Log("response body: ", rr.Body)
		t.Errorf("Method %s returned wrong status code: got %v want %v",
			data.Method, status, data.Status)
	}

	return rr.Result().Cookies()
}

func TestAuth(t *testing.T) {
	config := services.ReadConfig()
	db, err := services.NewDB(config)
	if err != nil {
		t.Fatal("Did not connected to database")
	}

	env := controllers.EnvironmentUser{db}

	for _, value := range gData {
		CheckStatus(value, env.ResponseRegisterHandler, t)
	}

	// Setup for next test
	// Case 1. Status 200
	body1 := []byte(`{
		"login":"dan_kokin",
		"password":"password"}`)
	gData[0].Body = body1

	// Case 2. Code 401
	body2 := []byte(`{
		"login":"dan_kokin2",
		"password":"password2"}`)
	gData[1].Body = body2
	gData[1].Status = http.StatusUnauthorized

	// Case 3. Code 401
	body3 := []byte(`{
		"login":"",
		"password":"password"}`)
	gData[2].Body = body3
	gData[2].Status = http.StatusUnauthorized

	gData = append(gData[:3])

	for _, value := range gData {
		CheckStatus(value, env.ResponseLoginHandler, t)
	}
}

func TestCRUD(t *testing.T) {
	config := services.ReadConfig()
	db, err := services.NewDB(config)
	if err != nil {
		t.Fatal(err)
	}

	env := controllers.EnvironmentUser{db}

	body4 := []byte(`{
		"login":"testuser4",
		"password":"password"}`)

	body5 := []byte(`{
		"login":"testuser5",
		"password":"password"}`)

	user4 := RequestData{Id: 4, Body: body4, Method: "POST", URL: "/user/4",
		Status: http.StatusOK, Cookie: nil}
	user5 := RequestData{Id: 5, Body: body5, Method: "POST", URL: "/user/5",
		Status: http.StatusOK, Cookie: nil}

	// Crutch because a cannot confirm user via mail
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), 8)
	err = db.QueryRow("INSERT INTO user_table (email, login, fullname, password, acc_verified) values ($1, $2, $3, $4, $5) RETURNING id",
		"testmail4", "testuser4", "test4 user4", hashedPassword, true).Scan(&user4.Id)

	if err != nil {
		t.Fatal(err)
	}

	hashedPassword, _ = bcrypt.GenerateFromPassword([]byte("password"), 8)
	err = db.QueryRow("INSERT INTO user_table (email, login, fullname, password, acc_verified) values ($1, $2, $3, $4, $5) RETURNING id",
		"testmail5", "testuser5", "test5 user5", hashedPassword, true).Scan(&user5.Id)
	if err != nil {
		t.Fatal(err)
	}

	// Login user4
	user4.Cookie = CheckStatus(user4, env.ResponseLoginHandler, t)

	user4.Body = nil
	user4.URL = "/user/4"
	user4.Method = "GET"

	// Get user4 with user4`s cookie
	CheckStatus(user4, env.GetUserHandler, t)

	// Try to get 5th user
	user4.URL = "/user/5"
	user4.Status = http.StatusUnauthorized
	user4.Id = 5
	CheckStatus(user4, env.GetUserHandler, t)

	// same with 5th user
	user5.Cookie = CheckStatus(user5, env.ResponseLoginHandler, t)

	user5.Body = nil
	user5.URL = "/user/5"
	user5.Method = "GET"

	// Get user5 with user5`s cookie
	CheckStatus(user4, env.GetUserHandler, t)

	// Try to get 4th user
	user5.URL = "/user/4"
	user5.Status = http.StatusUnauthorized
	user5.Id = 4
	CheckStatus(user4, env.GetUserHandler, t)


	// METHOD: PUT
	// Setup good cases
	body4 = []byte(`{
		"email":"testmail4",
		"login":"newtestuser4",
		"fullname":"Test User5",
		"password":"newpassword"}`)

	body5 = []byte(`{
		"email":"testmail5",
		"login":"newtestuser5",
		"fullname":"Test User5",
		"password":"newpassword"}`)

	user4 = RequestData{Id: 4, Body: body4, Method: "PUT", URL: "/user/4",
		Status: http.StatusOK, Cookie: user4.Cookie}
	user5 = RequestData{Id: 5, Body: body5, Method: "PUT", URL: "/user/5",
		Status: http.StatusOK, Cookie: user5.Cookie}

	CheckStatus(user4, env.UpdateUserHandler, t)
	CheckStatus(user5, env.UpdateUserHandler, t)

	// Setup bad cases
	// Code 401
	user4.URL = "/user/5"
	user4.Id = 5
	user4.Status = http.StatusUnauthorized

	user5.URL = "/user/4"
	user5.Id = 4
	user5.Status = http.StatusUnauthorized

	CheckStatus(user4, env.UpdateUserHandler, t)
	CheckStatus(user5, env.UpdateUserHandler, t)

	// Bad Request. Code 400
	user4.URL = "/user/4"
	user4.Id = 4
	user4.Status = http.StatusBadRequest
	user4.Body = []byte(`{
		"email":"testmail4",
		"login":"",
		"fullname":"Test User5",
		"password":"newpassword"}`)

	user5.URL = "/user/5"
	user5.Id = 5
	user5.Status = http.StatusBadRequest
	user5.Body = []byte(`{
		"email":"testmail5",
		"login":"testuser4",
		"fullname":"",
		"password":"newpassword"}`)

	CheckStatus(user4, env.UpdateUserHandler, t)
	CheckStatus(user5, env.UpdateUserHandler, t)

	// Method DELETE
	// Setup for 401 Code
	user4 = RequestData{Id: 5, Body: nil, Method: "DELETE", URL: "/user/5",
		Status: http.StatusUnauthorized, Cookie: user4.Cookie}
	user5 = RequestData{Id: 4, Body: nil, Method: "DELETE", URL: "/user/4",
		Status: http.StatusUnauthorized, Cookie: user5.Cookie}

	CheckStatus(user4, env.DeleteUserHandler, t)
	CheckStatus(user5, env.DeleteUserHandler, t)

	// Good case. Code 200
	user4.Id = 4
	user4.Status = http.StatusOK
	user4.URL = "/user/4"

	user5.Id = 5
	user5.Status = http.StatusOK
	user5.URL = "/user/5"

	CheckStatus(user4, env.DeleteUserHandler, t)
	CheckStatus(user5, env.DeleteUserHandler, t)
}
