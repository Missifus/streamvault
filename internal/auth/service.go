// logica de Autenticacion
package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt" //Implementación de JSON Web Tokens
	"golang.org/x/crypto/bcrypt" //Algoritmo para hashear contraseñas 
)

// Estructura del Servicio
type AuthService struct {
	store UserStore
}

// NewAuthService crea un nuevo servicio de autenticación
func NewAuthService(store UserStore) *AuthService {
	return &AuthService{store: store}
}

// Register maneja el registro de nuevos usuarios
func (s *AuthService) Register(email, password string) error {
	// Validar email
	if !isValidEmail(email) {
		return errors.New("formato de email inválido")
	}
	
	// Validar fortaleza de contraseña
	if !isStrongPassword(password) {
		return errors.New("la contraseña debe tener al menos 8 caracteres e incluir letras mayúsculas, minúsculas, números y/o símbolos")
	}

	// Verificar si el usuario ya existe
	if _, err := s.store.GetUserByEmail(email); err == nil {
		return errors.New("el usuario ya existe")
	}

	// Hashear contraseña con bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error al hashear contraseña: %w", err)
	}

	// Crear nuevo usuario
	newUser := &User{
		Email:        email,
		PasswordHash: string(hash),
		Role:         "user", // Rol por defecto
		IsVerified:   false,  // Requerirá verificación
		CreatedAt:    time.Now(),
	}

	// Guardar en almacenamiento
	if err := s.store.CreateUser(newUser); err != nil {
		return fmt.Errorf("error al crear usuario: %w", err)
	}

	// Enviar email de verificación
	sendVerificationEmail(newUser.Email)
	return nil
}

// Login verifica credenciales y genera JWT
func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.store.GetUserByEmail(email)
	if err != nil {
		return "", errors.New("credenciales inválidas")
	}

	// Verificar si el email está confirmado
	if !user.IsVerified {
		return "", errors.New("cuenta no verificada")
	}

	// Comparar contraseñas hasheadas
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("credenciales inválidas")
	}

	// Generar token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})

	// Firmar token con secreto
	tokenString, err := token.SignedString([]byte("secreto_muy_seguro"))
	if err != nil {
		return "", fmt.Errorf("error al generar token: %w", err)
	}

	return tokenString, nil
}

// VerifyEmail actualiza el estado de verificación
func (s *AuthService) VerifyEmail(email string) error {
	user, err := s.store.GetUserByEmail(email)
	if err != nil {
		return errors.New("usuario no encontrado")
	}

	user.IsVerified = true
	if err := s.store.UpdateUser(user); err != nil {
		return fmt.Errorf("error al actualizar usuario: %w", err)
	}

	return nil
}