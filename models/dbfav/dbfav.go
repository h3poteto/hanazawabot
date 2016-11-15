package dbfav

import (
	"time"

	"../basedb"

	"database/sql"
	"github.com/pkg/errors"
)

type Fav interface {
	Add(int, int) error
}

type DBFav struct {
	database *sql.DB
}

func (self *DBFav) Initialize() {
	self.database = basedb.SharedInstance().Connection
}

func NewDBFav() *DBFav {
	dbfav := &DBFav{}
	dbfav.Initialize()
	return dbfav
}

func (u *DBFav) Add(user_id int, youtube_movie_id int) error {
	_, err := u.database.Exec("insert into youtube_movie_favs (user_id, youtube_movie_id, created_at) values (?, ?, ?)", user_id, youtube_movie_id, time.Now())
	if err != nil {
		return errors.Wrap(err, "mysql insert error")
	}
	return nil
}
