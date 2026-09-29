package mongo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/Genedevelop/7solutions-backend-challenge/internal/user/domain/entity"
	"github.com/Genedevelop/7solutions-backend-challenge/internal/user/port/outbound"
	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/errs"
	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/mongodb"
)

var _ outbound.IUserRepository = (*UserRepository)(nil)

type UserRepository struct {
	coll *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{coll: db.Collection(mongodb.UsersCollection)}
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
	// Step 1: an id that is not a valid ObjectID can never exist, so it is a plain not-found
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, errs.ErrUserNotFound
	}

	// Step 2: map "no document" to the domain error
	var doc userDocument
	err = r.coll.FindOne(ctx, bson.D{{Key: "_id", Value: oid}}).Decode(&doc)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, errs.ErrUserNotFound
	}
	if err != nil {
		return nil, errs.Internal(err, "find user")
	}
	return doc.toEntity(), nil
}

func (r *UserRepository) List(ctx context.Context) ([]*entity.User, error) {
	// Step 1: oldest first so the order is stable between calls
	cur, err := r.coll.Find(ctx, bson.D{}, options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}}))
	if err != nil {
		return nil, errs.Internal(err, "find users")
	}

	// Step 2: decode everything and map to the domain model
	var docs []userDocument
	if err := cur.All(ctx, &docs); err != nil {
		return nil, errs.Internal(err, "decode users")
	}
	users := make([]*entity.User, 0, len(docs))
	for _, d := range docs {
		users = append(users, d.toEntity())
	}
	return users, nil
}

func (r *UserRepository) Update(ctx context.Context, u *entity.User) error {
	oid, err := bson.ObjectIDFromHex(u.ID)
	if err != nil {
		return errs.ErrUserNotFound
	}

	// Step 1: only name and email are editable, password and created_at are never overwritten here
	res, err := r.coll.UpdateOne(ctx,
		bson.D{{Key: "_id", Value: oid}},
		bson.D{{Key: "$set", Value: bson.D{{Key: "name", Value: u.Name}, {Key: "email", Value: u.Email}}}},
	)
	if mongo.IsDuplicateKeyError(err) {
		return errs.ErrEmailTaken
	}
	if err != nil {
		return errs.Internal(err, "update user")
	}

	// Step 2: the user may have been deleted between read and write
	if res.MatchedCount == 0 {
		return errs.ErrUserNotFound
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return errs.ErrUserNotFound
	}
	res, err := r.coll.DeleteOne(ctx, bson.D{{Key: "_id", Value: oid}})
	if err != nil {
		return errs.Internal(err, "delete user")
	}
	if res.DeletedCount == 0 {
		return errs.ErrUserNotFound
	}
	return nil
}

func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	n, err := r.coll.CountDocuments(ctx, bson.D{})
	if err != nil {
		return 0, errs.Internal(err, "count users")
	}
	return n, nil
}
