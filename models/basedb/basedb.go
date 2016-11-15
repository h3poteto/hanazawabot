package basedb

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	Connection *sql.DB
}

var sharedInstance = New()

// New prepare database instance
func New() *Database {
	username := os.Getenv("HANAZAWA_DB_USER")
	password := os.Getenv("HANAZAWA_DB_PASSWORD")
	host := os.Getenv("HANAZAWA_DB_HOST")
	port := os.Getenv("HANAZAWA_DB_PORT")
	pool := 5
	db, err := sql.Open("mysql", username+":"+password+"@tcp("+host+":"+port+")/hanazawa?charset=utf8mb4")
	if err != nil {
		panic(err.Error())
	}
	db.SetMaxIdleConns(pool)
	db.SetMaxOpenConns(pool)
	return &Database{
		Connection: db,
	}
}

// SharedInstance for database instance
func SharedInstance() *Database {
	return sharedInstance
}

// Close connection to database
func (d *Database) Close() error {
	return d.Connection.Close()
}
