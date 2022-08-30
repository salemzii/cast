package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/salemzii/cast.git/app"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	tykCollection *mongo.Collection
)

// connects to mongodb server and defines value for Collection
func PrepareMongo() {
	mongo_uri := os.Getenv("MONGO_URI")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// connect to mongodb client
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongo_uri))
	if err != nil {
		log.Println(err)
	}

	if err := MigrateMongodb(client); err != nil {
		log.Println(err)
	}
}

// Makes migration for mongodb collection
func MigrateMongodb(db *mongo.Client) error {
	tykCollection = db.Database("testing").Collection("cast")
	return nil
}

// Adds record to a mongodb collection
// implements CollectionApi type
func AddDataRecordMongodb(collection CollectionApi, data *app.Message) (CreatedData *app.Message, err error) {

	bson_data := bson.D{{Key: "app_id", Value: data.AppId}, {Key: "device_id", Value: data.ClientId}}
	result, err := collection.InsertOne(context.TODO(), bson_data)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	if result.InsertedID == 0 {
		log.Println(err)
		return nil, err
	}
	return data, nil
}

func VerifyAppId(app_id string) bool {
	coll := tykCollection

	var result bson.M
	err := coll.FindOne(context.TODO(), bson.D{{"app_id", app_id}}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return false
		}
		log.Println(err)
	}

	return true
}
