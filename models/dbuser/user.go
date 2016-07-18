package dbuser

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"

	"../basedb"
)

type User interface {
	Add(int64, string) error
	Select(int64, string) (*DBUser, error)
	SelectOrAdd(int64, string) (*DBUser, error)
}

type DBUser struct {
	Id         int
	ScreenName string
	TwitterID  int64
	dbobject   basedb.DB
}

func (self *DBUser) Initialize() {
	myDatabase := &basedb.Database{}
	var myDb basedb.DB = myDatabase
	self.dbobject = myDb
}

func NewDBUser() *DBUser {
	dbuser := &DBUser{}
	dbuser.Initialize()
	return dbuser
}

func (u *DBUser) Add(twitter_id int64, screen_name string) error {
	db := u.dbobject.Init()
	defer db.Close()

	_, err := db.Exec("insert into users (twitter_id, screen_name, created_at) values (?, ?, ?)", twitter_id, screen_name, time.Now())
	if err != nil {
		return errors.Wrap(err, "mysql insert error")
	}

	return nil
}

func (u *DBUser) Select(twitter_id int64, screen_name string) (*DBUser, error) {
	db := u.dbobject.Init()
	defer db.Close()

	var rows *sql.Rows
	var err error
	if twitter_id > 0 {
		rows, err = db.Query("select * from users where twitter_id = ?;", twitter_id)
		if err != nil {
			return nil, errors.Wrap(err, "mysql select error")
		}
	} else {
		rows, err = db.Query("select * from users where screen_name = ?;", screen_name)
		if err != nil {
			return nil, errors.Wrap(err, "mysql select error")
		}
	}

	var twitter int64
	id, twitter, screen_name, created_at, updated_at := 0, 0, "", "", ""
	for rows.Next() {
		err := rows.Scan(&id, &twitter, &screen_name, &created_at, &updated_at)
		if err != nil {
			return nil, errors.Wrap(err, "mysql scan error")
		}
	}

	user := DBUser{Id: id, TwitterID: twitter, ScreenName: screen_name}
	return &user, nil
}

func (u *DBUser) SelectOrAdd(twitter_id int64, screen_name string) (*DBUser, error) {
	user, err := u.Select(twitter_id, screen_name)
	if err != nil {
		err := u.Add(twitter_id, screen_name)
		if err != nil {
			return nil, err
		}
		user, err = u.Select(twitter_id, screen_name)
		if err != nil {
			return nil, err
		}
	}
	return user, nil
}
