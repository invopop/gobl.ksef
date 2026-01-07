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
	client               *Client
}

func (s *UploadSession) clientForRequests() (*Client, error) {
	if s == nil {
		return nil, fmt.Errorf("upload session is nil")
	}
	if s.client == nil {
		return nil, fmt.Errorf("upload session missing client")
	}
	return s.client, nil
}

// Note that there are more fields, but we only need these for now
type SessionStatus struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

// UPO = urzędowe potwierdzenie odbioru = confirmation that the invoice has been successfully received by the system
type SessionStatusUpoPage struct {
	ReferenceNumber           string `json:"referenceNumber"`
	DownloadURL               string `json:"downloadUrl"`
	DownloadURLExpirationDate string `json:"downloadUrlExpirationDate"`
}

type SessionStatusUpo struct {
	Pages []SessionStatusUpoPage `json:"pages"`
}

type SessionStatusResponse struct {
	Status                 *SessionStatus    `json:"status"`
	InvoiceCount           int               `json:"invoiceCount"`
	SuccessfulInvoiceCount int               `json:"successfulInvoiceCount"`
	FailedInvoiceCount     int               `json:"failedInvoiceCount"`
	Upo                    *SessionStatusUpo `json:"upo"`
}

// CreateSession opens a new upload session in online (interactive) mode, allowing to upload invoices one by one
// (There exists also a batch mode, where a ZIP file can be uploaded)
func (c *Client) CreateSession(ctx context.Context) (*UploadSession, error) {
	token, err := c.AccessTokenValue(ctx)
	if err != nil {
		return nil, err
	}
	publicKeyCertificate, err := GetRSAPublicKey(ctx, c)
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

	resp, err := c.Client.R().
		SetBody(request).
		SetResult(response).
		SetContext(ctx).
		SetAuthToken(token).
		Post(c.URL + "/sessions/online")

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
		client:               c,
	}, nil
}

// FinishUpload ends the current session. When the session is terminated, all uploaded invoices start
// to be processed by the KSeF system.
func (s *UploadSession) FinishUpload(ctx context.Context) error {
	c, err := s.clientForRequests()
	if err != nil {
		return err
	}

	token, err := c.AccessTokenValue(ctx)
	if err != nil {
		return err
	}

	resp, err := c.Client.R().
		SetContext(ctx).
		SetAuthToken(token).
		Post(c.URL + "/sessions/online/" + s.ReferenceNumber + "/close")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return newErrorResponse(resp)
	}

	return nil
}

// PollSessionStatus checks the status of an upload session, after upload is completed.
func (s *UploadSession) PollSessionStatus(ctx context.Context) (*SessionStatusResponse, error) {
	c, err := s.clientForRequests()
	if err != nil {
		return nil, err
	}

	token, err := c.AccessTokenValue(ctx)
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
		resp, err := c.Client.R().
			SetContext(ctx).
			SetAuthToken(token).
			SetResult(response).
			Get(c.URL + "/sessions/" + s.ReferenceNumber)
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
