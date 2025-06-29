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
	GetAllUsers() ([]models.User, error)
	DeleteUser(id int) error
	UpdateUserRole(id int, role string) error
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

func NewPostgresStore(connStr string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Init() error {
	// Usamos 'CREATE TABLE IF NOT EXISTS' para que la operación sea segura
	// y solo cree las tablas la primera vez que se ejecuta.
	createUsersTableSQL := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        username VARCHAR(50) UNIQUE NOT NULL,
        email VARCHAR(255) UNIQUE NOT NULL,
        password_hash VARCHAR(255) NOT NULL,
        role VARCHAR(20) NOT NULL CHECK (role IN ('user', 'admin')),
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );`

	createVideosTableSQL := `
    CREATE TABLE IF NOT EXISTS videos (
        id SERIAL PRIMARY KEY,
        title VARCHAR(255) NOT NULL,
        description TEXT,
        category VARCHAR(100) NOT NULL,
        file_path VARCHAR(255) NOT NULL,
        uploaded_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );`

	// Ejecutamos ambas consultas para crear las tablas.
	if _, err := s.db.Exec(createUsersTableSQL); err != nil {
		return fmt.Errorf("error al crear la tabla users: %w", err)
	}

	if _, err := s.db.Exec(createVideosTableSQL); err != nil {
		return fmt.Errorf("error al crear la tabla videos: %w", err)
	}

	return nil
}

// CreateUser se ha simplificado para coincidir con la nueva estructura de la BD.
func (s *PostgresStore) CreateUser(user *models.User) error {
	query := `INSERT INTO users (username, email, password_hash, role) VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	return s.db.QueryRow(query, user.Username, user.Email, user.Password, user.Role).Scan(&user.ID, &user.CreatedAt)
}

// GetUserByEmail se ha simplificado.
func (s *PostgresStore) GetUserByEmail(email string) (*models.User, error) {
	user := new(models.User)
	query := `SELECT id, username, email, password_hash, role FROM users WHERE email = $1`
	err := s.db.QueryRow(query, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("usuario no encontrado")
		}
		return nil, err
	}
	return user, nil
}

func (s *PostgresStore) GetAllUsers() ([]models.User, error) {
	query := `SELECT id, username, email, role, created_at FROM users`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (s *PostgresStore) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := s.db.Exec(query, id)
	return err
}

func (s *PostgresStore) UpdateUserRole(id int, role string) error {
	query := `UPDATE users SET role = $1 WHERE id = $2`
	_, err := s.db.Exec(query, role, id)
	return err
}

func (s *PostgresStore) CreateVideo(video *models.Video) error {
	query := `INSERT INTO videos (title, description, category, file_path) VALUES ($1, $2, $3, $4) RETURNING id, uploaded_at`
	return s.db.QueryRow(query, video.Title, video.Description, video.Category, video.FilePath).Scan(&video.ID, &video.UploadedAt)
}

func (s *PostgresStore) GetAllVideos() ([]*models.Video, error) {
	query := `SELECT id, title, description, category, file_path, uploaded_at FROM videos`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var videos []*models.Video
	for rows.Next() {
		video := new(models.Video)
		if err := rows.Scan(&video.ID, &video.Title, &video.Description, &video.Category, &video.FilePath, &video.UploadedAt); err != nil {
			return nil, err
		}
		videos = append(videos, video)
	}
	return videos, nil
}

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

func (s *PostgresStore) UpdateVideo(video *models.Video) error {
	query := `UPDATE videos SET title = $1, description = $2, category = $3 WHERE id = $4`
	_, err := s.db.Exec(query, video.Title, video.Description, video.Category, video.ID)
	return err
}

func (s *PostgresStore) DeleteVideo(id int) error {
	query := `DELETE FROM videos WHERE id = $1`
	_, err := s.db.Exec(query, id)
	return err
}
