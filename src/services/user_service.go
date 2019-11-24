package services

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/models"
)

type DatastoreUser interface {
	Login(string, string) (models.User, error)
	Register(models.User) (models.User, string, string)
	Confirm(string) error
	UpdateUser(int, models.User) (models.User, error)
	GetUser(int) (models.User, error)
	DeleteUser(int) error
}

// The function checks the username and password that comes
// from the request with the username and password that lies in the database
func (db *DB) Login(login string, password string) (models.User, error) {
	user := models.User{}
	row := db.QueryRow("SELECT * FROM user_table WHERE login = $1", login)
	err := row.Scan(&user.Id, &user.Email, &user.Login, &user.Fullname, &user.Password, &user.AccVerified)
	if err != nil {
		return user, err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return user, err
	}
	user.Password = ""
	return user, nil
}

func (db *DB) Register(user models.User) (models.User, string, string) {
	if len(user.Password) < 6 {
		return user, "Password must be more than 6 characters", "Bad Request"
	}
	// password hashing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 8)

	if user.Fullname == "" {
		user.Fullname = user.Login
	}
	query := "INSERT INTO user_table (email, login, fullname, password, acc_verified)" +
		"values ('%s', '%s', '%s', '%s', false)  RETURNING id"
	query = fmt.Sprintf(query, user.Email, user.Login, user.Fullname, hashedPassword)

	err = db.QueryRow(query).Scan(&user.Id)
	if err != nil {
		return user, "Query error", "Internal Server Error"
	}
	user.Password = ""
	// запись в вспомогательную таблицу которая хранит логин его хэш и срок протухания ссылка для подтверждения почты
	err = db.recordMailConfirm(user.Login)
	if err != nil {
		return user, "There was no record in the additional table", "Internal Server Error"
	}
	// отправка сообщения на почту юзера с таким логином
	err = db.sendMail(user.Login)
	if err != nil {
		return user, "Message was not sent", "Internal Server Error"
	}
	return user, "", ""
}

// get the user structure by id
func (db *DB) GetUser(id int) (models.User, error) {
	user := models.User{}
	row := db.QueryRow("SELECT * FROM user_table WHERE id = $1", id)
	err := row.Scan(&user.Id, &user.Email, &user.Login, &user.Fullname, &user.Password, &user.AccVerified)
	if err != nil {
		return user, err
	}
	user.Password = ""
	return user, nil
}

//update user data
func (db *DB) UpdateUser(id int, updateUser models.User) (models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateUser.Password), 8)
	_, err = db.Exec("UPDATE user_table SET email = $1, login = $2, fullname = $3,"+
		"password = $4 where id = $5", updateUser.Email,
		updateUser.Login, updateUser.Fullname, hashedPassword, id)
	if err != nil {
		return updateUser, err
	}
	updateUser.Id = id
	updateUser.Password = ""
	return updateUser, nil
}

func (db *DB) DeleteUser(id int) error {
	_, err := db.GetUser(id)
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE  FROM user_table WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) Confirm(hash string) (err error) {
	var conf models.AuthConfirmation
	row := db.QueryRow("SELECT * FROM auth_confirmation WHERE hash = $1", hash)
	err = row.Scan(&conf.Login, &conf.Hash, &conf.Deadline)
	if err != nil {
		fmt.Println(err)
		return err
	}
	if conf.Deadline.Before(time.Now()) {
		newHash := addressGenerator(conf.Login)
		err = db.confirmFieldUpdate(conf.Login, newHash)
		if err != nil {
			return err
		}
		err = db.sendMail(conf.Login)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	} else {
		_, err = db.Exec("UPDATE user_table SET acc_verified = true where login = $1", conf.Login)
		if err != nil {
			return err
		}
		_, err = db.Exec("DELETE  FROM auth_confirmation WHERE hash = $1", hash)
		if err != nil {
			return err
		}
		return nil
	}
}
