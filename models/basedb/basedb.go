package basedb

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

type DB interface {
	Init() *sql.DB
}

type Database struct {
}

func (u *Database) Init() *sql.DB {
	username := os.Getenv("HANAZAWA_DB_USER")
	password := os.Getenv("HANAZAWA_DB_PASSWORD")
	host := os.Getenv("HANAZAWA_DB_HOST")
	port := os.Getenv("HANAZAWA_DB_PORT")
	db, err := sql.Open("mysql", username+":"+password+"@tcp("+host+":"+port+")/hanazawa?charset=utf8")
	if err != nil {
		panic(err.Error())
	}
	return db
}
