package auth

import (
	"encoding/json"
	"net/http"
)

// AuthHandlers maneja los endpoints HTTP para autenticación
type AuthHandlers struct {
	service *AuthService
}

// NewAuthHandlers crea nuevos manejadores de autenticación
func NewAuthHandlers(service *AuthService) *AuthHandlers {
	return &AuthHandlers{service: service}
}

// RegisterHandler maneja el registro de usuarios
func (h *AuthHandlers) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	if err := h.service.Register(req.Email, req.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Usuario registrado. Por favor verifica tu email."})
}

// LoginHandler maneja el inicio de sesión
func (h *AuthHandlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	token, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(LoginResponse{Token: token})
}

// VerifyHandler maneja la verificación de email
func (h *AuthHandlers) VerifyHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email requerido", http.StatusBadRequest)
		return
	}

	if err := h.service.VerifyEmail(email); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Email verificado exitosamente"})
}

// ProtectedHandler es un endpoint protegido de ejemplo
func (h *AuthHandlers) ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Acceso concedido a ruta protegida",
		"user_id": userID,
	})
}

// AdminHandler es un endpoint solo para administradores
func (h *AuthHandlers) AdminHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Acceso concedido a ruta de administrador",
		"user_id": userID,
	})
}