package web

import (
	"log"
	"net/http"

	"streamvault/internal/auth"
	"streamvault/internal/content"
)

type WebServer struct {
	router *mux.Router
}

func NewWebServer(authService *auth.AuthService, contentService *content.VideoService) (*WebServer, error) {
	handlers, err := NewWebHandlers(authService, contentService)
	if err != nil {
		return nil, err
	}

	router := NewRouter(handlers)
	return &WebServer{router: router}, nil
}

func (s *WebServer) Start(addr string) error {
	log.Printf("Servidor web iniciado en %s", addr)
	return http.ListenAndServe(addr, s.router)
}