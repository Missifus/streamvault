// File: internal/api/middleware.go
package api

import (
	"context"
	"net/http"
	"strings"

	"streamvault/internal/models"

	"github.com/golang-jwt/jwt/v4"
)

// middleware encapsula la lógica de los middlewares.
type middleware struct {
	app *App
}

// AuthMiddleware verifica el JWT.
func (m *middleware) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Falta el encabezado de autorización", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &models.Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.app.JwtSecret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Token inválido", http.StatusUnauthorized)
			return
		}

		// Pasar los claims (información del usuario) al siguiente manejador a través del contexto.
		ctx := context.WithValue(r.Context(), "userClaims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AdminOnlyMiddleware verifica que el rol del usuario sea 'admin'.
func (m *middleware) AdminOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Se asume que AuthMiddleware ya se ejecutó y pobló el contexto.
		claims, ok := r.Context().Value("userClaims").(*models.Claims)
		if !ok || claims.Role != "admin" {
			http.Error(w, "Acceso denegado: se requiere rol de administrador", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
