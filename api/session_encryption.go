package api

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"fmt"
)

// buildSessionEncryption generates AES key/IV pair and encrypts the key using the RSA certificate.
// This is necessary when uploading and downloading invoices to/from the KSeF API.
func buildSessionEncryption(certificate string) (*CreateSessionEncryption, error) {
	publicKey, err := certificateToPublicKey(certificate)
	if err != nil {
		return nil, err
	}

	symmetricKey, err := randomBytes(32)
	if err != nil {
		return nil, err
	}

	initializationVector, err := randomBytes(16)
	if err != nil {
		return nil, err
	}

	encryptedSymmetricKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, symmetricKey, nil)
	if err != nil {
		return nil, fmt.Errorf("encrypt symmetric key: %w", err)
	}

	return &CreateSessionEncryption{
		EncryptedSymmetricKey: base64.StdEncoding.EncodeToString(encryptedSymmetricKey),
		InitializationVector:  base64.StdEncoding.EncodeToString(initializationVector),
	}, nil
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
