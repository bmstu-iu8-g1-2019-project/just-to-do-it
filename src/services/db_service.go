package services

import (
	"database/sql"
	"fmt"
	"io"
	"os"
)

type DB struct {
	*sql.DB
}

func NewDB(dbSourceName string) (*DB, error) {
	db, err := sql.Open("postgres",dbSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func OpenConfigFile(filename string) (config string){
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	data := make([]byte, 64)
	for {
		n, err := file.Read(data)
		if err == io.EOF {
			break
		}
		config += string(data[:n])
	}
	return config
}
