package authorization

import (
	"crypto/sha256"
	"encoding/hex"
	"net/smtp"
	"os"
	"time"
)

type Auth_confirmation struct {
	Login string `json:"login"`
	Hash string `json:"hash"`
	Deadline time.Time `json:"deadline"`
}

// Функция генерирует хэш, пока без соли:(
func addressGenerator(login string) (str string) {
	h := sha256.New()
	h.Write([]byte(login))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

// Функция записывает данные в таблицу которая хранит информацию о логине, хэше и деадлайне
func (db *DB) recordMailConfirm (login string, hash string) (err error){
	deadlineTime := time.Now().Add(24 * time.Hour)
	_, err = db.Exec("INSERT INTO auth_confirmation (login, hash, deadline) values ($1, $2, $3)",
		login, hash, deadlineTime)
	if err != nil {
		return err
	}
	return nil
}

// Функция отправки сообщения
func sendMail(login string, email string) (hash string, err error) {
	hash = addressGenerator(login)
	url := "\nlocalhost:3000/confirm?hash=" + hash
	from := os.Getenv("email")
	pass := os.Getenv("pass")

	msg := "\nFrom :" + from + "\n" +
		"To: " + email + "\n" +
		"Please confirm your email: " +
		url

	err = smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth(
			"",
			from,
			pass,
			"smtp.gmail.com"),
		from, []string{email}, []byte(msg))

	if err != nil {
		return "", err
	}
	return hash, nil
}