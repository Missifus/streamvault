// File: internal/models/models.go
package models

import "github.com/golang-jwt/jwt/v4"

// User representa la estructura de un usuario en la base de datos.
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role"`
}

// Video representa la estructura de un video en la base de datos.
type Video struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	FilePath    string `json:"file_path"`
}

// Claims define la estructura del payload de nuestro JWT.
type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}
