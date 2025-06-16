package content

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// MemoryVideoStore implementa VideoStore en memoria
type MemoryVideoStore struct {
	mu     sync.RWMutex
	videos map[string]*VideoMetadata
}

func NewMemoryVideoStore() *MemoryVideoStore {
	return &MemoryVideoStore{
		videos: make(map[string]*VideoMetadata),
	}
}

func (s *MemoryVideoStore) SaveMetadata(metadata *VideoMetadata) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	metadata.CreatedAt = time.Now().Format(time.RFC3339)
	s.videos[metadata.ID] = metadata
	return nil
}

func (s *MemoryVideoStore) GetVideoMetadata(id string) (*VideoMetadata, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	video, exists := s.videos[id]
	if !exists {
		return nil, errors.New("video no encontrado")
	}
	return video, nil
}

func (s *MemoryVideoStore) ListVideos(userID string) ([]*VideoMetadata, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	result := []*VideoMetadata{}
	for _, video := range s.videos {
		if video.OwnerID == userID {
			result = append(result, video)
		}
	}
	return result, nil
}

// JSONVideoStore implementa VideoStore con archivos JSON
type JSONVideoStore struct {
	filePath string
	mu       sync.RWMutex
}

func NewJSONVideoStore(filePath string) *JSONVideoStore {
	return &JSONVideoStore{filePath: filePath}
}

func (s *JSONVideoStore) load() (map[string]*VideoMetadata, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	data := make(map[string]*VideoMetadata)
	
	file, err := os.Open(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return data, nil
		}
		return nil, err
	}
	defer file.Close()
	
	err = json.NewDecoder(file).Decode(&data)
	return data, err
}

func (s *JSONVideoStore) save(data map[string]*VideoMetadata) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	// Crear directorio si no existe
	dir := filepath.Dir(s.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	
	file, err := os.Create(s.filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	
	return json.NewEncoder(file).Encode(data)
}

func (s *JSONVideoStore) SaveMetadata(metadata *VideoMetadata) error {
	data, err := s.load()
	if err != nil {
		return err
	}
	
	metadata.CreatedAt = time.Now().Format(time.RFC3339)
	data[metadata.ID] = metadata
	return s.save(data)
}

func (s *JSONVideoStore) GetVideoMetadata(id string) (*VideoMetadata, error) {
	data, err := s.load()
	if err != nil {
		return nil, err
	}
	
	video, exists := data[id]
	if !exists {
		return nil, errors.New("video no encontrado")
	}
	return video, nil
}

func (s *JSONVideoStore) ListVideos(userID string) ([]*VideoMetadata, error) {
	data, err := s.load()
	if err != nil {
		return nil, err
	}
	
	result := []*VideoMetadata{}
	for _, video := range data {
		if video.OwnerID == userID {
			result = append(result, video)
		}
	}
	return result, nil
}