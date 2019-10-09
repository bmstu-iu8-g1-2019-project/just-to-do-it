package main

import "golang.org/x/crypto/bcrypt"

func (db *DB) getUser (id int) (user User, err error) {
	row := db.QueryRow("SELECT * FROM usertab WHERE id = $1", id)
	err = row.Scan(&user.Id, &user.Email, &user.Login, &user.Fullname, &user.Password, &user.Acc_verified)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (db *DB) updateUser (id int, updateUser User) (err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateUser.Password), 8)
	_, err = db.Exec("UPDATE usertab SET email = $1, login = $2, fullname = $3, password = $4, acc_verified = $5 where id = $6", updateUser.Email,
		updateUser.Login, updateUser.Fullname, hashedPassword, updateUser.Acc_verified, id)
	if err != nil {
		return err
	}
	return nil
}

