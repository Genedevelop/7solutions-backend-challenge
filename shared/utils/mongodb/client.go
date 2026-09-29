package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"

	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/errs"
)

const UsersCollection = "user_storage"

func Connect(ctx context.Context, uri string) (*mongo.Client, error) {
	// Step 1: create the client, the driver connects lazily
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, errs.Internal(err, "mongo connect")
	}

	// Step 2: ping so a wrong URI fails at startup instead of on the first request
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		_ = client.Disconnect(context.Background())
		return nil, errs.Internal(err, "mongo ping")
	}
	return client, nil
}

func EnsureIndexes(ctx context.Context, db *mongo.Database) error {
	// Step 1: unique email index on the collection shared by user and authentication, it is what really blocks duplicate accounts under concurrent sign-ups
	_, err := db.Collection(UsersCollection).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("uniq_email"),
	})
	if err != nil {
		return errs.Internal(err, "create email index")
	}
	return nil
}
