package api

import (
	"context"
	"fmt"
)

// UploadInvoice uploads a serialized invoice using the provided upload session data.
func UploadInvoice(ctx context.Context, session *UploadSession, invoice []byte, s *Client) error {
	if session == nil {
		return fmt.Errorf("upload session is nil")
	}
	if session.ReferenceNumber == "" {
		return fmt.Errorf("upload session missing reference number")
	}

	token, err := s.AccessTokenValue(ctx)
	if err != nil {
		return err
	}

	request, err := buildUploadInvoiceRequest(session, invoice)
	if err != nil {
		return err
	}

	resp, err := s.Client.R().
		SetBody(request).
		SetContext(ctx).
		SetAuthToken(token).
		Post(s.URL + "/sessions/online/" + session.ReferenceNumber + "/invoices")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return newErrorResponse(resp)
	}

	return nil
}
