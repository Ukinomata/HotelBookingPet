package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

var (
	user     = "postgres"
	password = "ukino"
	host     = "localhost"
	dbname   = "HotelApp"
	sslmode  = "disable"
)

var DbInfo = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s", user, password, host, dbname, sslmode) //DataSourceName

func ConnectToDB() *sql.DB {
	db, err := sql.Open("postgres", DbInfo)
	if err != nil {
		panic(err)
	}

	return db
}
