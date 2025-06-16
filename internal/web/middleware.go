package web

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

func (h *WebHandlers) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Rutas públicas que no requieren autenticación
		publicRoutes := []string{
			"/login",
			"/register",
			"/static/",
		}
		
		for _, path := range publicRoutes {
			if strings.HasPrefix(r.URL.Path, path) {
				next.ServeHTTP(w, r)
				return
			}
		}

		// Verificar cookie de sesión
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Validar token JWT
		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			return []byte("secreto_muy_seguro"), nil
		})

		if err != nil || !token.Valid {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Obtener usuario del token
		user := &auth.User{
			ID:    int(claims["sub"].(float64)),
			Email: claims["email"].(string),
			Role:  claims["role"].(string),
		}

		// Agregar usuario al contexto
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}