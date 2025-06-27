package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"streamvault/internal/models"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type handler struct {
	app *App
}

// HandleRegisterUser maneja el registro de nuevos usuarios.
func (h *handler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Request inválido", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error al procesar la contraseña", http.StatusInternalServerError)
		return
	}

	user.Role = "user"
	query := `INSERT INTO users (username, password_hash, role) VALUES ($1, $2, $3) RETURNING id`
	err = h.app.DB.QueryRow(query, user.Username, string(hashedPassword), user.Role).Scan(&user.ID)
	if err != nil {
		log.Printf("Error al registrar usuario: %v", err)
		http.Error(w, "El nombre de usuario ya existe o hubo un error en la base de datos", http.StatusInternalServerError)
		return
	}

	user.Password = "" // No devolver el hash de la contraseña
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// HandleLoginUser maneja el inicio de sesión y la generación de JWT.
func (h *handler) HandleLoginUser(w http.ResponseWriter, r *http.Request) {
	var credentials models.User
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Request inválido", http.StatusBadRequest)
		return
	}

	var user models.User
	var hashedPassword string
	query := `SELECT id, username, password_hash, role FROM users WHERE username=$1`
	err := h.app.DB.QueryRow(query, credentials.Username).Scan(&user.ID, &user.Username, &hashedPassword, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Usuario o contraseña incorrectos", http.StatusUnauthorized)
		} else {
			log.Printf("Error al buscar usuario '%s': %v", credentials.Username, err)
			http.Error(w, "Error interno del servidor", http.StatusInternalServerError)
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(credentials.Password)); err != nil {
		http.Error(w, "Usuario o contraseña incorrectos", http.StatusUnauthorized)
		return
	}

	// Generar el token JWT
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &models.Claims{
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(h.app.JwtSecret))
	if err != nil {
		log.Printf("Error al firmar el token para el usuario '%s': %v", user.Username, err)
		http.Error(w, "Error al generar el token", http.StatusInternalServerError)
		return
	}

	// Escribir el token en la respuesta JSON.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}

// HandleListVideos lista todos los videos, con opción de filtrar por categoría.
func (h *handler) HandleListVideos(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")

	var rows *sql.Rows
	var err error

	query := `SELECT id, title, description, category, file_path FROM videos`
	if category != "" {
		query += " WHERE category = $1"
		rows, err = h.app.DB.Query(query, category)
	} else {
		rows, err = h.app.DB.Query(query)
	}

	if err != nil {
		log.Printf("Error al consultar videos: %v", err)
		http.Error(w, "Error al consultar los videos", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	videos := []models.Video{}
	for rows.Next() {
		var v models.Video
		if err := rows.Scan(&v.ID, &v.Title, &v.Description, &v.Category, &v.FilePath); err != nil {
			log.Printf("Error al escanear fila de video: %v", err)
			http.Error(w, "Error al procesar los videos", http.StatusInternalServerError)
			return
		}
		videos = append(videos, v)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(videos)
}

// HandleUploadVideo se encarga de la subida de archivos de video por parte de un admin.
func (h *handler) HandleUploadVideo(w http.ResponseWriter, r *http.Request) {
	// Aumentamos el límite por si acaso, ej. 50 MB.
	if err := r.ParseMultipartForm(50 << 20); err != nil {
		log.Printf("Error al parsear el formulario multipart: %v", err)
		http.Error(w, "El archivo es demasiado grande. Límite: 50MB.", http.StatusBadRequest)
		return
	}

	// --- INICIO: CÓDIGO DE DEPURACIÓN ---
	// Este bloque imprimirá en la consola del servidor exactamente lo que está recibiendo.
	log.Println("--- INICIO DEPURACIÓN DE SUBIDA ---")
	log.Println("Encabezado Content-Type:", r.Header.Get("Content-Type"))
	log.Println("Formulario parseado. Valores de texto recibidos:")
	for key, values := range r.MultipartForm.Value {
		log.Printf("  - Clave (Texto): '%s', Valores: %v\n", key, values)
	}
	log.Println("Archivos recibidos:")
	for key := range r.MultipartForm.File {
		log.Printf("  - Clave (Archivo): '%s'\n", key)
	}
	log.Println("--- FIN DEPURACIÓN DE SUBIDA ---")
	// --- FIN: CÓDIGO DE DEPURACIÓN ---

	file, handler, err := r.FormFile("video")
	if err != nil {
		// Añadimos más contexto al error que se muestra en la terminal.
		log.Printf("Error en r.FormFile('video'): %v. El cliente no envió un archivo con el nombre 'video' o hubo otro error.", err)
		http.Error(w, "Petición inválida: Falta el archivo con el nombre 'video'", http.StatusBadRequest)
		return
	}
	defer file.Close()

	title := r.FormValue("title")
	description := r.FormValue("description")
	category := r.FormValue("category")

	if title == "" || category == "" {
		http.Error(w, "Faltan los campos 'title' o 'category'", http.StatusBadRequest)
		return
	}

	fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(handler.Filename))
	filePath := filepath.Join(h.app.UploadDir, fileName)

	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Error interno: No se pudo crear el archivo en el servidor", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Error interno: No se pudo guardar el archivo", http.StatusInternalServerError)
		return
	}

	var videoID int
	query := `INSERT INTO videos (title, description, category, file_path) VALUES ($1, $2, $3, $4) RETURNING id`
	err = h.app.DB.QueryRow(query, title, description, category, fileName).Scan(&videoID)
	if err != nil {
		os.Remove(filePath)
		log.Printf("Error al guardar video en BD: %v", err)
		http.Error(w, "Error interno: No se pudo guardar la información en la base de datos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Video subido exitosamente",
		"videoId":  videoID,
		"filePath": fileName,
	})
}
