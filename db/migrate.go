package db

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// migrate defines a DB migration. It is run once upon every connection to the database.
func migrate() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	result, err := DB.Database("replier").Collection("users").UpdateMany(ctx, bson.M{
		"reply.mode": bson.M{"$exists": false}, "$and": bson.A{
			bson.M{
				"$or": bson.A{
					bson.M{"reply.start": nil},
					bson.M{"reply.start": time.Time{}},
				},
			},
			bson.M{
				"$or": bson.A{
					bson.M{"reply.end": nil},
					bson.M{"reply.end": time.Time{}},
				},
			},
		},
	}, bson.M{
		"$set": bson.M{"reply.mode": ReplyModeManual},
	})
	if err != nil {
		return err
	}

	log.Printf("Successfully migrated %d DB records to ReplyModeManual", result.ModifiedCount)

	result, err = DB.Database("replier").Collection("users").UpdateMany(ctx, bson.M{
		"reply.mode": bson.M{"$exists": false}, "$or": bson.A{
			bson.M{"reply.start": bson.M{"$ne": nil}},
			bson.M{"reply.end": bson.M{"$ne": nil}},
			bson.M{"reply.start": bson.M{"$ne": time.Time{}}},
			bson.M{"reply.end": bson.M{"$ne": time.Time{}}},
		},
	}, bson.M{
		"$set": bson.M{"reply.mode": ReplyModeDate},
	})
	if err != nil {
		return err
	}

	log.Printf("Successfully migrated %d DB records to ReplyModeDate", result.ModifiedCount)

	return nil
}
