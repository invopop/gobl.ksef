package api

import (
	"context"
	"fmt"
)

// UploadInvoice uploads a serialized invoice using the provided upload session data.
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

	token, err := c.AccessTokenValue(ctx)
	if err != nil {
		return err
	}

	request, err := buildUploadInvoiceRequest(s, invoice)
	if err != nil {
		return err
	}

	resp, err := c.Client.R().
		SetBody(request).
		SetContext(ctx).
		SetAuthToken(token).
		Post(c.URL + "/sessions/online/" + s.ReferenceNumber + "/invoices")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return newErrorResponse(resp)
	}

	return nil
}
