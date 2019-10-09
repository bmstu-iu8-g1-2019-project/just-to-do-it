package auth

import (
	"database/sql"
	"net/http"
	"time"
)

type Datastore interface {
	Login(string) (User, error)
	Register(User) (error)
	Confirm(string) (error)
	updateUser(int, User) (error)
	getUser(int) (User, error)
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

type Auth_confirmation struct {
	Login string `json:"login"`
	Hash string `json:"hash"`
	Deadline time.Time `json:"deadline"`
}
