package main

import (
	"log"
	"net/http"
	
	"streamvault/internal/auth"
	"streamvault/internal/content"
	"streamvault/internal/web"
)

func main() {
	// Configurar servicios
	authStore := auth.NewMemoryUserStore()
	authService := auth.NewAuthService(authStore)
	
	// Registrar usuario administrador de ejemplo
	if err := authService.Register("admin@streamvault.com", "AdminPass123!", "admin"); err != nil {
		log.Printf("Error registrando admin: %v", err)
	}
	authService.VerifyEmail("admin@streamvault.com")
	
	// Configurar servicio de contenido
	contentStore := content.NewMemoryVideoStore()
	contentService, err := content.NewVideoService(content.VideoServiceConfig{
		StoragePath:   "./videos",
		MetadataStore: contentStore,
		Transcoder:    &content.MockTranscoder{},
		EncryptionKey: "clave-secreta-de-32-bytes-123456",
	})
	if err != nil {
		log.Fatal(err)
	}
	
	// Configurar servidor web
	webHandlers, err := web.NewWebHandlers(authService, contentService)
	if err != nil {
		log.Fatal(err)
	}
	
	router := web.NewRouter(webHandlers)
	
	// Iniciar servidor
	log.Println("Servidor iniciado en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}