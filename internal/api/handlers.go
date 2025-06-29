// El paquete 'api' contiene toda la lógica relacionada con el manejo de las
// peticiones HTTP, la definición de las rutas y las respuestas al cliente.
package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"streamvault/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// handler es una estructura que encapsula la aplicación.
// Los métodos de esta estructura (los manejadores) tendrán acceso a todas las dependencias de App.
type handler struct {
	app *App
}

// --- FUNCIONES AUXILIARES DE RESPUESTA ---

// respondWithError es una función auxiliar para enviar respuestas de error estandarizadas en formato JSON.
// Esto asegura que el frontend siempre reciba un JSON válido, incluso si hay un error.
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"message": message})
}

// respondWithJSON es una función auxiliar para serializar el payload y enviarlo como una respuesta JSON.
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// --- FUNCIÓN AUXILIAR PARA TAREAS EN SEGUNDO PLANO ---

// processVideoInBackground simula una tarea de larga duración que se ejecuta en segundo plano.
// Es nuestra demostración del concepto de CONCURRENCIA.
func processVideoInBackground(videoID int) {
	// Se registra el inicio de la tarea en la consola del servidor.
	log.Printf("[Video ID: %d] Iniciando procesamiento en segundo plano (ej: generar miniaturas)...", videoID)
	// time.Sleep simula una tarea que consume tiempo, como podría ser la transcodificación de un video
	// o la generación de diferentes calidades.
	time.Sleep(10 * time.Second)
	// Se registra la finalización de la tarea.
	log.Printf("[Video ID: %d] ...Procesamiento en segundo plano finalizado.", videoID)
}

// --- HANDLERS DE AUTENTICACIÓN Y USUARIOS ---

// HandleRegisterUser procesa el registro de un nuevo usuario.
func (h *handler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Request inválido")
		return
	}
	if user.Username == "" || user.Email == "" || user.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Faltan campos obligatorios")
		return
	}
	// Hashea la contraseña del usuario con bcrypt. NUNCA se debe guardar una contraseña en texto plano.
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	user.Role = "user"

	// Usa la interfaz DataStore para crear el usuario. El handler no sabe qué base de datos se usa (abstracción).
	if err := h.app.Store.CreateUser(&user); err != nil {
		respondWithError(w, http.StatusInternalServerError, "El email o nombre de usuario ya está en uso")
		return
	}
	respondWithJSON(w, http.StatusCreated, map[string]string{"message": "Registro exitoso. ¡Ahora puedes iniciar sesión!"})
}

// HandleLoginUser procesa el inicio de sesión de un usuario existente.
func (h *handler) HandleLoginUser(w http.ResponseWriter, r *http.Request) {
	var reqUser models.User
	if err := json.NewDecoder(r.Body).Decode(&reqUser); err != nil {
		respondWithError(w, http.StatusBadRequest, "Request inválido")
		return
	}
	// Obtiene el usuario por su email desde la capa de datos.
	user, err := h.app.Store.GetUserByEmail(reqUser.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Email o contraseña incorrectos")
		return
	}
	// Compara de forma segura la contraseña enviada con el hash guardado en la BD.
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqUser.Password)); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Email o contraseña incorrectos")
		return
	}
	// Si las credenciales son correctas, crea las "claims" para el token JWT.
	claims := &models.Claims{
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	// Crea y firma el token con la clave secreta.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(h.app.JwtSecret))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error interno al generar el token")
		return
	}
	// Envía el token al cliente.
	respondWithJSON(w, http.StatusOK, map[string]string{"token": tokenString})
}

// --- HANDLERS PÚBLICOS DE VIDEOS ---

// HandleListVideos obtiene y devuelve la lista de todos los videos.
func (h *handler) HandleListVideos(w http.ResponseWriter, r *http.Request) {
	videos, err := h.app.Store.GetAllVideos()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "No se pudieron obtener los videos")
		return
	}
	respondWithJSON(w, http.StatusOK, videos)
}

// HandleGetVideoByID obtiene los detalles de un solo video por su ID.
func (h *handler) HandleGetVideoByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "ID de video inválido")
		return
	}
	video, err := h.app.Store.GetVideoByID(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Video no encontrado")
		return
	}
	respondWithJSON(w, http.StatusOK, video)
}

// HandleStreamVideo sirve el contenido de un video para su reproducción.
func (h *handler) HandleStreamVideo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileName := vars["filename"]
	videoPath := filepath.Join(h.app.UploadDir, fileName)
	// http.ServeFile es una función de Go que se encarga de servir un archivo.
	// Soporta 'Range requests', crucial para que los navegadores puedan buscar (seek) en el video.
	http.ServeFile(w, r, videoPath)
}

// --- HANDLERS DE ADMINISTRACIÓN (PROTEGIDOS) ---

// HandleUploadVideo maneja la subida de un archivo de video.
func (h *handler) HandleUploadVideo(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(50 << 20); err != nil {
		respondWithError(w, http.StatusBadRequest, "El archivo es demasiado grande (límite 50MB)")
		return
	}
	title := r.FormValue("title")
	category := r.FormValue("category")
	if title == "" || category == "" {
		respondWithError(w, http.StatusBadRequest, "Faltan los campos 'title' o 'category'")
		return
	}
	file, fileHandler, err := r.FormFile("video")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Petición inválida: Falta el archivo con la clave 'video'")
		return
	}
	defer file.Close()
	// Genera un nombre de archivo único para evitar colisiones.
	fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(fileHandler.Filename))
	filePath := filepath.Join(h.app.UploadDir, fileName)
	dst, err := os.Create(filePath)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error interno al guardar el archivo")
		return
	}
	defer dst.Close()
	if _, err := io.Copy(dst, file); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error interno al procesar el archivo")
		return
	}
	video := &models.Video{
		Title:       title,
		Description: r.FormValue("description"),
		Category:    category,
		FilePath:    fileName,
	}
	// Guarda los metadatos en la base de datos a través de la interfaz.
	if err := h.app.Store.CreateVideo(video); err != nil {
		os.Remove(filePath) // Limpia el archivo si la BD falla.
		respondWithError(w, http.StatusInternalServerError, "Error al guardar la información del video")
		return
	}
	// Inicia una tarea en segundo plano (goroutine) para "procesar" el video.
	go processVideoInBackground(video.ID)
	respondWithJSON(w, http.StatusCreated, video)
}

// HandleUpdateVideo actualiza los detalles de un video existente.
func (h *handler) HandleUpdateVideo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "ID de video inválido")
		return
	}
	var updatedVideo models.Video
	if err := json.NewDecoder(r.Body).Decode(&updatedVideo); err != nil {
		respondWithError(w, http.StatusBadRequest, "Request inválido")
		return
	}
	updatedVideo.ID = id
	if err := h.app.Store.UpdateVideo(&updatedVideo); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error al actualizar el video")
		return
	}
	respondWithJSON(w, http.StatusOK, updatedVideo)
}

// HandleDeleteVideo elimina un video de la base de datos y del sistema de archivos.
func (h *handler) HandleDeleteVideo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "ID de video inválido")
		return
	}
	// Obtiene los datos del video ANTES de borrarlo de la BD para saber el nombre del archivo.
	video, err := h.app.Store.GetVideoByID(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Video no encontrado")
		return
	}
	// Elimina el registro de la BD.
	if err := h.app.Store.DeleteVideo(id); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error al eliminar el video")
		return
	}
	// Si la eliminación de la BD fue exitosa, elimina el archivo físico.
	filePath := filepath.Join(h.app.UploadDir, video.FilePath)
	os.Remove(filePath)
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Video eliminado exitosamente"})
}

// HandleListAllUsers devuelve una lista de todos los usuarios registrados (solo para admins).
func (h *handler) HandleListAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.app.Store.GetAllUsers()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error al obtener la lista de usuarios")
		return
	}
	respondWithJSON(w, http.StatusOK, users)
}

// HandleAdminUpdateUserRole actualiza el rol de un usuario (solo para admins).
func (h *handler) HandleAdminUpdateUserRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "ID de usuario inválido")
		return
	}
	var payload struct {
		Role string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		respondWithError(w, http.StatusBadRequest, "Cuerpo de la petición inválido")
		return
	}
	if payload.Role != "user" && payload.Role != "admin" {
		respondWithError(w, http.StatusBadRequest, "Rol inválido. Debe ser 'user' o 'admin'.")
		return
	}
	if err := h.app.Store.UpdateUserRole(id, payload.Role); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error al actualizar el rol del usuario")
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Rol del usuario actualizado exitosamente"})
}

// HandleAdminDeleteUser elimina un usuario del sistema (solo para admins).
func (h *handler) HandleAdminDeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "ID de usuario inválido")
		return
	}
	if err := h.app.Store.DeleteUser(id); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error al eliminar el usuario")
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Usuario eliminado exitosamente"})
}
