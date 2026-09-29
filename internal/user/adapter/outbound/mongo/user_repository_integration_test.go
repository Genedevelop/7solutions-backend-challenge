//go:build integration

package mongo

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/Genedevelop/7solutions-backend-challenge/internal/user/domain/entity"
	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/errs"
	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/mongodb"
)

func TestUserRepoAgainstMongo(t *testing.T) {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongodb.Connect(ctx, uri)
	if err != nil {
		t.Fatal(err)
	}
	db := client.Database(fmt.Sprintf("it_%d", time.Now().UnixNano()))
	t.Cleanup(func() {
		_ = db.Drop(context.Background())
		_ = client.Disconnect(context.Background())
	})
	if err := mongodb.EnsureIndexes(ctx, db); err != nil {
		t.Fatal(err)
	}
	r := NewUserRepository(db)
	seed := func(name, email string) *entity.User {
		res, err := r.coll.InsertOne(ctx, userDocument{Name: name, Email: email, CreatedAt: time.Now().UTC()})
		if err != nil {
			t.Fatalf("seed %s: %v", email, err)
		}
		return &entity.User{ID: res.InsertedID.(bson.ObjectID).Hex(), Name: name, Email: email}
	}

	a := seed("A", "a@example.com")
	got, err := r.GetByID(ctx, a.ID)
	if err != nil || got.Email != "a@example.com" {
		t.Fatalf("GetByID = %+v, %v", got, err)
	}
	if _, err := r.GetByID(ctx, "not-an-object-id"); !errors.Is(err, errs.ErrUserNotFound) {
		t.Errorf("bad id: want ErrUserNotFound, got %v", err)
	}

	b := seed("B", "b@example.com")
	b.Email = "a@example.com"
	if err := r.Update(ctx, b); !errors.Is(err, errs.ErrEmailTaken) {
		t.Errorf("update to taken email: want ErrEmailTaken, got %v", err)
	}
	b.Name, b.Email = "Bee", "bee@example.com"
	if err := r.Update(ctx, b); err != nil {
		t.Fatalf("Update: %v", err)
	}

	users, err := r.List(ctx)
	if err != nil || len(users) != 2 || users[0].ID != a.ID || users[1].Name != "Bee" {
		t.Fatalf("List = %+v, %v", users, err)
	}
	if n, err := r.Count(ctx); err != nil || n != 2 {
		t.Errorf("Count = %d, %v", n, err)
	}

	if err := r.Delete(ctx, a.ID); err != nil {
		t.Fatal(err)
	}
	if err := r.Delete(ctx, a.ID); !errors.Is(err, errs.ErrUserNotFound) {
		t.Errorf("second delete: want ErrUserNotFound, got %v", err)
	}
}
