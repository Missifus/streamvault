package content

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// FFmpegTranscoder implementa Transcoder usando ffmpeg
type FFmpegTranscoder struct{}

func (t *FFmpegTranscoder) Transcode(inputPath, outputDir string) error {
	// Asegurar que el directorio de salida existe
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("no se pudo crear el directorio de salida: %w", err)
	}

	// Crear contexto con timeout para evitar procesos colgados
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	// Generar el comando ffmpeg para crear HLS
	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-i", inputPath,              // Archivo de entrada
		"-c:v", "libx264",            // Códec de video
		"-crf", "23",                 // Calidad (0-51, menor es mejor)
		"-preset", "medium",          // Compresión/velocidad
		"-c:a", "aac",                // Códec de audio
		"-b:a", "128k",               // Tasa de audio
		"-f", "hls",                  // Formato de salida HLS
		"-hls_time", "10",            // Duración de segmento (segundos)
		"-hls_list_size", "0",        // Máximo número de segmentos en playlist (0=infinito)
		"-hls_segment_filename", filepath.Join(outputDir, "segment_%03d.ts"), // Nombre de segmentos
		filepath.Join(outputDir, "playlist.m3u8"), // Archivo playlist
	)

	// Capturar stdout y stderr para diagnóstico
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Ejecutar el comando
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error en ffmpeg: %w", err)
	}

	return nil
}

// MockTranscoder implementa Transcoder para pruebas sin ffmpeg
type MockTranscoder struct{}

func (t *MockTranscoder) Transcode(inputPath, outputDir string) error {
	// Crear directorio de salida
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("no se pudo crear directorio HLS: %w", err)
	}

	// Crear archivo playlist.m3u8
	playlistPath := filepath.Join(outputDir, "playlist.m3u8")
	playlistContent := `#EXTM3U
#EXT-X-VERSION:3
#EXT-X-TARGETDURATION:10
#EXT-X-MEDIA-SEQUENCE:0
#EXTINF:10.000000,
segment0.ts
#EXTINF:10.000000,
segment1.ts
#EXT-X-ENDLIST`

	if err := os.WriteFile(playlistPath, []byte(playlistContent), 0644); err != nil {
		return fmt.Errorf("error al crear playlist: %w", err)
	}

	// Crear segmentos simulados
	for i := 0; i < 2; i++ {
		segmentPath := filepath.Join(outputDir, fmt.Sprintf("segment%d.ts", i))
		content := fmt.Sprintf("Segmento simulado #%d para %s", i, filepath.Base(inputPath))
		
		if err := os.WriteFile(segmentPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("error al crear segmento %d: %w", i, err)
		}
	}

	return nil
}