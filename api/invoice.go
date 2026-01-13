package api

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// UploadInvoice uploads a serialized invoice in the provided upload session.
func (s *UploadSession) UploadInvoice(ctx context.Context, invoice []byte) error {
	if s == nil {
		return fmt.Errorf("upload session is nil")
	}
	if s.ReferenceNumber == "" {
		return fmt.Errorf("upload session missing reference number")
	}

	c, err := s.clientForRequests()
	if err != nil {
		return err
	}

	token, err := c.getAccessToken(ctx)
	if err != nil {
		return err
	}

	request, err := s.buildUploadInvoiceRequest(invoice)
	if err != nil {
		return err
	}

	resp, err := c.client.R().
		SetBody(request).
		SetContext(ctx).
		SetAuthToken(token).
		Post(c.url + "/sessions/online/" + s.ReferenceNumber + "/invoices")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return newErrorResponse(resp)
	}

	return nil
}

type uploadInvoiceRequest struct {
	InvoiceHash             string `json:"invoiceHash"`
	InvoiceSize             int    `json:"invoiceSize"`
	EncryptedInvoiceHash    string `json:"encryptedInvoiceHash"`
	EncryptedInvoiceSize    int    `json:"encryptedInvoiceSize"`
	EncryptedInvoiceContent string `json:"encryptedInvoiceContent"`
	OfflineMode             bool   `json:"offlineMode"`
}

func (s *UploadSession) buildUploadInvoiceRequest(invoice []byte) (*uploadInvoiceRequest, error) {
	if len(invoice) == 0 {
		return nil, fmt.Errorf("invoice payload is empty")
	}
	if len(s.SymmetricKey) != 32 {
		return nil, fmt.Errorf("symmetric key must be 32 bytes, got %d", len(s.SymmetricKey))
	}
	if len(s.InitializationVector) != aes.BlockSize {
		return nil, fmt.Errorf("initialization vector must be %d bytes, got %d", aes.BlockSize, len(s.InitializationVector))
	}

	invoiceHash := sha256.Sum256(invoice)
	encryptedInvoice, err := encryptInvoice(s.SymmetricKey, s.InitializationVector, invoice)
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
