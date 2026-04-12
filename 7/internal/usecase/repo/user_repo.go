package repo

import (
	"errors"
	"sync"

	"practice-7/internal/entity"

	"github.com/google/uuid"
)

type UserRepo struct {
	mu    sync.RWMutex
	users map[uuid.UUID]*entity.User
}

func NewUserRepo() *UserRepo {
	return &UserRepo{users: make(map[uuid.UUID]*entity.User)}
}

func (r *UserRepo) SeedAdmin(admin *entity.User) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[admin.ID] = admin
}

func (r *UserRepo) RegisterUser(user *entity.User) (*entity.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, u := range r.users {
		if u.Username == user.Username {
			return nil, errors.New("username already taken")
		}
		if u.Email == user.Email {
			return nil, errors.New("email already taken")
		}
	}

	user.ID = uuid.New()
	r.users[user.ID] = user
	return user, nil
}

func (r *UserRepo) GetByUsername(username string) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, u := range r.users {
		if u.Username == username {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *UserRepo) GetByID(id uuid.UUID) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return u, nil
}

func (r *UserRepo) UpdateRole(id uuid.UUID, role string) (*entity.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	u, ok := r.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	u.Role = role
	return u, nil
}
