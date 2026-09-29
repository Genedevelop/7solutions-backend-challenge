package mock

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/Genedevelop/7solutions-backend-challenge/internal/user/domain/entity"
	"github.com/Genedevelop/7solutions-backend-challenge/internal/user/port/outbound"
	"github.com/Genedevelop/7solutions-backend-challenge/shared/utils/errs"
)

var _ outbound.IUserRepository = (*UserRepository)(nil)

type UserRepository struct {
	mu     sync.Mutex
	users  map[string]entity.User
	nextID int
	Err    error
}

func NewUserRepository() *UserRepository {
	return &UserRepository{users: map[string]entity.User{}}
}

func (r *UserRepository) Create(_ context.Context, u *entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.Err != nil {
		return r.Err
	}
	if r.emailTaken(u.Email, "") {
		return errs.ErrEmailTaken
	}
	r.nextID++
	u.ID = fmt.Sprintf("user-%d", r.nextID)
	r.users[u.ID] = *u
	return nil
}

func (r *UserRepository) GetByID(_ context.Context, id string) (*entity.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.Err != nil {
		return nil, r.Err
	}
	u, ok := r.users[id]
	if !ok {
		return nil, errs.ErrUserNotFound
	}
	return &u, nil
}

func (r *UserRepository) List(_ context.Context) ([]*entity.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.Err != nil {
		return nil, r.Err
	}
	out := make([]*entity.User, 0, len(r.users))
	for _, u := range r.users {
		out = append(out, &u)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

func (r *UserRepository) Update(_ context.Context, u *entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.Err != nil {
		return r.Err
	}
	if _, ok := r.users[u.ID]; !ok {
		return errs.ErrUserNotFound
	}
	if r.emailTaken(u.Email, u.ID) {
		return errs.ErrEmailTaken
	}
	r.users[u.ID] = *u
	return nil
}

func (r *UserRepository) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.Err != nil {
		return r.Err
	}
	if _, ok := r.users[id]; !ok {
		return errs.ErrUserNotFound
	}
	delete(r.users, id)
	return nil
}

func (r *UserRepository) Count(_ context.Context) (int64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.Err != nil {
		return 0, r.Err
	}
	return int64(len(r.users)), nil
}

func (r *UserRepository) emailTaken(email, exceptID string) bool {
	for id, u := range r.users {
		if u.Email == email && id != exceptID {
			return true
		}
	}
	return false
}
