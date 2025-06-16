package auth

import (
	"errors"
	"sync"
)

// MemoryUserStore implementa UserStore en memoria
type MemoryUserStore struct {
	mu     sync.RWMutex
	users  map[string]*User
	lastID int
}

func NewMemoryUserStore() *MemoryUserStore {
	return &MemoryUserStore{
		users: make(map[string]*User),
	}
}

func (s *MemoryUserStore) CreateUser(user *User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Verificar si el usuario ya existe
	if _, exists := s.users[user.Email]; exists {
		return errors.New("el usuario ya existe")
	}

	s.lastID++
	user.ID = s.lastID
	s.users[user.Email] = user
	return nil
}

func (s *MemoryUserStore) GetUserByEmail(email string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[email]
	if !exists {
		return nil, errors.New("usuario no encontrado")
	}
	return user, nil
}

func (s *MemoryUserStore) UpdateUser(user *User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[user.Email]; !exists {
		return errors.New("usuario no encontrado")
	}
	s.users[user.Email] = user
	return nil
}