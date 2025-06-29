package storage

import (
	"database/sql"
	"fmt"
	"streamvault/internal/models"

	_ "github.com/lib/pq"
)

/*
DataStore es la INTERFAZ que define el "contrato" para nuestro almacenamiento de datos.
Cualquier tipo que implemente todos estos métodos se considera un 'DataStore'.
Esta es la clave para la abstracción y el polimorfismo: los handlers dependerán
de esta interfaz, no de una base de datos específica como PostgreSQL.
*/

type DataStore interface {
	// Métodos de Usuario
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	VerifyUser(token string) (*models.User, error)
	SetUserVerificationStatus(userID int, isVerified bool) error
	// Métodos de Video
	CreateVideo(video *models.Video) error
	GetAllVideos() ([]*models.Video, error)
	GetVideoByID(id int) (*models.Video, error)
	UpdateVideo(video *models.Video) error
	DeleteVideo(id int) error
}

// PostgresStore es la IMPLEMENTACIÓN CONCRETA de la interfaz DataStore.
// Contiene la conexión a la base de datos PostgreSQL.
type PostgresStore struct {
	db *sql.DB
}

// NewPostgresStore es una función "constructora" que crea y devuelve una nueva instancia de PostgresStore.
// Recibe la cadena de conexión, establece la conexión con la base de datos y la verifica.
func NewPostgresStore(connStr string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	// db.Ping() es crucial para verificar que la conexión a la base de datos es válida.
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStore{db: db}, nil
}

// --- Implementaciones de Métodos de Usuario ---

// CreateUser implementa el método de la interfaz para guardar un nuevo usuario en la BD.
// Usa RETURNING en la consulta SQL para obtener de vuelta el ID y la fecha de creación generados por la BD.
func (s *PostgresStore) CreateUser(user *models.User) error {
	query := `INSERT INTO users (username, email, password_hash, role, is_verified, verification_token) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created_at`
	// QueryRow se usa porque esperamos que la consulta devuelva exactamente una fila.
	// Los argumentos (user.Username, etc.) se pasan de forma segura para prevenir inyección SQL.
	return s.db.QueryRow(query, user.Username, user.Email, user.Password, user.Role, user.IsVerified, user.VerificationToken).Scan(&user.ID, &user.CreatedAt)
}

// GetUserByEmail implementa la búsqueda de un usuario por su email.
func (s *PostgresStore) GetUserByEmail(email string) (*models.User, error) {
	user := new(models.User)
	query := `SELECT id, username, email, password_hash, role, is_verified FROM users WHERE email = $1`
	err := s.db.QueryRow(query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &user.IsVerified)
	if err != nil {
		// Es importante manejar el caso 'sql.ErrNoRows' para saber si el usuario simplemente no existe.
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("usuario no encontrado")
		}
		return nil, err
	}
	return user, nil
}

// VerifyUser busca un usuario por su token de verificación.
func (s *PostgresStore) VerifyUser(token string) (*models.User, error) {
	user := new(models.User)
	query := `SELECT id, is_verified FROM users WHERE verification_token = $1`
	err := s.db.QueryRow(query, token).Scan(&user.ID, &user.IsVerified)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("token inválido o expirado")
		}
		return nil, err
	}
	return user, nil
}

// SetUserVerificationStatus actualiza el estado de verificación de un usuario y limpia el token.
func (s *PostgresStore) SetUserVerificationStatus(userID int, isVerified bool) error {
	// Se establece verification_token a NULL para que no pueda ser reutilizado.
	query := `UPDATE users SET is_verified = $1, verification_token = NULL WHERE id = $2`
	// db.Exec se usa para consultas que no devuelven filas (UPDATE, DELETE, etc.).
	_, err := s.db.Exec(query, isVerified, userID)
	return err
}

// --- Implementaciones de Métodos de Video ---

// CreateVideo implementa el método para guardar un nuevo video en la BD.
func (s *PostgresStore) CreateVideo(video *models.Video) error {
	query := `INSERT INTO videos (title, description, category, file_path) VALUES ($1, $2, $3, $4) RETURNING id, uploaded_at`
	return s.db.QueryRow(query, video.Title, video.Description, video.Category, video.FilePath).Scan(&video.ID, &video.UploadedAt)
}

// GetAllVideos implementa el método para obtener todos los videos de la BD.
func (s *PostgresStore) GetAllVideos() ([]*models.Video, error) {
	query := `SELECT id, title, description, category, file_path, uploaded_at FROM videos`
	// db.Query se usa porque esperamos múltiples filas como resultado.
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	// Es fundamental cerrar las filas después de usarlas para liberar la conexión.
	defer rows.Close()

	var videos []*models.Video
	// Se itera sobre cada fila del resultado.
	for rows.Next() {
		video := new(models.Video)
		// Se escanean los valores de la fila actual en el struct de video.
		if err := rows.Scan(&video.ID, &video.Title, &video.Description, &video.Category, &video.FilePath, &video.UploadedAt); err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}
	return videos, nil
}

// GetVideoByID implementa la búsqueda de un video por su ID.
func (s *PostgresStore) GetVideoByID(id int) (*models.Video, error) {
	video := new(models.Video)
	query := `SELECT id, title, description, category, file_path, uploaded_at FROM videos WHERE id = $1`
	err := s.db.QueryRow(query, id).Scan(&video.ID, &video.Title, &video.Description, &video.Category, &video.FilePath, &video.UploadedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("video no encontrado")
		}
		return nil, err
	}
	return video, nil
}

// UpdateVideo implementa la actualización de los detalles de un video.
func (s *PostgresStore) UpdateVideo(video *models.Video) error {
	query := `UPDATE videos SET title = $1, description = $2, category = $3 WHERE id = $4`
	_, err := s.db.Exec(query, video.Title, video.Description, video.Category, video.ID)
	return err
}

// DeleteVideo implementa la eliminación de un video por su ID.
func (s *PostgresStore) DeleteVideo(id int) error {
	query := `DELETE FROM videos WHERE id = $1`
	_, err := s.db.Exec(query, id)
	return err
}
