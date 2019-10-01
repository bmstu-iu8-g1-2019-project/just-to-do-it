package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockDB struct{}

func (mdb *mockDB) GetData() ([]*Object, error) {
	obj := make([]*Object, 0)
	obj = append(obj, &Object{"string"})
	return obj, nil
}

func (mdb *mockDB) GetOneData(arg string) (Object, error) {
	obj := Object{"string"}
	return obj, nil
}

func TestRespAllData(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/serv", nil)

	env := Environment{db: &mockDB{}}
	http.HandlerFunc(env.respAllData).ServeHTTP(rec, req)

	expected := "string"
	if expected != rec.Body.String() {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, rec.Body.String())
	}
}

func TestRespOneData (t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/serv/string", nil)

	env := Environment{db: &mockDB{}}
	http.HandlerFunc(env.respOneData).ServeHTTP(rec, req)

	expected := "string"
	if expected != rec.Body.String() {
		t.Errorf("\n...expected = %v\n...obtained = %v", expected, rec.Body.String())
	}
}

