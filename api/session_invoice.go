package api

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

type uploadInvoiceRequest struct {
	InvoiceHash             string `json:"invoiceHash"`
	InvoiceSize             int    `json:"invoiceSize"`
	EncryptedInvoiceHash    string `json:"encryptedInvoiceHash"`
	EncryptedInvoiceSize    int    `json:"encryptedInvoiceSize"`
	EncryptedInvoiceContent string `json:"encryptedInvoiceContent"`
	OfflineMode             bool   `json:"offlineMode"`
}

func buildUploadInvoiceRequest(session *UploadSession, invoice []byte) (*uploadInvoiceRequest, error) {
	if len(invoice) == 0 {
		return nil, fmt.Errorf("invoice payload is empty")
	}
	if len(session.SymmetricKey) != 32 {
		return nil, fmt.Errorf("symmetric key must be 32 bytes, got %d", len(session.SymmetricKey))
	}
	if len(session.InitializationVector) != aes.BlockSize {
		return nil, fmt.Errorf("initialization vector must be %d bytes, got %d", aes.BlockSize, len(session.InitializationVector))
	}

	invoiceHash := sha256.Sum256(invoice)
	encryptedInvoice, err := encryptInvoice(session.SymmetricKey, session.InitializationVector, invoice)
	if err != nil {
		return nil, err
	}
	encryptedInvoiceHash := sha256.Sum256(encryptedInvoice)

	return &uploadInvoiceRequest{
		InvoiceHash:             base64.StdEncoding.EncodeToString(invoiceHash[:]),
		InvoiceSize:             len(invoice),
		EncryptedInvoiceHash:    base64.StdEncoding.EncodeToString(encryptedInvoiceHash[:]),
		EncryptedInvoiceSize:    len(encryptedInvoice),
		EncryptedInvoiceContent: base64.StdEncoding.EncodeToString(encryptedInvoice),
		OfflineMode:             false,
	}, nil
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
