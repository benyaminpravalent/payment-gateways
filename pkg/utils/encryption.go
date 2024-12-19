package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"payment-gateway/models"
)

func EncryptAES(plaintext, key string) (string, error) {
	// Convert the key to a 32-byte array (AES-256 requires a 32-byte key)
	if len(key) != 32 {
		return "", errors.New("encryption key must be 32 bytes")
	}
	keyBytes := []byte(key)

	// Convert the plaintext to bytes and calculate padding
	plaintextBytes := []byte(plaintext)
	padding := aes.BlockSize - len(plaintextBytes)%aes.BlockSize
	paddedText := append(plaintextBytes, bytes.Repeat([]byte{byte(padding)}, padding)...) // PKCS#7 padding

	// Create the AES cipher block
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", errors.New("failed to create AES cipher block")
	}

	// Generate a random initialization vector (IV)
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", errors.New("failed to generate IV")
	}

	// Encrypt the plaintext using CBC mode
	ciphertext := make([]byte, len(paddedText))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, paddedText)

	// Prepend the IV to the ciphertext
	finalCiphertext := append(iv, ciphertext...)

	// Encode the ciphertext in Base64
	return base64.StdEncoding.EncodeToString(finalCiphertext), nil
}

// DecryptAES decrypts a Base64-encoded, AES-encrypted string
func DecryptAES(encryptedValue, key string) (string, error) {
	// Decode the Base64-encoded value
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedValue)
	if err != nil {
		return "", errors.New("failed to decode Base64 encrypted value")
	}

	// Convert the key to a 32-byte array (AES-256 requires a 32-byte key)
	if len(key) != 32 {
		return "", errors.New("encryption key must be 32 bytes")
	}
	keyBytes := []byte(key)

	// Create the AES cipher block
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", errors.New("failed to create AES cipher block")
	}

	// Verify the ciphertext length
	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext is too short")
	}

	// Separate the initialization vector (IV) from the ciphertext
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// Decrypt the ciphertext using CBC mode
	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	// Remove padding (PKCS#7)
	plaintext, err = removePadding(plaintext)
	if err != nil {
		return "", errors.New("failed to remove padding")
	}

	return string(plaintext), nil
}

// removePadding removes PKCS#7 padding from the decrypted plaintext
func removePadding(plaintext []byte) ([]byte, error) {
	paddingLen := int(plaintext[len(plaintext)-1])
	if paddingLen > len(plaintext) || paddingLen == 0 {
		return nil, errors.New("invalid padding length")
	}
	for _, padByte := range plaintext[len(plaintext)-paddingLen:] {
		if padByte != byte(paddingLen) {
			return nil, errors.New("invalid padding")
		}
	}
	return plaintext[:len(plaintext)-paddingLen], nil
}

func EncryptTransactionRequest(request *models.TransactionRequest, gatewayPrivateKey string) (string, error) {
	// Serialize the request struct to JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("failed to serialize request to JSON: %w", err)
	}

	// Encrypt the JSON data using the provided gateway private key
	encryptedData, err := EncryptAES(string(jsonData), gatewayPrivateKey)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt request: %w", err)
	}

	return encryptedData, nil
}
