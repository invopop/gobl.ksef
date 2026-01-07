package api

import (
	"context"
	"fmt"
	"time"
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

// Note that there are more fields, but we only need these for now
type SessionStatus struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

type SessionStatusResponse struct {
	Status *SessionStatus `json:"status"`
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

// TerminateSession ends the current session. When the session is terminated, all uploaded invoices start
// to be processed by the KSeF system.
func TerminateSession(session *UploadSession, ctx context.Context, s *Client) error {
	token, err := s.AccessTokenValue(ctx)
	if err != nil {
		return err
	}

	resp, err := s.Client.R().
		SetContext(ctx).
		SetAuthToken(token).
		Post(s.URL + "/sessions/online/" + session.ReferenceNumber + "/close")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return newErrorResponse(resp)
	}

	return nil
}

// PollSessionStatus checks the status of a session after upload is completed.
func PollSessionStatus(ctx context.Context, session *UploadSession, s *Client) (*SessionStatusResponse, error) {
	if session == nil {
		return nil, fmt.Errorf("upload session is nil")
	}

	token, err := s.AccessTokenValue(ctx)
	if err != nil {
		return nil, err
	}

	attempt := 0
	for {
		attempt++
		if attempt > 30 {
			return nil, fmt.Errorf("session polling count exceeded")
		}

		response := &SessionStatusResponse{}
		resp, err := s.Client.R().
			SetContext(ctx).
			SetAuthToken(token).
			SetResult(response).
			Get(s.URL + "/sessions/" + session.ReferenceNumber)
		if err != nil {
			return nil, err
		}
		if resp.IsError() {
			return nil, newErrorResponse(resp)
		}

		if response.Status == nil {
			return nil, fmt.Errorf("session status missing in response")
		}

		switch response.Status.Code {
		case 100, 150, 170: // still processing
			time.Sleep(2 * time.Second)
			continue
		case 200:
			return response, nil
		default:
			return nil, fmt.Errorf("session failed: %s", response.Status.Description)
		}
	}
}
