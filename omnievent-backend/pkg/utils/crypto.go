package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strings"

	"golang.org/x/crypto/pbkdf2"
)

var ErrCiphertextInvalid = errors.New("ciphertext is invalid")

func SHA256(s string) string {
	h := sha256.Sum256([]byte(s))
	return base64.StdEncoding.EncodeToString(h[:])
}

func SHA256Bytes(data []byte) []byte {
	h := sha256.Sum256(data)
	return h[:]
}

func MD5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func MD5Encode(data []byte) []byte {
	m := md5.New()
	m.Write(data)
	return m.Sum(nil)
}

func MD5EncodeToString(data []byte) string {
	hash := MD5Encode(data)
	return hex.EncodeToString(hash)
}

func GenerateRandomBytes(n int) ([]byte, error) {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func GenerateRandomString(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)[:length]
}

func GenerateSalt(length int) string {
	return GenerateRandomString(length)
}

func AESGCMEncrypt(key []byte, plainText []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)

	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)

	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesgcm.NonceSize())

	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := aesgcm.Seal(nil, nonce, plainText, nil)
	result := append(nonce, ciphertext...)

	return result, nil
}

func AESGCMDecrypt(key []byte, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)

	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)

	if err != nil {
		return nil, err
	}

	nonceSize := aesgcm.NonceSize()

	if len(ciphertext)-nonceSize <= 0 {
		return nil, ErrCiphertextInvalid
	}

	nonce := ciphertext[:nonceSize]
	ciphertext = ciphertext[nonceSize:]

	plainText, err := aesgcm.Open(nil, nonce, ciphertext, nil)

	if err != nil {
		return nil, err
	}

	return plainText, nil
}

func EncodePassword(password string, salt string) string {
	encodedPassword := pbkdf2.Key([]byte(password), []byte(salt), 10000, 48, sha256.New)
	return strings.TrimRight(base64.StdEncoding.EncodeToString(encodedPassword), "=")
}

func EncryptSecret(secret string, key string) (string, error) {
	encryptedSecret, err := AESGCMEncrypt(MD5Encode([]byte(key)), []byte(secret))

	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(encryptedSecret), nil
}

func DecryptSecret(encyptedSecret string, key string) (string, error) {
	encyptedData, err := base64.StdEncoding.DecodeString(encyptedSecret)

	if err != nil {
		return "", err
	}

	secret, err := AESGCMDecrypt(MD5Encode([]byte(key)), []byte(encyptedData))

	if err != nil {
		return "", err
	}

	return string(secret), nil
}
