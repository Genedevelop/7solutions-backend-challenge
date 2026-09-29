package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/Genedevelop/7solutions-backend-challenge/internal/user/domain/entity"
)

type userDocument struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	Name      string        `bson:"name"`
	Email     string        `bson:"email"`
	CreatedAt time.Time     `bson:"created_at"`
}

func (d userDocument) toEntity() *entity.User {
	return &entity.User{
		ID:        d.ID.Hex(),
		Name:      d.Name,
		Email:     d.Email,
		CreatedAt: d.CreatedAt.UTC(),
	}
}
