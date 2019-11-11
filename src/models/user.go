package models

import "time"

type User struct {
	Id int            `json:"id"`
	Email string      `json:"email"`
	Login string      `json:"login"`
	Fullname string   `json:"fullname"`
	Password string   `json:"password"`
	AccVerified bool  `json:"acc_verified"`
}

// helper structure for mail confirmation
type AuthConfirmation struct {
	Login string       `json:"login"`
	Hash string        `json:"hash"`
	Deadline time.Time `json:"deadline"`
}
