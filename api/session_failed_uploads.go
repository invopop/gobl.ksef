package api

import (
	"context"
)

// FailedUploadInvoiceStatus describes the status payload for invoices that failed to upload.
type FailedUploadInvoiceStatus struct {
	Code        int      `json:"code"`
	Description string   `json:"description"`
	Details     []string `json:"details"`
}

// FailedUploadInvoice contains a single failed invoice entry returned by the API.
type FailedUploadInvoice struct {
	OrdinalNumber   int                        `json:"ordinalNumber"`
	ReferenceNumber string                     `json:"referenceNumber"`
	InvoiceHash     string                     `json:"invoiceHash"`
	Status          *FailedUploadInvoiceStatus `json:"status"`
}

type failedUploadInvoicesResponse struct {
	ContinuationToken string                `json:"continuationToken"`
	Invoices          []FailedUploadInvoice `json:"invoices"`
}

// GetFailedUploadData lists invoices that failed during upload for the session, following continuation tokens if needed.
func (s *UploadSession) GetFailedUploadData(ctx context.Context) ([]FailedUploadInvoice, error) {
	c, err := s.clientForRequests()
	if err != nil {
		return nil, err
	}

	token, err := c.getAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	var (
		allInvoices       []FailedUploadInvoice
		continuationToken string
	)

	for {
		response := &failedUploadInvoicesResponse{}

		req := c.client.R().
			SetContext(ctx).
			SetAuthToken(token).
			SetResult(response)
		if continuationToken != "" {
			req.SetHeader("x-continuation-token", continuationToken)
		}

		resp, err := req.Get(c.url + "/sessions/" + s.ReferenceNumber + "/invoices/failed")
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
