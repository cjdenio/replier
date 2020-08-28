package db

import (
	"context"
	"crypto/rand"
	"encoding/hex"
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

// AddInstallation adds an installation to the database
func AddInstallation(installation Installation) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := DB.Database("replier").Collection("installations").UpdateOne(ctx, bson.M{"team_id": installation.TeamID}, bson.M{"$set": installation}, options.Update().SetUpsert(true))

	if err != nil {
		return err
	}

	return nil
}

// GetInstallation gets an installation
func GetInstallation(teamID string) (*Installation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var installation *Installation

	err := DB.Database("replier").Collection("installations").FindOne(ctx, bson.M{"team_id": teamID}).Decode(&installation)

	if err != nil {
		return &Installation{}, err
	}

	return installation, nil
}

// AddUser adds a user
func AddUser(user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := DB.Database("replier").Collection("users").UpdateOne(ctx, bson.D{{Key: "user_id", Value: user.UserID}}, bson.D{{Key: "$set", Value: bson.D{{Key: "user_id", Value: user.UserID}, {Key: "token", Value: user.Token}, {Key: "scopes", Value: user.Scopes}, {Key: "team_id", Value: user.TeamID}}}}, options.Update().SetUpsert(true))

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

// GetUserByToken gets a User, provided an API token
func GetUserByToken(token string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var result *User

	err := DB.Database("replier").Collection("users").FindOne(ctx, bson.D{{Key: "api_token", Value: token}}).Decode(&result)

	if err != nil {
		return &User{}, err
	}
	return result, nil
}

// GetUserAPIToken gets a user's API token, creating it if necessary.
func GetUserAPIToken(userID string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var result *User

	err := DB.Database("replier").Collection("users").FindOne(ctx, bson.D{{Key: "user_id", Value: userID}}).Decode(&result)

	if err != nil {
		return "", err
	}

	if result.APIToken == "" {
		token := generateAPIToken()
		setUserAPIToken(userID, token)
		return token, nil
	}
	return result.APIToken, nil
}

func setUserAPIToken(userID, token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := DB.Database("replier").Collection("users").UpdateOne(ctx, bson.M{"user_id": userID}, bson.M{"$set": bson.M{"api_token": token}})

	return err
}

func generateAPIToken() string {
	token := make([]byte, 16)
	rand.Read(token)
	return hex.EncodeToString(token)
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

	return err
}

// SetUserDates sets a user's start/end dates
func SetUserDates(start, end time.Time, userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := DB.Database("replier").Collection("users").UpdateOne(ctx, bson.M{"user_id": userID}, bson.M{"$set": bson.M{"reply.start": start, "reply.end": end}})

	return err
}

// ToggleReplyActive toggle's the activity of a user's autoreply
func ToggleReplyActive(userID string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	user, err := GetUser(userID)

	if err != nil {
		fmt.Println(err)
	}

	_, err = DB.Database("replier").Collection("users").UpdateOne(ctx, bson.D{{Key: "user_id", Value: userID}}, bson.M{"$set": bson.M{"reply.active": !user.Reply.Active}})
	if err != nil {
		fmt.Println(err)
	}
}

// GetConversationLastPostedOn gets the Time that the conversation was last autoreplied to.
func GetConversationLastPostedOn(conversationID, userID string) time.Time {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var conversation *Conversation

	err := DB.Database("replier").Collection("conversations").FindOne(ctx, bson.M{"user_id": userID, "conversation_id": conversationID}).Decode(&conversation)

	if err != nil {
		return time.Time{}
	}

	result := time.Unix(conversation.LastPostedOn, 0)

	return result
}

// SetConversationLastPostedOn sets the Time above ^
func SetConversationLastPostedOn(conversationID, userID string, lastPostedOn time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := DB.Database("replier").Collection("conversations").UpdateOne(ctx, bson.D{
		{Key: "conversation_id", Value: conversationID},
		{Key: "user_id", Value: userID},
	}, bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "last_posted_on", Value: lastPostedOn.Unix()},
		}},
	}, options.Update().SetUpsert(true))

	return err
}
