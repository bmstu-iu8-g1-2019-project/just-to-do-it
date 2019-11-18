package models

type Message struct {
	Err error  `json:"error"`
	Msg string `json:"message"`
}
