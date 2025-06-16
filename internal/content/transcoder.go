package content

import (
	"context"
	"errors"
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
	// Simulación: crear archivos HLS de ejemplo
	// ...
	return nil
}