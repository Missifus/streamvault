package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt" // Paquete para trabajar con JWT
)

// AuthMiddleware crea un middleware para autenticación JWT
func (s *AuthService) AuthMiddleware(requiredRole string) func(http.Handler) http.Handler {
	// La función retornada es un constructor de middlewares
	return func(next http.Handler) http.Handler {
		//ejecutará en cada solicitud
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extraer token del encabezado Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "No autorizado", http.StatusUnauthorized)
				return
			}

			// El formato debe ser: Bearer <token>
			tokenString := ""
			if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				tokenString = authHeader[7:]
			}

			if tokenString == "" {
				http.Error(w, "Formato de token inválido", http.StatusUnauthorized)
				return
			}

			// Validar token
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Verificar el método de firma usado
			// jwt.SigningMethodHMAC: Algoritmos como HS256, HS384, HS512
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("método de firma inesperado: %v", token.Header["alg"])
				}
				return []byte("secreto_muy_seguro"), nil
			})
			// Verificar si el token es válido (no expirado, firma correcta)
			if err != nil || !token.Valid {
				http.Error(w, "Token inválido", http.StatusUnauthorized)
				return
			}
			// Extraer y verificar los claims (datos) del token
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Error en los claims del token", http.StatusUnauthorized)
				return
			}

			// Verificar rol requerido
			if requiredRole != "" {
				userRole, ok := claims["role"].(string)
				if !ok || userRole != requiredRole {
					http.Error(w, "Permisos insuficientes", http.StatusForbidden)
					return
				}
			}

			// Agregar información del usuario al contexto
			userID := int(claims["sub"].(float64))
			ctx := context.WithValue(r.Context(), "userID", userID)
			
			// Autenticación exitosa, continuar
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}