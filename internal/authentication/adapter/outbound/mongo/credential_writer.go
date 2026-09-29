package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/Genedevelop/7solutions-backend-challenge/internal/authentication/port/outbound"
	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/errs"
	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/mongodb"
)

var _ outbound.ICredentialWriter = (*CredentialWriter)(nil)

type CredentialWriter struct {
	coll *mongo.Collection
}

func NewCredentialWriter(db *mongo.Database) *CredentialWriter {
	return &CredentialWriter{coll: db.Collection(mongodb.UsersCollection)}
}

func (w *CredentialWriter) Create(ctx context.Context, account *outbound.Account) error {
	// Step 1: insert and let Mongo generate the ObjectID
	res, err := w.coll.InsertOne(ctx, accountDocument{
		Name:         account.Name,
		Email:        account.Email,
		PasswordHash: account.PasswordHash,
		CreatedAt:    account.CreatedAt,
	})

	// Step 2: the unique email index rejects a second account with the same email
	if mongo.IsDuplicateKeyError(err) {
		return errs.ErrEmailTaken
	}
	if err != nil {
		return errs.Internal(err, "insert account")
	}

	// Step 3: hand the generated id back to the caller
	account.ID = res.InsertedID.(bson.ObjectID).Hex()
	return nil
}
