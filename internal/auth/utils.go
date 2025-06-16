package auth

import (
	"regexp"
	"strings"
)

// isValidEmail verifica si un email tiene formato válido
func isValidEmail(email string) bool {
	// Expresión regular simple para validar email
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(pattern, email)
	return match
}

// isStrongPassword verifica si una contraseña es segura
func isStrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false
	
	for _, c := range password {
		switch {
		case 'A' <= c && c <= 'Z':
			hasUpper = true
		case 'a' <= c && c <= 'z':
			hasLower = true
		case '0' <= c && c <= '9':
			hasDigit = true
		case strings.ContainsRune("!@#$%^&*()_+=-{}[]|\\:;\"'<>,.?/", c):
			hasSpecial = true
		}
	}
	
	// Requerir al menos 3 de los 4 tipos
	count := 0
	if hasUpper { count++ }
	if hasLower { count++ }
	if hasDigit { count++ }
	if hasSpecial { count++ }
	
	return count >= 3
}

// sendVerificationEmail simula el envío de email de verificación
func sendVerificationEmail(email string) {
	// En producción se enviaría un email real con un enlace de verificación
	// Esta es una implementación simulada para desarrollo
	println("\n=== EMAIL DE VERIFICACIÓN ===")
	println("Para: ", email)
	println("Asunto: Por favor verifica tu cuenta")
	println("Contenido: Haz clic en este enlace para verificar tu cuenta:")
	println("http://localhost:8080/verify?token=simulated-token&email=" + email)
	println("============================\n")
}