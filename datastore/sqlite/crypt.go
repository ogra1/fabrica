package sqlite

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"log"
)

func generateSecret() (string, error) {
	rb := make([]byte, keyLength)
	_, err := rand.Read(rb)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(rb), nil
}

// encryptKey uses symmetric encryption to encrypt the data for storage
func encryptKey(plainTextKey, keyText string) (string, error) {
	// The AES key needs to be 16 or 32 bytes i.e. AES-128 or AES-256
	aesKey := padRight(keyText, "x", 32)

	block, err := aes.NewCipher([]byte(aesKey))
	if err != nil {
		log.Printf("Error creating the cipher block: %v", err)
		return "", err
	}

	// The IV needs to be unique, but not secure. Including it at the start of the plaintext
	ciphertext := make([]byte, aes.BlockSize+len(plainTextKey))
	iv := ciphertext[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		log.Printf("Error creating the IV for the cipher: %v", err)
		return "", err
	}

	// Use CFB mode for the encryption
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plainTextKey))

	return string(ciphertext), nil
}

// decryptKey handles the decryption of a sealed signing key
func decryptKey(sealedKey []byte, keyText string) (string, error) {
	aesKey := padRight(keyText, "x", 32)

	block, err := aes.NewCipher([]byte(aesKey))
	if err != nil {
		log.Printf("Error creating the cipher block: %v", err)
		return "", err
	}

	if len(sealedKey) < aes.BlockSize {
		return "", errors.New("cipher text too short")
	}

	iv := sealedKey[:aes.BlockSize]
	sealedKey = sealedKey[aes.BlockSize:]

	// Use CFB mode for the decryption
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(sealedKey, sealedKey)

	return string(sealedKey), nil
}

// padRight truncates a string to a specific length, padding with a named
// character for shorter strings.
func padRight(str, pad string, length int) string {
	for {
		str += pad
		if len(str) > length {
			return str[0:length]
		}
	}
}
