package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"          // Enrutador HTTP poderoso
	"streamvault/internal/auth"      // Tu paquete de autenticación
)

func main() {
	// Configurar almacenamiento
	store := auth.NewMemoryUserStore()
	
	// Crear servicio de autenticación
	authService := auth.NewAuthService(store)
	
	// Crear manejadores HTTP
	authHandlers := auth.NewAuthHandlers(authService)
	
	// Configurar router
	router := mux.NewRouter()
	
	// Rutas públicas
	router.HandleFunc("/register", authHandlers.RegisterHandler).Methods("POST")
	router.HandleFunc("/login", authHandlers.LoginHandler).Methods("POST")
	router.HandleFunc("/verify", authHandlers.VerifyHandler).Methods("GET")
	
	// Ruta protegida para usuarios normales
	protectedRouter := router.PathPrefix("/protected").Subrouter()
	protectedRouter.Use(authService.AuthMiddleware(""))
	protectedRouter.HandleFunc("", authHandlers.ProtectedHandler).Methods("GET")
	
	// Ruta solo para administradores
	adminRouter := router.PathPrefix("/admin").Subrouter()
	adminRouter.Use(authService.AuthMiddleware("admin"))
	adminRouter.HandleFunc("", authHandlers.AdminHandler).Methods("GET")
	
	// Mensaje de inicio
	log.Println("Servidor de autenticación iniciado en :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}