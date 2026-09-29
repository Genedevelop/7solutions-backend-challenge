//go:build integration

package mongo

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Genedevelop/7solutions-backend-challenge/internal/authentication/port/outbound"
	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/errs"
	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/mongodb"
)

func TestCredentialsAgainstMongo(t *testing.T) {
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
	w, r := NewCredentialWriter(db), NewCredentialReader(db)

	a := &outbound.Account{Name: "A", Email: "a@example.com", PasswordHash: "hash-a", CreatedAt: time.Now().UTC()}
	if err := w.Create(ctx, a); err != nil || a.ID == "" {
		t.Fatalf("Create: id=%q err=%v", a.ID, err)
	}
	if err := w.Create(ctx, &outbound.Account{Name: "dup", Email: "a@example.com"}); !errors.Is(err, errs.ErrEmailTaken) {
		t.Fatalf("duplicate: want ErrEmailTaken, got %v", err)
	}

	cred, err := r.GetByEmail(ctx, "a@example.com")
	if err != nil || cred.UserID != a.ID || cred.PasswordHash != "hash-a" {
		t.Fatalf("GetByEmail = %+v, %v", cred, err)
	}
	if _, err := r.GetByEmail(ctx, "nobody@example.com"); !errors.Is(err, errs.ErrCredentialNotFound) {
		t.Errorf("unknown email: want ErrCredentialNotFound, got %v", err)
	}
}
