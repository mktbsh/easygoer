// Package envelope provides envelope encryption and decryption functionality.
//
// Envelope encryption is a cryptographic technique where:
//   - Data is encrypted with a randomly generated Data Encryption Key (DEK)
//   - The DEK itself is encrypted with a Key Encryption Key (KEK)
//   - Both the encrypted data and the encrypted DEK are stored together
//
// This approach provides several benefits:
//   - Large amounts of data are encrypted with symmetric encryption (fast)
//   - Only the small DEK needs to be encrypted with the KEK
//   - The KEK can be rotated without re-encrypting all the data
//   - Different KEKs can be used for different data while maintaining efficiency
//
// Example usage:
//
//	// Generate a KEK (or load from secure storage)
//	kek, err := GenerateKEK()
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Encrypt some data
//	data := []byte("Secret message")
//	env, err := Encrypt(data, kek)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Later, decrypt the data
//	decrypted, err := Decrypt(env, kek)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// Security:
//   - Uses AES-256-GCM for both data and key encryption
//   - Generates cryptographically secure random keys and nonces
//   - Provides authenticated encryption (integrity and confidentiality)
//   - KEK must be 32 bytes (256 bits) for AES-256
package envelope
