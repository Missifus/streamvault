package api

import (
	"log"
	"net/http"
	"streamvault/internal/storage"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// App contiene las dependencias de la aplicación que los manejadores necesitan para funcionar.
// Al inyectar dependencias de esta manera, el código se vuelve más modular y fácil de probar.

type App struct {
	Store                   storage.DataStore
	UploadDir               string
	JwtSecret               string
	EnableEmailVerification bool
}

// loggingMiddleware es un middleware simple que imprime en la consola cada petición recibida.
// Esto es increíblemente útil para depurar problemas de enrutamiento.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Imprime el método (GET, POST, etc.) y la ruta solicitada.
		log.Printf("Petición recibida: %s %s", r.Method, r.URL.Path)
		// Continúa con el siguiente manejador en la cadena.
		next.ServeHTTP(w, r)
	})
}

// NewRouter configura y devuelve un nuevo enrutador con todas las rutas de la API y el frontend.
func NewRouter(app *App) http.Handler {
	r := mux.NewRouter()

	// Aplicamos nuestro nuevo middleware de logging a todas las rutas.
	r.Use(loggingMiddleware)

	h := &handler{app: app}
	m := &middleware{app: app}

	// Agrupamos todas las rutas de la API bajo el prefijo "/api".
	apiRouter := r.PathPrefix("/api").Subrouter()

	// Definimos las rutas públicas de la API.
	apiRouter.HandleFunc("/register", h.HandleRegisterUser).Methods("POST")
	apiRouter.HandleFunc("/login", h.HandleLoginUser).Methods("POST")
	apiRouter.HandleFunc("/videos", h.HandleListVideos).Methods("GET")
	apiRouter.HandleFunc("/videos/{id:[0-9]+}", h.HandleGetVideoByID).Methods("GET")

	// Definimos las rutas de administrador protegidas.
	adminRoutes := apiRouter.PathPrefix("/admin").Subrouter()
	adminRoutes.Use(m.AuthMiddleware, m.AdminOnlyMiddleware)
	adminRoutes.HandleFunc("/upload", h.HandleUploadVideo).Methods("POST")
	adminRoutes.HandleFunc("/videos/{id:[0-9]+}", h.HandleUpdateVideo).Methods("PUT")
	adminRoutes.HandleFunc("/videos/{id:[0-9]+}", h.HandleDeleteVideo).Methods("DELETE")
	adminRoutes.HandleFunc("/users", h.HandleListAllUsers).Methods("GET")
	adminRoutes.HandleFunc("/users/{id:[0-9]+}/role", h.HandleAdminUpdateUserRole).Methods("PUT")
	adminRoutes.HandleFunc("/users/{id:[0-a-9]+}", h.HandleAdminDeleteUser).Methods("DELETE")
	// La ruta de streaming es una ruta especial para servir archivos.
	r.HandleFunc("/stream/{filename}", h.HandleStreamVideo).Methods("GET")

	// Servidor de archivos estáticos para el frontend (debe ser la última regla de enrutamiento).
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/")))

	// --- Configuración de CORS ---
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Authorization", "Content-Type"})

	return handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders)(r)
}
