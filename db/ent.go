package db

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Accepts any type that implements method InsertOne()
type CollectionApi interface {
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
}

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
