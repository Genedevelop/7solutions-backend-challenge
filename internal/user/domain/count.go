package domain

import "context"

func (s *UserService) Count(ctx context.Context) (int64, error) {
	return s.repo.Count(ctx)
}
