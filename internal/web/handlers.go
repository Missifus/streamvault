package web

import (
	"html/template"
	"net/http"
	"path/filepath"
	
	"streamvault/internal/auth"
	"streamvault/internal/content"
)

type WebHandlers struct {
	authService    *auth.AuthService
	contentService *content.VideoService
	templates      map[string]*template.Template
}

func NewWebHandlers(authService *auth.AuthService, contentService *content.VideoService) (*WebHandlers, error) {
	// Cargar plantillas
	templates := make(map[string]*template.Template)
	
	tmplFiles := []struct {
		name     string
		files    []string
	}{
		{"login", []string{"templates/base.html", "templates/login.html"}},
		{"register", []string{"templates/base.html", "templates/register.html"}},
		{"videos", []string{"templates/base.html", "templates/videos.html"}},
		{"player", []string{"templates/base.html", "templates/player.html"}},
	}
	
	for _, tmpl := range tmplFiles {
		t, err := template.ParseFiles(tmpl.files...)
		if err != nil {
			return nil, err
		}
		templates[tmpl.name] = t
	}
	
	return &WebHandlers{
		authService:    authService,
		contentService: contentService,
		templates:      templates,
	}, nil
}

func (h *WebHandlers) LoginPage(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{"Title": "Iniciar Sesión"}
	h.templates["login"].ExecuteTemplate(w, "base", data)
}

func (h *WebHandlers) RegisterPage(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{"Title": "Registro"}
	h.templates["register"].ExecuteTemplate(w, "base", data)
}

func (h *WebHandlers) VideosPage(w http.ResponseWriter, r *http.Request) {
	// Obtener usuario del contexto (seteado por middleware)
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

	data := map[string]interface{}{
		"Title":  "Mis Videos",
		"User":   user,
		"Videos": videos,
	}

	h.templates["videos"].ExecuteTemplate(w, "base", data)
}

func (h *WebHandlers) PlayerPage(w http.ResponseWriter, r *http.Request) {
	videoID := r.PathValue("id")
	video, err := h.contentService.GetVideoMetadata(videoID)
	if err != nil {
		http.Error(w, "Video no encontrado", http.StatusNotFound)
		return
	}

	user, _ := r.Context().Value("user").(*auth.User)
	
	data := map[string]interface{}{
		"Title": video.Title,
		"User":  user,
		"Video": video,
	}

	h.templates["player"].ExecuteTemplate(w, "base", data)
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
		data := map[string]interface{}{
			"Title": "Iniciar Sesión",
			"Error": "Credenciales inválidas",
		}
		h.templates["login"].ExecuteTemplate(w, "base", data)
		return
	}

	// Establecer cookie de sesión
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
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
		data := map[string]interface{}{
			"Title": "Registro",
			"Error": "Error en el registro: " + err.Error(),
		}
		h.templates["register"].ExecuteTemplate(w, "base", data)
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
		MaxAge:   -1,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *WebHandlers) StreamVideo(w http.ResponseWriter, r *http.Request) {
	videoID := r.PathValue("id")
	
	w.Header().Set("Content-Type", "video/mp4")
	w.Header().Set("Cache-Control", "public, max-age=31536000")
	
	if err := h.contentService.StreamVideo(videoID, w); err != nil {
		http.Error(w, "Error transmitiendo video: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *WebHandlers) ServeStatic(w http.ResponseWriter, r *http.Request) {
	// Servir archivos estáticos
	http.ServeFile(w, r, filepath.Join("internal", "web", r.URL.Path))
}