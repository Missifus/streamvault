// File: internal/api/handlers.go
package api

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"path/filepath"
	"strconv"
	"streamvault/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// handler es una estructura que contiene las dependencias de la aplicación,
// en este caso, la estructura 'App' que tiene el acceso a la base de datos y la configuración.
type handler struct {
	app *App
}

// --- FUNCIONES AUXILIARES ---

// generateSecureToken crea una cadena de texto aleatoria y segura para usarla como token.
// Es importante usar crypto/rand para asegurar que el token sea criptográficamente impredecible.
func generateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		// Si falla la generación de aleatoriedad, es un problema serio.
		// Devolvemos una cadena vacía para manejar el error.
		return ""
	}
	return hex.EncodeToString(b)
}

// sendVerificationEmail construye y envía el email de verificación usando una plantilla HTML.
func sendVerificationEmail(user *models.User, token string) error {
	// 1. Lee las credenciales y URLs desde las variables de entorno para mayor seguridad.
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	frontendURL := os.Getenv("FRONTEND_URL")

	// 2. Construye el enlace de verificación único para el usuario.
	verificationLink := fmt.Sprintf("%s#verify?token=%s", frontendURL, token)

	// 3. Define los datos que se insertarán en los marcadores de la plantilla HTML.
	data := struct {
		Username         string
		VerificationLink string
	}{
		Username:         user.Username,
		VerificationLink: verificationLink,
	}

	// 4. Lee y procesa el archivo de la plantilla HTML.
	t, err := template.ParseFiles("templates/verification_email.html")
	if err != nil {
		log.Printf("Error al parsear la plantilla de email: %v", err)
		return err
	}

	// 5. Prepara el cuerpo del correo en un buffer para un manejo eficiente de la memoria.
	var body bytes.Buffer
	// Se definen las cabeceras para que los clientes de correo lo interpreten como HTML.
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Verifica tu cuenta en StreamVault\r\n%s", mimeHeaders)))

	// Ejecuta la plantilla, rellenando los datos (Username y VerificationLink).
	if err := t.Execute(&body, data); err != nil {
		log.Printf("Error al ejecutar la plantilla de email: %v", err)
		return err
	}

	// 6. Se autentica con el servidor SMTP y envía el correo.
	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, []string{user.Email}, body.Bytes())
	if err != nil {
		log.Printf("Error al enviar email a %s: %v", user.Email, err)
		return err
	}

	log.Printf("Email de verificación con plantilla HTML enviado a %s", user.Email)
	return nil
}

// --- HANDLERS DE AUTENTICACIÓN Y USUARIOS ---

// HandleRegisterUser procesa el registro de un nuevo usuario.
func (h *handler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	// 1. Decodifica el cuerpo de la petición JSON en un struct de Usuario.
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Request inválido", http.StatusBadRequest)
		return
	}
	// 2. Valida que los campos necesarios no estén vacíos.
	if user.Username == "" || user.Email == "" || user.Password == "" {
		http.Error(w, "Faltan campos obligatorios", http.StatusBadRequest)
		return
	}
	// 3. Hashea la contraseña del usuario con bcrypt. NUNCA se debe guardar una contraseña en texto plano.
	// bcrypt es seguro porque es lento y resistente a ataques de fuerza bruta.
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	// 4. Asigna valores por defecto al nuevo usuario.
	user.Role = "user"
	user.IsVerified = false
	user.VerificationToken = generateSecureToken(32)

	// 5. Usa la interfaz DataStore para crear el usuario. El handler no sabe qué base de datos se usa (abstracción).
	if err := h.app.Store.CreateUser(&user); err != nil {
		http.Error(w, "El email o nombre de usuario ya está en uso", http.StatusInternalServerError)
		return
	}

	// 6. Ejecuta el envío de email en una goroutine (concurrencia).
	// Esto permite que la respuesta al usuario sea inmediata, sin esperar a que el email se envíe.
	go sendVerificationEmail(&user, user.VerificationToken)

	// 7. Responde al cliente con un mensaje de éxito.
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Registro exitoso. Por favor, revisa tu email para verificar tu cuenta.",
	})
}

// HandleLoginUser procesa el inicio de sesión.
func (h *handler) HandleLoginUser(w http.ResponseWriter, r *http.Request) {
	var reqUser models.User
	if err := json.NewDecoder(r.Body).Decode(&reqUser); err != nil {
		http.Error(w, "Request inválido", http.StatusBadRequest)
		return
	}

	// 1. Obtiene el usuario por su email desde la capa de datos.
	user, err := h.app.Store.GetUserByEmail(reqUser.Email)
	if err != nil {
		http.Error(w, "Email o contraseña incorrectos", http.StatusUnauthorized)
		return
	}

	// 2. Comprueba si la cuenta ha sido verificada. Es una capa de seguridad importante.
	if !user.IsVerified {
		http.Error(w, "Tu cuenta no ha sido verificada. Por favor, revisa tu email.", http.StatusForbidden)
		return
	}

	// 3. Compara de forma segura la contraseña enviada con el hash guardado en la BD.
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqUser.Password)); err != nil {
		http.Error(w, "Email o contraseña incorrectos", http.StatusUnauthorized)
		return
	}

	// 4. Si las credenciales son correctas, crea las "claims" para el token JWT.
	// Las claims son la información que irá dentro del token.
	claims := &models.Claims{
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // El token expira en 24 horas.
		},
	}

	// 5. Crea y firma el token con la clave secreta.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(h.app.JwtSecret))

	// 6. Envía el token al cliente.
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

// HandleVerifyEmail verifica la cuenta del usuario a través del token.
func (h *handler) HandleVerifyEmail(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Request inválido", http.StatusBadRequest)
		return
	}

	// 1. Busca al usuario en la BD usando el token de verificación.
	user, err := h.app.Store.VerifyUser(payload.Token)
	if err != nil {
		http.Error(w, "Token de verificación inválido o expirado.", http.StatusUnauthorized)
		return
	}

	// 2. Previene que una cuenta ya verificada se vuelva a verificar.
	if user.IsVerified {
		http.Error(w, "Esta cuenta ya ha sido verificada.", http.StatusBadRequest)
		return
	}

	// 3. Actualiza el estado del usuario a verificado y limpia el token.
	if err := h.app.Store.SetUserVerificationStatus(user.ID, true); err != nil {
		http.Error(w, "No se pudo verificar la cuenta.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "¡Cuenta verificada exitosamente! Ahora puedes iniciar sesión.",
	})
}

// --- HANDLERS DE GESTIÓN DE VIDEOS ---

// HandleListVideos obtiene y devuelve la lista de todos los videos.
func (h *handler) HandleListVideos(w http.ResponseWriter, r *http.Request) {
	videos, err := h.app.Store.GetAllVideos()
	if err != nil {
		http.Error(w, "No se pudieron obtener los videos", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(videos)
}

// HandleUploadVideo maneja la subida de un archivo de video.
func (h *handler) HandleUploadVideo(w http.ResponseWriter, r *http.Request) {
	// 1. Parsea el formulario, que puede contener tanto texto como archivos (multipart).
	// Se establece un límite de tamaño (ej. 50MB) para proteger al servidor.
	r.ParseMultipartForm(50 << 20)

	// 2. Obtiene el archivo del formulario. "video" es el 'name' del campo en el HTML.
	file, fileHandler, err := r.FormFile("video")
	if err != nil {
		http.Error(w, "Petición inválida: Falta el archivo 'video'", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 3. Genera un nombre de archivo único para evitar colisiones y guardarlo en el servidor.
	fileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filepath.Base(fileHandler.Filename))
	filePath := filepath.Join(h.app.UploadDir, fileName)

	// 4. Crea el archivo de destino en el sistema de archivos.
	dst, _ := os.Create(filePath)
	defer dst.Close()

	// 5. Copia el contenido del archivo subido al archivo de destino.
	io.Copy(dst, file)

	// 6. Crea un modelo de video con la información recibida.
	video := &models.Video{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Category:    r.FormValue("category"),
		FilePath:    fileName,
	}

	// 7. Guarda los metadatos del video en la base de datos a través de la interfaz.
	if err := h.app.Store.CreateVideo(video); err != nil {
		// Si falla la inserción en la BD, eliminamos el archivo físico para no dejar basura.
		os.Remove(filePath)
		http.Error(w, "Error al guardar la información del video", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(video)
}

// HandleGetVideoByID obtiene los detalles de un solo video por su ID.
func (h *handler) HandleGetVideoByID(w http.ResponseWriter, r *http.Request) {
	// 1. Obtiene las variables de la URL, en este caso el {id}.
	vars := mux.Vars(r)
	// 2. Convierte el ID de string a int.
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID de video inválido", http.StatusBadRequest)
		return
	}

	// 3. Obtiene el video desde la capa de datos.
	video, err := h.app.Store.GetVideoByID(id)
	if err != nil {
		http.Error(w, "Video no encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(video)
}

// HandleUpdateVideo actualiza los detalles de un video existente.
func (h *handler) HandleUpdateVideo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID de video inválido", http.StatusBadRequest)
		return
	}

	// 1. Se asegura de que el video a actualizar realmente existe.
	_, err = h.app.Store.GetVideoByID(id)
	if err != nil {
		http.Error(w, "Video no encontrado", http.StatusNotFound)
		return
	}

	// 2. Decodifica los nuevos datos del video desde el cuerpo de la petición.
	var updatedVideo models.Video
	if err := json.NewDecoder(r.Body).Decode(&updatedVideo); err != nil {
		http.Error(w, "Request inválido", http.StatusBadRequest)
		return
	}

	updatedVideo.ID = id // Se asegura de que el ID sea el correcto.

	// 3. Llama a la capa de datos para realizar la actualización.
	if err := h.app.Store.UpdateVideo(&updatedVideo); err != nil {
		log.Printf("Error al actualizar video: %v", err)
		http.Error(w, "Error al actualizar el video", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedVideo)
}

// HandleDeleteVideo elimina un video de la base de datos y del sistema de archivos.
func (h *handler) HandleDeleteVideo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID de video inválido", http.StatusBadRequest)
		return
	}

	// 1. Importante: Obtiene los datos del video ANTES de borrarlo de la BD
	// para poder saber el nombre del archivo físico a eliminar.
	video, err := h.app.Store.GetVideoByID(id)
	if err != nil {
		http.Error(w, "Video no encontrado", http.StatusNotFound)
		return
	}

	// 2. Elimina el registro del video de la base de datos.
	if err := h.app.Store.DeleteVideo(id); err != nil {
		log.Printf("Error al eliminar video de la BD: %v", err)
		http.Error(w, "Error al eliminar el video", http.StatusInternalServerError)
		return
	}

	// 3. Si la eliminación de la BD fue exitosa, elimina el archivo físico.
	filePath := filepath.Join(h.app.UploadDir, video.FilePath)
	if err := os.Remove(filePath); err != nil {
		// Este no es un error fatal para el cliente, pero se debe registrar en el log
		// para que un administrador pueda limpiar el archivo manualmente.
		log.Printf("ADVERTENCIA: No se pudo eliminar el archivo físico %s: %v", filePath, err)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Video eliminado exitosamente"})
}
