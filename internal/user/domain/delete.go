package domain

import (
	"context"
	"strings"

	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/errs"
)

func (s *UserService) Delete(ctx context.Context, actorID, id string) error {
	// Step 1: users can only delete their own account, ids are hex so letter case does not matter
	if !strings.EqualFold(actorID, id) {
		return errs.ErrForbidden
	}
	return s.repo.Delete(ctx, actorID)
}
