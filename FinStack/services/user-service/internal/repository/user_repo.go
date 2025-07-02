package repository

import (
	"errors"
	"sync"
)

type User struct {
	ID       string
	Name     string
	Email    string
	Password string
}

type UserRepository interface {
	Create(user User) error
	FindByEmail(email string) (*User, error)
}

// InMemoryUserRepo is a thread-safe in-memory store for testing.
type InMemoryUserRepo struct {
	users map[string]User
	mu    sync.RWMutex
}

func NewInMemoryUserRepo() *InMemoryUserRepo {
	return &InMemoryUserRepo{
		users: make(map[string]User),
	}
}

func (r *InMemoryUserRepo) Create(user User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.users[user.Email]; exists {
		return errors.New("user already exists")
	}
	r.users[user.Email] = user
	return nil
}

func (r *InMemoryUserRepo) FindByEmail(email string) (*User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	user, exists := r.users[email]
	if !exists {
		return nil, errors.New("user not found")
	}
	return &user, nil
}
