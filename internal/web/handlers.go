package web

import (
	"html/template"
	"net/http"
	"path/filepath"
	"tuproyecto/internal/auth"
	"tuproyecto/internal/content"
)

type WebHandlers struct {
	authService    *auth.AuthService
	contentService *content.VideoService
	templates      *template.Template
}

type TemplateData struct {
	Title  string
	User   *auth.User
	Error  string
	Videos []*content.VideoMetadata
	Video  *content.VideoMetadata
}

func NewWebHandlers(authService *auth.AuthService, contentService *content.VideoService) (*WebHandlers, error) {
	// Cargar plantillas
	templates, err := template.ParseGlob("internal/web/templates/*.html")
	if err != nil {
		return nil, err
	}

	return &WebHandlers{
		authService:    authService,
		contentService: contentService,
		templates:      templates,
	}, nil
}

func (h *WebHandlers) LoginPage(w http.ResponseWriter, r *http.Request) {
	data := TemplateData{
		Title: "Iniciar Sesión",
	}
	
	if err := h.templates.ExecuteTemplate(w, "login.html", data); err != nil {
		http.Error(w, "Error al renderizar plantilla", http.StatusInternalServerError)
	}
}

func (h *WebHandlers) RegisterPage(w http.ResponseWriter, r *http.Request) {
	data := TemplateData{
		Title: "Registro",
	}
	
	if err := h.templates.ExecuteTemplate(w, "register.html", data); err != nil {
		http.Error(w, "Error al renderizar plantilla", http.StatusInternalServerError)
	}
}

func (h *WebHandlers) VideosPage(w http.ResponseWriter, r *http.Request) {
	// Obtener usuario del contexto (seteado por el middleware)
	user, ok := r.Context().Value("user").(*auth.User)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	videos, err := h.contentService.ListVideos(user.Email)
	if err != nil {
		http.Error(w, "Error al obtener videos", http.StatusInternalServerError)
		return
	}

	data := TemplateData{
		Title:  "Mis Videos",
		User:   user,
		Videos: videos,
	}

	if err := h.templates.ExecuteTemplate(w, "videos.html", data); err != nil {
		http.Error(w, "Error al renderizar plantilla", http.StatusInternalServerError)
	}
}

func (h *WebHandlers) PlayerPage(w http.ResponseWriter, r *http.Request) {
	videoID := r.PathValue("id")
	video, err := h.contentService.GetVideoMetadata(videoID)
	if err != nil {
		http.Error(w, "Video no encontrado", http.StatusNotFound)
		return
	}

	data := TemplateData{
		Title: video.Title,
		User:  r.Context().Value("user").(*auth.User),
		Video: video,
	}

	if err := h.templates.ExecuteTemplate(w, "player.html", data); err != nil {
		http.Error(w, "Error al renderizar plantilla", http.StatusInternalServerError)
	}
}

func (h *WebHandlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error procesando formulario", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	token, err := h.authService.Login(email, password)
	if err != nil {
		data := TemplateData{
			Title: "Iniciar Sesión",
			Error: "Credenciales inválidas",
		}
		
		h.templates.ExecuteTemplate(w, "login.html", data)
		return
	}

	// Establecer cookie de sesión
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // En producción debe ser true (HTTPS)
		SameSite: http.SameSiteStrictMode,
	})

	http.Redirect(w, r, "/videos", http.StatusSeeOther)
}

func (h *WebHandlers) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error procesando formulario", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	if err := h.authService.Register(email, password, "user"); err != nil {
		data := TemplateData{
			Title: "Registro",
			Error: "Error en el registro: " + err.Error(),
		}
		
		h.templates.ExecuteTemplate(w, "register.html", data)
		return
	}

	// Enviar a login después de registro exitoso
	http.Redirect(w, r, "/login?registered=true", http.StatusSeeOther)
}

func (h *WebHandlers) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Eliminar cookie de sesión
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   -1, // Eliminar la cookie
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *WebHandlers) StreamVideo(w http.ResponseWriter, r *http.Request) {
	videoID := r.PathValue("id")
	
	w.Header().Set("Content-Type", "video/mp4")
	w.Header().Set("Cache-Control", "public, max-age=31536000") // Cache por 1 año
	
	if err := h.contentService.StreamVideo(videoID, w); err != nil {
		http.Error(w, "Error transmitiendo video: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *WebHandlers) ServeStatic(w http.ResponseWriter, r *http.Request) {
	// Servir archivos estáticos desde el directorio
	fs := http.FileServer(http.Dir("internal/web/static"))
	fs.ServeHTTP(w, r)
}