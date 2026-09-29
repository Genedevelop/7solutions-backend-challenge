package mongo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/Genedevelop/7solutions-backend-challenge/internal/authentication/port/outbound"
	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/errs"
	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/mongodb"
)

var _ outbound.ICredentialReader = (*CredentialReader)(nil)

type CredentialReader struct {
	coll *mongo.Collection
}

func NewCredentialReader(db *mongo.Database) *CredentialReader {
	return &CredentialReader{coll: db.Collection(mongodb.UsersCollection)}
}

func (r *CredentialReader) GetByEmail(ctx context.Context, email string) (*outbound.Credential, error) {
	// Step 1: project only the id and password hash, the unique email index makes this a single index lookup
	var doc credentialDocument
	opts := options.FindOne().SetProjection(bson.D{{Key: "password_hash", Value: 1}})
	err := r.coll.FindOne(ctx, bson.D{{Key: "email", Value: email}}, opts).Decode(&doc)

	// Step 2: translate "no document" into the port error
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errs.ErrCredentialNotFound
	}
	if err != nil {
		return nil, errs.Internal(err, "find credential")
	}
	return &outbound.Credential{UserID: doc.ID.Hex(), PasswordHash: doc.PasswordHash}, nil
}
