package content

import (
	"io"
	"time"
)

// VideoMetadata contiene metadatos de videos
type VideoMetadata struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	OwnerID     string `json:"owner_id"`
	FilePath    string `json:"file_path"` // Ruta relativa del video cifrado
	HLSPlaylist string `json:"hls_playlist"`
	CreatedAt   string `json:"created_at"`
}

// VideoStore define la interfaz para almacenamiento de metadatos de videos
type VideoStore interface {
	SaveMetadata(metadata *VideoMetadata) error
	GetVideoMetadata(id string) (*VideoMetadata, error)
	ListVideos(userID string) ([]*VideoMetadata, error)
}

// Transcoder define la interfaz para transcodificación de video
type Transcoder interface {
	Transcode(inputPath, outputDir string) error
}

// EncryptionService define la interfaz para cifrado/descifrado
type EncryptionService interface {
	EncryptFile(inputPath, outputPath string) error
	DecryptStream(r io.Reader, w io.Writer) error
}