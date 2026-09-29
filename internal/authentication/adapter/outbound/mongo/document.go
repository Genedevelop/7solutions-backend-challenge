package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type credentialDocument struct {
	ID           bson.ObjectID `bson:"_id"`
	PasswordHash string        `bson:"password_hash"`
}

type accountDocument struct {
	Name         string    `bson:"name"`
	Email        string    `bson:"email"`
	PasswordHash string    `bson:"password_hash"`
	CreatedAt    time.Time `bson:"created_at"`
}
