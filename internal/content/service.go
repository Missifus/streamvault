package content

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

// VideoService maneja las operaciones con videos
type VideoService struct {
	storagePath   string
	encryptionKey []byte
}

// NewVideoService crea un nuevo servicio de videos
func NewVideoService(storagePath string, encryptionKey string) (*VideoService, error) {
	if len(encryptionKey) != 32 {
		return nil, errors.New("la clave de cifrado debe tener 32 bytes")
	}

	return &VideoService{
		storagePath:   storagePath,
		encryptionKey: []byte(encryptionKey),
	}, nil
}

// UploadVideo procesa y almacena un nuevo video
func (s *VideoService) UploadVideo(
	file multipart.File,
	header *multipart.FileHeader,
	metadata *VideoMetadata,
) error {
	// Generar ID único
	metadata.ID = uuid.New().String()

	// Crear directorio estructurado: /año/mes/día/id_video/
	dirPath := filepath.Join(s.storagePath, metadata.ID)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("error al crear directorio: %w", err)
	}

	// Guardar archivo original temporalmente
	tempPath := filepath.Join(dirPath, "original.tmp")
	if err := saveUploadedFile(file, tempPath); err != nil {
		return err
	}

	// Cifrar video
	encryptedPath := filepath.Join(dirPath, "encrypted.mp4")
	if err := s.encryptFile(tempPath, encryptedPath); err != nil {
		return fmt.Errorf("error cifrando video: %w", err)
	}
	metadata.FilePath = encryptedPath

	// Transcodificar a HLS
	hlsPath := filepath.Join(dirPath, "hls")
	if err := s.transcoder.Transcode(encryptedPath, hlsPath); err != nil {
	}
	metadata.HLSPlaylist = filepath.Join(hlsPath, "playlist.m3u8")
	// Limpiar archivo temporal
	os.Remove(tempPath)

	return nil
}

// StreamVideo entrega el video para transmisión
func (s *VideoService) StreamVideo(videoID string, w io.Writer) error {
	metadata, err := s.metadataStore.GetVideoMetadata(videoID) // Corregido: s.metadataStore
	if err != nil {
		return fmt.Errorf("error obteniendo metadatos: %w", err)
	}

	file, err := os.Open(metadata.FilePath)
	if err != nil {
		return fmt.Errorf("error abriendo archivo: %w", err)
	}
	defer file.Close()

	if err := s.encryptionService.DecryptStream(file, w); err != nil {
		return fmt.Errorf("error descifrando stream: %w", err)
	}

	return nil
}

// --- Funciones auxiliares ---

// encryptFile cifra un archivo usando AES-256
func (s *VideoService) encryptFile(inputPath, outputPath string) error {
	input, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(s.encryptionKey)
	if err != nil {
		return err
	}

	// El IV debe ser único pero no secreto
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return err
	}

	// Añadir IV al inicio del archivo cifrado
	output := make([]byte, len(input)+aes.BlockSize)
	copy(output[:aes.BlockSize], iv)

	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(output[aes.BlockSize:], input)

	return ioutil.WriteFile(outputPath, output, 0644)
}

// decryptStream descifra un stream para transmisión
func (s *VideoService) decryptStream(r io.Reader, w io.Writer) error {
	block, err := aes.NewCipher(s.encryptionKey)
	if err != nil {
		return err
	}

	// Leer IV del inicio del stream
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(r, iv); err != nil {
		return err
	}

	stream := cipher.NewCTR(block, iv)
	buf := make([]byte, 32*1024) // Buffer de 32KB

	for {
		n, err := r.Read(buf)
		if n > 0 {
			stream.XORKeyStream(buf[:n], buf[:n])
			if _, err := w.Write(buf[:n]); err != nil {
				return err
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// transcodeToHLS convierte un video a formato HLS (simulado)
func transcodeToHLS(inputPath, outputDir string) error {
	// En producción usaríamos ffmpeg:
	// ffmpeg -i input.mp4 -codec: copy -start_number 0 -hls_time 10 -hls_list_size 0 -f hls output.m3u8

	// Simulación: crear archivos HLS de ejemplo
	os.Mkdir(outputDir, 0755)
	playlist := `#EXTM3U
#EXT-X-VERSION:3
#EXT-X-TARGETDURATION:10
#EXT-X-MEDIA-SEQUENCE:0
#EXTINF:10.000000,
segment0.ts
#EXTINF:10.000000,
segment1.ts
#EXT-X-ENDLIST`

	if err := ioutil.WriteFile(filepath.Join(outputDir, "playlist.m3u8"), []byte(playlist), 0644); err != nil {
		return err
	}

	// Crear segmentos simulados
	for i := 0; i < 2; i++ {
		segmentPath := filepath.Join(outputDir, fmt.Sprintf("segment%d.ts", i))
		if err := ioutil.WriteFile(segmentPath, []byte(fmt.Sprintf("Segment %d data", i)), 0644); err != nil {
			return err
		}
	}
	return nil
}

func saveUploadedFile(file multipart.File, path string) error {
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	return err
}
