package auth

import (
	"time"
)

// User representa a un usuario con seguridad adecuada
type User struct {
    ID           int       `json:"id"`           // ID único (como número de cédula)
    Email        string    `json:"email"`         // Correo del usuario
    PasswordHash string    `json:"-"`             // Contraseña encriptada (NO se muestra)
    Role         string    `json:"role"`          // Rol: "user", "admin", etc.
    IsVerified   bool      `json:"verified"`     // ¿Verificó su email?
    CreatedAt    time.Time `json:"created_at"`   // Fecha de registro
}

// UserStore define la interfaz para almacenamiento de usuarios
type UserStore interface {
	CreateUser(user *User) error              // Crear nuevo usuario
    GetUserByEmail(email string) (*User, error) // Buscar por email
    UpdateUser(user *User) error              // Actualizar datos
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

