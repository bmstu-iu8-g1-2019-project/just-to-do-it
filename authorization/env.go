package authorization

import (
	"database/sql"
	"net/http"
)

type Datastore interface {
	Login(string) (User, error)
	Register(r *http.Request) (error)
	Confirm(string) (error)
}

type Environment struct {
	db Datastore
}

type DB struct {
	*sql.DB
}
type User struct {
	Id int            `json:"id"`
	Email string      `json:"email"`
	Login string      `json:"login"`
	Fullname string   `json:"fullname"`
	Password string   `json:"password"`
	Acc_verified bool `json:"acc_verified"`
}
