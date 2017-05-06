package dbserif

import (
	"database/sql"
	"time"

	"github.com/h3poteto/hanazawabot/models/basedb"

	"github.com/pkg/errors"
)

type Serif interface {
	Add(string) error
	SelectRandom() (tweet string, err error)
}

type DBSerif struct {
	database *sql.DB
}

func (self *DBSerif) Initialize() {
	self.database = basedb.SharedInstance().Connection
}

func NewDBSerif() *DBSerif {
	dbserif := &DBSerif{}
	dbserif.Initialize()
	return dbserif
}

func (u *DBSerif) Add(body string) error {
	_, err := u.database.Exec("insert into serifs (body, created_at) values (?, ?)", body, time.Now())
	if err != nil {
		return errors.Wrap(err, "mysql insert error")
	}

	return nil
}

func (u *DBSerif) SelectRandom() (tweet string, err error) {
	rows, err := u.database.Query("select * from serifs order by rand() limit 1;")
	if err != nil {
		return "", errors.Wrap(err, "mysql select error")
	}

	id, body, created_at, updated_at := 0, "", "", ""
	for rows.Next() {
		err = rows.Scan(&id, &body, &created_at, &updated_at)
		if err != nil {
			return "", errors.Wrap(err, "mysql scan error")
		}
	}

	return body, nil
}
