// File: internal/api/router.go
package api

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// App contiene las dependencias de la aplicación, como la conexión a la base de datos.
// Esto facilita las pruebas y la inyección de dependencias.
type App struct {
	DB        *sql.DB
	UploadDir string
	JwtSecret string
}

// NewRouter ahora devuelve un http.Handler para ser compatible con el middleware de CORS.
func NewRouter(app *App) http.Handler {
	r := mux.NewRouter()

	// Asignar los manejadores a la App para que tengan acceso a la BD.
	h := &handler{app: app}
	m := &middleware{app: app}

	// Rutas públicas
	r.HandleFunc("/register", h.HandleRegisterUser).Methods("POST")
	r.HandleFunc("/login", h.HandleLoginUser).Methods("POST")
	r.HandleFunc("/videos", h.HandleListVideos).Methods("GET")

	// Ruta para servir los archivos de video.
	fileServer := http.StripPrefix("/stream/", http.FileServer(http.Dir(app.UploadDir)))
	r.PathPrefix("/stream/").Handler(fileServer)

	// Rutas protegidas para administradores
	adminRoutes := r.PathPrefix("/admin").Subrouter()
	adminRoutes.Use(m.AuthMiddleware, m.AdminOnlyMiddleware)
	adminRoutes.HandleFunc("/upload", h.HandleUploadVideo).Methods("POST")

	// --- CONFIGURACIÓN DE CORS ---
	// Define de qué orígenes (dominios) se aceptarán las peticiones.
	// Usamos "*" para permitir desde cualquier origen durante el desarrollo.
	// "null" es importante para permitir peticiones desde archivos locales (file://).
	allowedOrigins := handlers.AllowedOrigins([]string{"*", "null"})

	// Define qué métodos HTTP están permitidos.
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"})

	// Define qué cabeceras HTTP puede enviar el cliente.
	// "Authorization" es crucial para el token JWT. "Content-Type" para los JSON.
	allowedHeaders := handlers.AllowedHeaders([]string{"Authorization", "Content-Type"})

	// Aplicamos el middleware de CORS a nuestro enrutador principal 'r'.
	// La función handlers.CORS() envuelve nuestro enrutador y maneja las peticiones OPTIONS y añade las cabeceras.
	return handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders)(r)
}
