package db

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"
)

// DB is the Mongo database
var DB *mongo.Client

// Connect to the database
func Connect() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("DB_URL")))
	if err != nil {
		log.Fatal(err)
	}
	DB = client
}

func AddUser(user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := DB.Database("replier").Collection("users").UpdateOne(ctx, bson.D{{Key: "user_id", Value: user.UserID}}, bson.D{{Key: "$set", Value: bson.D{{Key: "user_id", Value: user.UserID}, {Key: "token", Value: user.Token}}}}, options.Update().SetUpsert(true))

	return err
}
