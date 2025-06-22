package content

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"os"
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
	// Leer archivo de entrada
	input, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}

	// Crear cifrador AES
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return err
	}

	// Generar IV (Initialization Vector)
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return err
	}

	// Crear cifrador en modo CTR
	stream := cipher.NewCTR(block, iv)

	// Cifrar datos
	ciphertext := make([]byte, len(input))
	stream.XORKeyStream(ciphertext, input)

	// Combinar IV + datos cifrados
	output := append(iv, ciphertext...)

	// Escribir archivo cifrado
	return os.WriteFile(outputPath, output, 0644)
}

func (s *AESEncryptionService) DecryptStream(r io.Reader, w io.Writer) error {
	// Leer IV (primeros 16 bytes)
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(r, iv); err != nil {
		return err
	}

	// Crear cifrador AES
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return err
	}

	// Crear descifrador en modo CTR
	stream := cipher.NewCTR(block, iv)

	// Buffer para descifrar en streaming
	buf := make([]byte, 4096)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			// Descifrar el chunk
			stream.XORKeyStream(buf[:n], buf[:n])
			
			// Escribir datos descifrados
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