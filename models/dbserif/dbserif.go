package dbserif

import (
	"../basedb"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type Serif interface {
	Add(string) bool
	SelectRandom() (tweet string, err error)
}

type DBSerif struct {
	dbobject basedb.DB
}

func (self *DBSerif) Initialize() {
	myDatabase := &basedb.Database{}
	var myDb basedb.DB = myDatabase
	self.dbobject = myDb
}

func NewDBSerif() *DBSerif {
	dbserif := &DBSerif{}
	dbserif.Initialize()
	return dbserif
}

func (u *DBSerif) Add(body string) bool {
	db := u.dbobject.Init()

	_, err := db.Exec("insert into serifs (body, created_at) values (?, ?)", body, time.Now())
	if err != nil {
		fmt.Printf("mysql connect error: %v \n", err)
		return false
	}

	defer db.Close()

	return true
}

func (u *DBSerif) SelectRandom() (tweet string, err error) {
	db := u.dbobject.Init()

	rows, err := db.Query("select * from serifs order by rand() limit 1;")
	if err != nil {
		fmt.Printf("mysql connect error: %v \n", err)
		return "", err
	}

	defer db.Close()

	id, body, created_at, updated_at := 0, "", "", ""
	for rows.Next() {
		err = rows.Scan(&id, &body, &created_at, &updated_at)
		if err != nil {
			return "", err
		}
	}

	return body, nil
}
