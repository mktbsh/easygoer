package envelope

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	// Generate a KEK
	kek, err := GenerateKEK()
	if err != nil {
		t.Fatalf("GenerateKEK failed: %v", err)
	}

	testData := []byte("Hello, World! This is a test message for envelope encryption.")

	// Encrypt
	env, err := Encrypt(testData, kek)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// Verify envelope structure
	if len(env.EncryptedData) == 0 {
		t.Fatal("EncryptedData is empty")
	}
	if len(env.EncryptedKey) == 0 {
		t.Fatal("EncryptedKey is empty")
	}
	if len(env.DataNonce) == 0 {
		t.Fatal("DataNonce is empty")
	}
	if len(env.KeyNonce) == 0 {
		t.Fatal("KeyNonce is empty")
	}

	// Decrypt
	decrypted, err := Decrypt(env, kek)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	// Verify
	if !bytes.Equal(testData, decrypted) {
		t.Fatalf("Decrypted data does not match original.\nExpected: %s\nGot: %s", testData, decrypted)
	}
}

func TestEncrypt_InvalidKEKSize(t *testing.T) {
	testData := []byte("test data")

	// Test with invalid KEK sizes
	invalidKEKs := [][]byte{
		make([]byte, 16), // Too short
		make([]byte, 24), // Wrong size
		make([]byte, 64), // Too long
		[]byte{},         // Empty
	}

	for _, kek := range invalidKEKs {
		_, err := Encrypt(testData, kek)
		if err == nil {
			t.Errorf("Expected error for KEK size %d, got nil", len(kek))
		}
	}
}

func TestDecrypt_InvalidKEKSize(t *testing.T) {
	env := &Envelope{
		EncryptedData: []byte("dummy"),
		EncryptedKey:  []byte("dummy"),
		DataNonce:     []byte("dummy"),
		KeyNonce:      []byte("dummy"),
	}

	// Test with invalid KEK sizes
	invalidKEKs := [][]byte{
		make([]byte, 16), // Too short
		make([]byte, 24), // Wrong size
		make([]byte, 64), // Too long
		[]byte{},         // Empty
	}

	for _, kek := range invalidKEKs {
		_, err := Decrypt(env, kek)
		if err == nil {
			t.Errorf("Expected error for KEK size %d, got nil", len(kek))
		}
	}
}

func TestDecrypt_NilEnvelope(t *testing.T) {
	kek, err := GenerateKEK()
	if err != nil {
		t.Fatalf("GenerateKEK failed: %v", err)
	}

	_, err = Decrypt(nil, kek)
	if err == nil {
		t.Fatal("Expected error for nil envelope, got nil")
	}
}

func TestDecrypt_WrongKEK(t *testing.T) {
	kek1, err := GenerateKEK()
	if err != nil {
		t.Fatalf("GenerateKEK failed: %v", err)
	}

	kek2, err := GenerateKEK()
	if err != nil {
		t.Fatalf("GenerateKEK failed: %v", err)
	}

	testData := []byte("Secret message")

	// Encrypt with kek1
	env, err := Encrypt(testData, kek1)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// Try to decrypt with kek2 (wrong key)
	_, err = Decrypt(env, kek2)
	if err == nil {
		t.Fatal("Expected error when decrypting with wrong KEK, got nil")
	}
}

func TestGenerateKEK(t *testing.T) {
	kek, err := GenerateKEK()
	if err != nil {
		t.Fatalf("GenerateKEK failed: %v", err)
	}

	if len(kek) != 32 {
		t.Fatalf("Expected KEK length 32, got %d", len(kek))
	}

	// Generate another KEK and ensure they're different (extremely unlikely to be equal)
	kek2, err := GenerateKEK()
	if err != nil {
		t.Fatalf("GenerateKEK failed: %v", err)
	}

	if bytes.Equal(kek, kek2) {
		t.Fatal("Two generated KEKs are identical (should be virtually impossible)")
	}
}

func TestEnvelope_JSON_Serialization(t *testing.T) {
	kek, err := GenerateKEK()
	if err != nil {
		t.Fatalf("GenerateKEK failed: %v", err)
	}

	testData := []byte("JSON serialization test")

	// Encrypt
	env, err := Encrypt(testData, kek)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(env)
	if err != nil {
		t.Fatalf("JSON Marshal failed: %v", err)
	}

	// Unmarshal from JSON
	var env2 Envelope
	if err := json.Unmarshal(jsonData, &env2); err != nil {
		t.Fatalf("JSON Unmarshal failed: %v", err)
	}

	// Verify fields match
	if !bytes.Equal(env.EncryptedData, env2.EncryptedData) {
		t.Error("EncryptedData mismatch after JSON round-trip")
	}
	if !bytes.Equal(env.EncryptedKey, env2.EncryptedKey) {
		t.Error("EncryptedKey mismatch after JSON round-trip")
	}
	if !bytes.Equal(env.DataNonce, env2.DataNonce) {
		t.Error("DataNonce mismatch after JSON round-trip")
	}
	if !bytes.Equal(env.KeyNonce, env2.KeyNonce) {
		t.Error("KeyNonce mismatch after JSON round-trip")
	}

	// Decrypt the unmarshaled envelope
	decrypted, err := Decrypt(&env2, kek)
	if err != nil {
		t.Fatalf("Decrypt after JSON round-trip failed: %v", err)
	}

	if !bytes.Equal(testData, decrypted) {
		t.Fatalf("Data mismatch after JSON round-trip and decryption")
	}
}

func TestEncryptDecrypt_EmptyData(t *testing.T) {
	kek, err := GenerateKEK()
	if err != nil {
		t.Fatalf("GenerateKEK failed: %v", err)
	}

	testData := []byte{}

	// Encrypt empty data
	env, err := Encrypt(testData, kek)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// Decrypt
	decrypted, err := Decrypt(env, kek)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	// Verify
	if !bytes.Equal(testData, decrypted) {
		t.Fatalf("Decrypted data does not match original empty data")
	}
}

func TestEncryptDecrypt_LargeData(t *testing.T) {
	kek, err := GenerateKEK()
	if err != nil {
		t.Fatalf("GenerateKEK failed: %v", err)
	}

	// Create 1MB of test data
	testData := make([]byte, 1024*1024)
	for i := range testData {
		testData[i] = byte(i % 256)
	}

	// Encrypt
	env, err := Encrypt(testData, kek)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// Decrypt
	decrypted, err := Decrypt(env, kek)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	// Verify
	if !bytes.Equal(testData, decrypted) {
		t.Fatalf("Decrypted large data does not match original")
	}
}

func TestDecrypt_CorruptedData(t *testing.T) {
	kek, err := GenerateKEK()
	if err != nil {
		t.Fatalf("GenerateKEK failed: %v", err)
	}

	testData := []byte("Test data for corruption")

	// Encrypt
	env, err := Encrypt(testData, kek)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// Corrupt the encrypted data
	if len(env.EncryptedData) > 0 {
		env.EncryptedData[0] ^= 0xFF
	}

	// Try to decrypt corrupted data
	_, err = Decrypt(env, kek)
	if err == nil {
		t.Fatal("Expected error when decrypting corrupted data, got nil")
	}
}

func TestDecrypt_CorruptedKey(t *testing.T) {
	kek, err := GenerateKEK()
	if err != nil {
		t.Fatalf("GenerateKEK failed: %v", err)
	}

	testData := []byte("Test data for key corruption")

	// Encrypt
	env, err := Encrypt(testData, kek)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// Corrupt the encrypted key
	if len(env.EncryptedKey) > 0 {
		env.EncryptedKey[0] ^= 0xFF
	}

	// Try to decrypt with corrupted key
	_, err = Decrypt(env, kek)
	if err == nil {
		t.Fatal("Expected error when decrypting with corrupted key, got nil")
	}
}
