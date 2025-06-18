package auth

import (
	"errors" // Paquete para manejar errores
	"sync" // Paquete para sincronización concurrente
)

// MemoryUserStore implementa UserStore en memoria
type MemoryUserStore struct {
	mu     sync.RWMutex  // Evita accesos simultáneos
    users  map[string]*User  // Almacén: Email -> Usuario
    lastID int            // Autoincremental para IDs
}
// constructor 
func NewMemoryUserStore() *MemoryUserStore {
	return &MemoryUserStore{
		users: make(map[string]*User),
	}
}
//registrar un nuevo usuario
func (s *MemoryUserStore) CreateUser(user *User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Verificar si el usuario ya existe
	if _, exists := s.users[user.Email]; exists {
		return errors.New("el usuario ya existe")
	}
	// Asignar nuevo id
	s.lastID++
	user.ID = s.lastID
	s.users[user.Email] = user
	return nil
}
//Buscar Usuario por Email
func (s *MemoryUserStore) GetUserByEmail(email string) (*User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[email]
	if !exists {
		return nil, errors.New("usuario no encontrado")
	}
	return user, nil
}
//Actualizar Usuario
func (s *MemoryUserStore) UpdateUser(user *User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.users[user.Email]; !exists {
		return errors.New("usuario no encontrado")
	}
	s.users[user.Email] = user
	return nil
}