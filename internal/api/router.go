// File: internal/api/router.go
package api

import (
	"net/http"
	"streamvault/internal/storage"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// App contiene las dependencias de la aplicación, como la conexión a la base de datos.
// Esto facilita las pruebas y la inyección de dependencias.
type App struct {
	Store     storage.DataStore // <-- AHORA DEPENDE DE LA INTERFAZ
	UploadDir string
	JwtSecret string
}

// NewRouter ahora devuelve un http.Handler para ser compatible con el middleware de CORS.
func NewRouter(app *App) http.Handler {
	r := mux.NewRouter()

	// Asignar los manejadores a la App para que tengan acceso a la BD.
	h := &handler{app: app}
	m := &middleware{app: app}
	// --- Rutas Públicas ---
	r.HandleFunc("/register", h.HandleRegisterUser).Methods("POST")
	r.HandleFunc("/login", h.HandleLoginUser).Methods("POST")
	r.HandleFunc("/verify", h.HandleVerifyEmail).Methods("POST")
	r.HandleFunc("/videos", h.HandleListVideos).Methods("GET")
	r.HandleFunc("/videos/{id:[0-9]+}", h.HandleGetVideoByID).Methods("GET") // NUEVA RUTA PÚBLICA

	fileServer := http.StripPrefix("/stream/", http.FileServer(http.Dir(app.UploadDir)))
	r.PathPrefix("/stream/").Handler(fileServer)

	// --- Rutas de Administrador Protegidas ---
	adminRoutes := r.PathPrefix("/admin").Subrouter()
	adminRoutes.Use(m.AuthMiddleware, m.AdminOnlyMiddleware)
	adminRoutes.HandleFunc("/upload", h.HandleUploadVideo).Methods("POST")
	// NUEVAS RUTAS DE ADMIN
	adminRoutes.HandleFunc("/videos/{id:[0-9]+}", h.HandleUpdateVideo).Methods("PUT")
	adminRoutes.HandleFunc("/videos/{id:[0-9]+}", h.HandleDeleteVideo).Methods("DELETE")

	// --- Configuración de CORS ---
	allowedOrigins := handlers.AllowedOrigins([]string{"*", "null"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Authorization", "Content-Type"})

	return handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders)(r)
}
