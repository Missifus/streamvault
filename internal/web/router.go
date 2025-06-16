package web

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(handlers *WebHandlers) *mux.Router {
	r := mux.NewRouter()

	// Aplicar middleware de autenticación
	r.Use(handlers.AuthMiddleware)

	// Rutas públicas
	r.HandleFunc("/login", handlers.LoginPage).Methods("GET")
	r.HandleFunc("/register", handlers.RegisterPage).Methods("GET")
	r.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	r.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
	r.PathPrefix("/static/").HandlerFunc(handlers.ServeStatic)

	// Rutas protegidas
	r.HandleFunc("/videos", handlers.VideosPage).Methods("GET")
	r.HandleFunc("/player/{id}", handlers.PlayerPage).Methods("GET")
	r.HandleFunc("/videos/{id}", handlers.StreamVideo).Methods("GET")
	r.HandleFunc("/logout", handlers.LogoutHandler).Methods("GET")

	// Redirigir raíz a videos
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/videos", http.StatusFound)
	})

	return r
}