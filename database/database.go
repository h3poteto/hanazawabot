package database

import (
	"os"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type DB interface {
	Init() *sql.DB
}

type Database struct {
}

func (u *Database) Init() *sql.DB {
	username := os.Getenv("HANAZAWA_DB_USER")
	password := os.Getenv("HANAZAWA_DB_PASSWORD")
	db, err := sql.Open("mysql", username + ":" + password + "@/hanazawa?charset=utf8")
	if err != nil {
		panic(err.Error())
	}
	return db
}
