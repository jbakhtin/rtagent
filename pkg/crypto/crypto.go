package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/go-faster/errors"
	"os"
)

var AESKeyLength = 16

func generateCryptoKey(n int) ([]byte, error) {
	data := make([]byte, n)
	_, err := rand.Read(data)
	if err != nil {
		return nil, fmt.Errorf("genCryptoKey: %w", err)
	}
	return data, nil
}

func generateAESKey() ([]byte, error) {
	return generateCryptoKey(AESKeyLength)
}

func GetEncryptedMessage(rsaPubKey *rsa.PublicKey, data []byte) ([]byte, string, error) {
	aesKey, err := generateAESKey()
	if err != nil {
		return nil, "", fmt.Errorf("AES gen failed, %w", err)
	}

	aesMessage, err := aesEncrypt(aesKey, data)
	if err != nil {
		return nil, "", fmt.Errorf("aes encryption failed, %w", err)
	}

	encryptedKey, err := rsaEncrypt(rsaPubKey, aesKey)
	if err != nil {
		return nil, "", fmt.Errorf("rsa encryption failed, %w", err)
	}

	msg := make([]byte, base64.RawStdEncoding.EncodedLen(len(aesMessage)))
	base64.RawStdEncoding.Encode(msg, aesMessage)

	key := base64.RawStdEncoding.EncodeToString(encryptedKey)
	return msg, key, nil
}

func GetDecryptedMessage(rsaPrivateKey *rsa.PrivateKey, cypher []byte, aesEncryptedKey string) ([]byte, error) {
	encryptedMsg := make([]byte, base64.RawStdEncoding.DecodedLen(len(cypher)))
	_, err := base64.RawStdEncoding.Decode(encryptedMsg, cypher)
	if err != nil {
		return nil, err
	}

	encryptedKey, err := base64.RawStdEncoding.DecodeString(aesEncryptedKey)
	if err != nil {
		return nil, err
	}

	key, err := rsaDecrypt(rsaPrivateKey, encryptedKey)
	if err != nil {
		return nil, fmt.Errorf("rsa decryption failed, %w", err)
	}

	msg, err := aesDecrypt(key, encryptedMsg)
	if err != nil {
		return nil, fmt.Errorf("aes decryption failed, %w", err)
	}

	return msg, nil
}

func rsaEncrypt(key *rsa.PublicKey, data []byte) ([]byte, error) {
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, key, data, nil)
}

func rsaDecrypt(key *rsa.PrivateKey, data []byte) ([]byte, error) {
	return rsa.DecryptOAEP(sha256.New(), rand.Reader, key, data, nil)
}

func aesEncrypt(aesKey []byte, plaintext []byte) ([]byte, error) {
	aesBlock, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	return ciphertext, nil
}

func aesDecrypt(aesKey []byte, ciphertext []byte) ([]byte, error) {
	aesBlock, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

var ErrNoRSAKey = errors.New("no key provided")

func GetPublicKey(path string) (*rsa.PublicKey, error) {
	if len(path) == 0 {
		return nil, ErrNoRSAKey
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("GetPublicKey: %w", err)
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("GetPublicKey: failed to decode PEM block containing public key")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("GetPublicKey: %w", err)
	}
	return pub.(*rsa.PublicKey), nil
}

func GetPrivateKey(path string) (*rsa.PrivateKey, error) {
	if len(path) == 0 {
		return nil, ErrNoRSAKey
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("GetPrivateKey: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil  {
		return nil, fmt.Errorf("GetPrivateKey: failed to decode PEM block containing private key")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("GetPrivateKey: %w", err)
	}

	return priv, nil
}