package envelope

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

// Envelope represents encrypted data with its encrypted key.
// The data is encrypted with a Data Encryption Key (DEK),
// and the DEK itself is encrypted with a Key Encryption Key (KEK).
type Envelope struct {
	// EncryptedData is the ciphertext of the original data
	EncryptedData []byte `json:"encrypted_data"`
	// EncryptedKey is the encrypted DEK
	EncryptedKey []byte `json:"encrypted_key"`
	// Nonce for the data encryption
	DataNonce []byte `json:"data_nonce"`
	// Nonce for the key encryption
	KeyNonce []byte `json:"key_nonce"`
}

// Encrypt performs envelope encryption on the provided data.
// It generates a random Data Encryption Key (DEK), encrypts the data with it,
// then encrypts the DEK with the provided Key Encryption Key (KEK).
//
// The KEK must be 32 bytes (256 bits) for AES-256.
func Encrypt(data []byte, kek []byte) (*Envelope, error) {
	if len(kek) != 32 {
		return nil, errors.New("key encryption key (kek) must be 32 bytes (AES-256)")
	}

	// Generate a random DEK (32 bytes for AES-256)
	dek := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, dek); err != nil {
		return nil, fmt.Errorf("failed to generate DEK: %w", err)
	}

	// Encrypt the data with the DEK
	encryptedData, dataNonce, err := encryptAESGCM(data, dek)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt data: %w", err)
	}

	// Encrypt the DEK with the KEK
	encryptedKey, keyNonce, err := encryptAESGCM(dek, kek)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt DEK: %w", err)
	}

	return &Envelope{
		EncryptedData: encryptedData,
		EncryptedKey:  encryptedKey,
		DataNonce:     dataNonce,
		KeyNonce:      keyNonce,
	}, nil
}

// Decrypt performs envelope decryption.
// It first decrypts the DEK using the KEK, then uses the DEK to decrypt the data.
func Decrypt(env *Envelope, kek []byte) ([]byte, error) {
	if env == nil {
		return nil, errors.New("envelope is nil")
	}
	if len(kek) != 32 {
		return nil, errors.New("key encryption key (kek) must be 32 bytes (AES-256)")
	}

	// Decrypt the DEK with the KEK
	dek, err := decryptAESGCM(env.EncryptedKey, env.KeyNonce, kek)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt DEK: %w", err)
	}

	// Decrypt the data with the DEK
	data, err := decryptAESGCM(env.EncryptedData, env.DataNonce, dek)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}

	return data, nil
}

// encryptAESGCM encrypts data using AES-GCM.
func encryptAESGCM(plaintext, key []byte) (ciphertext, nonce []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	nonce = make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}

	ciphertext = gcm.Seal(nil, nonce, plaintext, nil)
	return ciphertext, nonce, nil
}

// decryptAESGCM decrypts data using AES-GCM.
func decryptAESGCM(ciphertext, nonce, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(nonce) != gcm.NonceSize() {
		return nil, errors.New("invalid nonce size")
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// GenerateKEK generates a secure random Key Encryption Key (KEK).
// Returns a 32-byte key suitable for AES-256.
func GenerateKEK() ([]byte, error) {
	kek := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, kek); err != nil {
		return nil, fmt.Errorf("failed to generate KEK: %w", err)
	}
	return kek, nil
}
