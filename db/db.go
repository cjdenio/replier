package db

import (
	"context"
	"fmt"
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

// AddUser adds a user
func AddUser(user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := DB.Database("replier").Collection("users").UpdateOne(ctx, bson.D{{Key: "user_id", Value: user.UserID}}, bson.D{{Key: "$set", Value: bson.D{{Key: "user_id", Value: user.UserID}, {Key: "token", Value: user.Token}}}}, options.Update().SetUpsert(true))

	return err
}

// GetUser gets a user based off of a user_id
func GetUser(userID string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var result *User

	err := DB.Database("replier").Collection("users").FindOne(ctx, bson.D{{Key: "user_id", Value: userID}}).Decode(&result)

	if err != nil {
		return &User{}, err
	}
	return result, nil
}

// SetUserMessage sets a users message
func SetUserMessage(userID string, message string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := DB.Database("replier").Collection("users").UpdateOne(ctx, bson.D{{Key: "user_id", Value: userID}}, bson.D{{Key: "$set", Value: bson.D{{Key: "reply.message", Value: message}}}})

	if err != nil {
		return err
	}

	return nil
}

// SetUserWhitelist sets a user's whitelist
func SetUserWhitelist(userID string, whitelist []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := DB.Database("replier").Collection("users").UpdateOne(ctx, bson.D{{Key: "user_id", Value: userID}}, bson.D{{Key: "$set", Value: bson.D{{Key: "reply.whitelist", Value: whitelist}}}})

	if err != nil {
		return err
	}

	return nil
}

// ToggleReplyActive toggle's the activity of a user's autoreply
func ToggleReplyActive(userID string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := DB.Database("replier").Collection("users").UpdateOne(ctx, bson.D{{Key: "user_id", Value: userID}}, bson.M{"$bit": bson.M{"reply.active": bson.M{"xor": 1}}})
	if err != nil {
		fmt.Println(err)
	}
}
