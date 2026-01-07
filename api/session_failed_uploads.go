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

// GetFailedUploadData lists invoices that failed during upload for the session, following continuation tokens if needed.
func GetFailedUploadData(ctx context.Context, session *UploadSession, s *Client) ([]FailedUploadInvoice, error) {
	if session == nil {
		return nil, fmt.Errorf("upload session is nil")
	}

	token, err := s.AccessTokenValue(ctx)
	if err != nil {
		return nil, err
	}

	var (
		allInvoices       []FailedUploadInvoice
		continuationToken string
	)

	for {
		response := &FailedUploadInvoicesResponse{}

		req := s.Client.R().
			SetContext(ctx).
			SetAuthToken(token).
			SetResult(response)
		if continuationToken != "" {
			req.SetHeader("x-continuation-token", continuationToken)
		}

		resp, err := req.Get(s.URL + "/sessions/" + session.ReferenceNumber + "/invoices/failed")
		if err != nil {
			return nil, err
		}
		if resp.IsError() {
			return nil, newErrorResponse(resp)
		}

		allInvoices = append(allInvoices, response.Invoices...)

		if response.ContinuationToken == "" {
			break
		}
		continuationToken = response.ContinuationToken
	}

	return allInvoices, nil
}
