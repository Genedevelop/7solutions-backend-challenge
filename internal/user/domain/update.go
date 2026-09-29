package domain

import (
	"context"
	"strings"

	"github.com/Genedevelop/7solutions-backend-challenge/internal/user/domain/entity"
	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/errs"
)

func (s *UserService) Update(ctx context.Context, actorID, id string, name, email *string) (*entity.User, error) {
	// Step 1: users can only edit their own account, ids are hex so letter case does not matter
	if !strings.EqualFold(actorID, id) {
		return nil, errs.ErrForbidden
	}

	// Step 2: load the current user and apply the validated changes
	user, err := s.repo.GetByID(ctx, actorID)
	if err != nil {
		return nil, err
	}
	if err := user.ApplyUpdate(name, email); err != nil {
		return nil, err
	}

	// Step 3: save, a taken email is still caught by the unique index
	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}
