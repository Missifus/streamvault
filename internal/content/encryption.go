package content

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

// AESEncryptionService implementa EncryptionService con AES-256-CTR
type AESEncryptionService struct {
	key []byte
}

func NewAESEncryptionService(key string) (*AESEncryptionService, error) {
	if len(key) != 32 {
		return nil, errors.New("la clave debe tener 32 bytes")
	}
	return &AESEncryptionService{key: []byte(key)}, nil
}

func (s *AESEncryptionService) EncryptFile(inputPath, outputPath string) error {
	// Implementación completa en el módulo anterior
	// ...
}

func (s *AESEncryptionService) DecryptStream(r io.Reader, w io.Writer) error {
	// Implementación completa en el módulo anterior
	// ...
}