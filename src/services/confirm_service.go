package services

import (
	"fmt"
	"net/smtp"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/models"
	)

const (
	url = "\nhttp://jtdi.ru/confirm?hash="
	msgConst = "\nFrom :%s\nTo: %s\nPlease confirm your email: %s"
)

//Генерирование хэша логниа
func addressGenerator(login string) (str string) {
	hashedLogin, _ := bcrypt.GenerateFromPassword([]byte(login), 4)
	return string(hashedLogin)
}

//запись в вспомогательную таблицу подтверждения почты
func (db *DB) recordMailConfirm (login string) (err error){
	secret := string(addressGenerator(login))
	deadlineTime := time.Now().Add(24 * time.Hour)
	_, err = db.Exec("INSERT INTO auth_confirmation (login, hash, deadline) values ($1, $2, $3)",
		login, string(secret), deadlineTime)
	if err != nil {
		return err
	}
	return nil
}

//обновление данных если переход был по старой ссылке
func (db *DB) confirmFieldUpdate(login string, hash string) (err error) {
	_, err = db.Exec("UPDATE auth_confirmation SET hash = $1, deadline = $2 where login = $3", hash, time.Now().Add(24 *time.Hour), login)
	if err != nil {
		return err
	}
	return nil
}

// Mail sending function
func (db *DB) sendMail(login string) (err error) {
	rowUser := db.QueryRow("SELECT * FROM user_table WHERE login = $1", login)
	var user models.User
	err = rowUser.Scan(&user.Id, &user.Email, &user.Login, &user.Fullname, &user.Password, &user.AccVerified)
	if err != nil {
		return err
	}
	row := db.QueryRow("SELECT * FROM auth_confirmation WHERE login = $1", login)
	var obj models.AuthConfirmation
	err = row.Scan(&obj.Login, &obj.Hash, &obj.Deadline)
	if err != nil {
		return err
	}
	from, _ := os.LookupEnv("from")
	pass, _ := os.LookupEnv("pass")
	fmt.Println(from, pass)
	msg := fmt.Sprintf(msgConst, from, user.Email, url + obj.Hash)

	err = smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth(
			"",
			from,
			pass,
			"smtp.gmail.com"),
		from, []string{user.Email}, []byte(msg))

	if err != nil {
		return err
	}
	return nil
}
