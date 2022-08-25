package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Ride struct {
	Id       int
	RideId   string
	Clientid string
	DriverId string
	Status   string
}

var (
	RideRespository *SQLiteRepository
)

func init() {

	db, err := sql.Open("sqlite3", "ride.db")
	if err != nil {
		log.Fatal(err)
	}

	RideRespository = NewSqliteRepository(db)
	if err := RideRespository.Migrate(); err != nil {
		log.Fatal(err)
	}

}
