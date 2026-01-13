package api

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"fmt"
)

// buildSessionEncryption generates AES key/IV pair and encrypts the key using the RSA certificate.
// It returns both the payload required by the API and the raw values so they can encrypt invoice data later.
func buildSessionEncryption(certificate string) (*createSessionEncryption, []byte, []byte, error) {
	publicKey, err := certificateToPublicKey(certificate)
	if err != nil {
		return nil, nil, nil, err
	}

	symmetricKey, err := randomBytes(32)
	if err != nil {
		return nil, nil, nil, err
	}

	initializationVector, err := randomBytes(16)
	if err != nil {
		return nil, nil, nil, err
	}

	encryptedSymmetricKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, symmetricKey, nil)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("encrypt symmetric key: %w", err)
	}

	return &createSessionEncryption{
		EncryptedSymmetricKey: base64.StdEncoding.EncodeToString(encryptedSymmetricKey),
		InitializationVector:  base64.StdEncoding.EncodeToString(initializationVector),
	}, symmetricKey, initializationVector, nil
}

func certificateToPublicKey(encoded string) (*rsa.PublicKey, error) {
	certBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("decode RSA certificate: %w", err)
	}

	certificate, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return nil, fmt.Errorf("parse RSA certificate: %w", err)
	}

	publicKey, ok := certificate.PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("certificate public key is not RSA")
	}

	return publicKey, nil
}

func randomBytes(length int) ([]byte, error) {
	buffer := make([]byte, length)
	if _, err := rand.Read(buffer); err != nil {
		return nil, fmt.Errorf("read random bytes: %w", err)
	}
	return buffer, nil
}

func encryptInvoice(key, iv, invoice []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create cipher: %w", err)
	}

	padded := pkcs7Pad(invoice, aes.BlockSize)
	ciphertext := make([]byte, len(padded))
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(ciphertext, padded)

	return ciphertext, nil
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	if padding == 0 {
		padding = blockSize
	}

	out := make([]byte, len(data)+padding)
	copy(out, data)
	for i := len(data); i < len(out); i++ {
		out[i] = byte(padding)
	}
	return out
}
