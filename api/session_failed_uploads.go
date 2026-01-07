package api

import (
	"context"
	"fmt"
)

type FailedUploadInvoiceStatus struct {
	Code        int      `json:"code"`
	Description string   `json:"description"`
	Details     []string `json:"details"`
}

type FailedUploadInvoice struct {
	OrdinalNumber   int                        `json:"ordinalNumber"`
	ReferenceNumber string                     `json:"referenceNumber"`
	InvoiceHash     string                     `json:"invoiceHash"`
	Status          *FailedUploadInvoiceStatus `json:"status"`
}

type FailedUploadInvoicesResponse struct {
	ContinuationToken string                `json:"continuationToken"`
	Invoices          []FailedUploadInvoice `json:"invoices"`
}

// GetFailedUploadData lists invoices that failed during upload for the session.
func GetFailedUploadData(ctx context.Context, session *UploadSession, s *Client) (*FailedUploadInvoicesResponse, error) {
	if session == nil {
		return nil, fmt.Errorf("upload session is nil")
	}

	token, err := s.AccessTokenValue(ctx)
	if err != nil {
		return nil, err
	}

	response := &FailedUploadInvoicesResponse{}
	resp, err := s.Client.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetResult(response).
		Get(s.URL + "/sessions/" + session.ReferenceNumber + "/invoices/failed")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, newErrorResponse(resp)
	}

	return response, nil
}
