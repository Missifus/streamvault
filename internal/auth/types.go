package auth

import (
	"time"
)

// User representa a un usuario con seguridad adecuada
type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	IsVerified   bool      `json:"verified"`
	CreatedAt    time.Time `json:"created_at"`
}

// UserStore define la interfaz para almacenamiento de usuarios
type UserStore interface {
	CreateUser(user *User) error
	GetUserByEmail(email string) (*User, error)
	UpdateUser(user *User) error
}

// RegisterRequest representa la solicitud de registro
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest representa la solicitud de inicio de sesión
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse representa la respuesta de inicio de sesión
type LoginResponse struct {
	Token string `json:"token"`
}

