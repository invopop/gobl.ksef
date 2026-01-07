package api

import (
	"context"
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
func (s *UploadSession) GetFailedUploadData(ctx context.Context) ([]FailedUploadInvoice, error) {
	c, err := s.clientForRequests()
	if err != nil {
		return nil, err
	}

	token, err := c.AccessTokenValue(ctx)
	if err != nil {
		return nil, err
	}

	var (
		allInvoices       []FailedUploadInvoice
		continuationToken string
	)

	for {
		response := &FailedUploadInvoicesResponse{}

		req := c.Client.R().
			SetContext(ctx).
			SetAuthToken(token).
			SetResult(response)
		if continuationToken != "" {
			req.SetHeader("x-continuation-token", continuationToken)
		}

		resp, err := req.Get(c.URL + "/sessions/" + s.ReferenceNumber + "/invoices/failed")
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
