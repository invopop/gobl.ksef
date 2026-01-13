package api

import (
	"context"
	"fmt"
	"time"
)

// CreateSessionFormCode identifies the document schema that will be submitted during a session.
type CreateSessionFormCode struct {
	SystemCode    string `json:"systemCode"`
	SchemaVersion string `json:"schemaVersion"`
	Value         string `json:"value"`
}

type createSessionEncryption struct {
	EncryptedSymmetricKey string `json:"encryptedSymmetricKey"`
	InitializationVector  string `json:"initializationVector"`
}

type createSessionRequest struct {
	FormCode   CreateSessionFormCode   `json:"formCode"`
	Encryption createSessionEncryption `json:"encryption"`
}

type createSessionResponse struct {
	ReferenceNumber string `json:"referenceNumber"`
	ValidUntil      string `json:"validUntil"`
}

// UploadSession represents a live KSeF invoice upload session, including encryption material and metadata.
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

// SessionStatus contains basic status information for a session returned by the API.
type SessionStatus struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

// SessionStatusUpoPage stores a single confirmation (UPO) download page returned by the service.
// UPO = urzędowe potwierdzenie odbioru = confirmation that the invoice has been successfully received by the system.
type SessionStatusUpoPage struct {
	ReferenceNumber           string `json:"referenceNumber"`
	DownloadURL               string `json:"downloadUrl"`
	DownloadURLExpirationDate string `json:"downloadUrlExpirationDate"`
}

// SessionStatusUpo groups all UPO pages associated with a session.
type SessionStatusUpo struct {
	Pages []SessionStatusUpoPage `json:"pages"`
}

// SessionStatusResponse summarizes the result of polling a session, including invoice stats and UPO links.
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
	token, err := c.getAccessToken(ctx)
	if err != nil {
		return nil, err
	}
	publicKeyCertificate, err := getRSAPublicKey(ctx, c)
	if err != nil {
		return nil, err
	}

	encryption, symmetricKey, initializationVector, err := buildSessionEncryption(publicKeyCertificate.Certificate)
	if err != nil {
		return nil, err
	}

	request := &createSessionRequest{
		FormCode: CreateSessionFormCode{
			SystemCode:    "FA (3)",
			SchemaVersion: "1-0E",
			Value:         "FA",
		},
		Encryption: *encryption,
	}
	response := &createSessionResponse{}

	resp, err := c.client.R().
		SetBody(request).
		SetResult(response).
		SetContext(ctx).
		SetAuthToken(token).
		Post(c.url + "/sessions/online")

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

	token, err := c.getAccessToken(ctx)
	if err != nil {
		return err
	}

	resp, err := c.client.R().
		SetContext(ctx).
		SetAuthToken(token).
		Post(c.url + "/sessions/online/" + s.ReferenceNumber + "/close")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return newErrorResponse(resp)
	}

	return nil
}

// PollStatus checks the status of an upload session, after upload is completed.
func (s *UploadSession) PollStatus(ctx context.Context) (*SessionStatusResponse, error) {
	c, err := s.clientForRequests()
	if err != nil {
		return nil, err
	}

	token, err := c.getAccessToken(ctx)
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
		resp, err := c.client.R().
			SetContext(ctx).
			SetAuthToken(token).
			SetResult(response).
			Get(c.url + "/sessions/" + s.ReferenceNumber)
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
