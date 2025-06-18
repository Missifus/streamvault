package auth

import (
	"encoding/json"
	"net/http"
)

// AuthHandlers es una estructura que agrupa todos los manejadores HTTP
type AuthHandlers struct {
	service *AuthService
}

// NewAuthHandlers es el constructor que crea una nueva instancia de AuthHandlers.
// Recibe el servicio de autenticación como dependencia
func NewAuthHandlers(service *AuthService) *AuthHandlers {
	return &AuthHandlers{service: service}
}

// RegisterHandler maneja el registro de usuarios
func (h *AuthHandlers) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest // "Prepara un formulario vacío"
	// Leer y decodificar el cuerpo JSON de la solicitud
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}

	//Delegar la operación de registro al servicio
	if err := h.service.Register(req.Email, req.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// 4. Registro exitoso: responder con código 201
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Usuario registrado. Por favor verifica tu email."})
}

// LoginHandler maneja el inicio de sesión
func (h *AuthHandlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	//Decodificar el cuerpo JSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Solicitud inválida", http.StatusBadRequest)
		return
	}
	//Intentar autenticar al usuario
	token, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(LoginResponse{Token: token})
}

// VerifyHandler maneja la verificación de email
func (h *AuthHandlers) VerifyHandler(w http.ResponseWriter, r *http.Request) {
	//Obtener el email desde los parámetros de la URL
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email requerido", http.StatusBadRequest)
		return
	}
	//Solicitar al servicio que verifique el email
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