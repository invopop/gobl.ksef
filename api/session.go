package api

import (
	"context"
	"encoding/xml"
)

type CreateSessionFormCode struct {
	SystemCode    string `json:"systemCode"`
	SchemaVersion string `json:"schemaVersion"`
	Value         string `json:"value"`
}

type CreateSessionEncryption struct {
	EncryptedSymmetricKey string `json:"encryptedSymmetricKey"`
	InitializationVector  string `json:"initializationVector"`
}

type CreateSessionRequest struct {
	FormCode   CreateSessionFormCode   `json:"formCode"`
	Encryption CreateSessionEncryption `json:"encryption"`
}

type CreateSessionResponse struct {
	ReferenceNumber string `json:"referenceNumber"`
	ValidUntil      string `json:"validUntil"`
}

type UploadSession struct {
	ReferenceNumber      string
	ValidUntil           string
	SymmetricKey         []byte
	InitializationVector []byte
}

// CreateSession opens a new upload session in online (interactive) mode, allowing to upload invoices one by one
// (There exists also a batch mode, where a ZIP file can be uploaded)
func CreateSession(ctx context.Context, s *Client) (*UploadSession, error) {
	token, err := s.AccessTokenValue(ctx)
	if err != nil {
		return nil, err
	}
	publicKeyCertificate, err := GetRSAPublicKey(ctx, s)
	if err != nil {
		return nil, err
	}

	encryption, symmetricKey, initializationVector, err := buildSessionEncryption(publicKeyCertificate.Certificate)
	if err != nil {
		return nil, err
	}

	request := &CreateSessionRequest{
		FormCode: CreateSessionFormCode{
			SystemCode:    "FA (3)",
			SchemaVersion: "1-0E",
			Value:         "FA",
		},
		Encryption: *encryption,
	}
	response := &CreateSessionResponse{}

	resp, err := s.Client.R().
		SetBody(request).
		SetResult(response).
		SetContext(ctx).
		SetAuthToken(token).
		Post(s.URL + "/sessions/online")

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, newErrorResponse(resp)
	}

	return &UploadSession{
		ReferenceNumber:      response.ReferenceNumber,
		ValidUntil:           response.ValidUntil,
		SymmetricKey:         symmetricKey,
		InitializationVector: initializationVector,
	}, nil
}

// UploadInvoice uploads a new invoice.
func UploadInvoice(ctx context.Context, s *Client) {
	// TODO complete
}

// TerminateSession ends the current session. When the session is terminated, all uploaded invoices start
// to be processed by the KSeF system.
func TerminateSession(referenceNumber string, ctx context.Context, s *Client) error {
	token, err := s.AccessTokenValue(ctx)
	if err != nil {
		return err
	}

	resp, err := s.Client.R().
		SetContext(ctx).
		SetAuthToken(token).
		Post(s.URL + "/sessions/online/" + referenceNumber + "/close")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return newErrorResponse(resp)
	}

	return nil
}

func bytes(d InitSessionTokenRequest) ([]byte, error) {
	bytes, err := xml.MarshalIndent(d, "", "  ")
	if err != nil {
		return nil, err
	}
	return append([]byte(`<?xml version="1.0" encoding="utf-8" standalone="yes"?>`+"\n"), bytes...), nil
}
