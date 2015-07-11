package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type DB interface {
	Init() *sql.DB
}

type Database struct {
}

func (u *Database) Init() *sql.DB {
	db, err := sql.Open("mysql", "root:@/hanazawa?charset=utf8")
	if err != nil {
		panic(err.Error())
	}
	return db
}
