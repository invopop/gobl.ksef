package api

import (
	"context"
	"fmt"
)

// UploadInvoice uploads a serialized invoice using the provided upload session data.
func (c *Client) UploadInvoice(ctx context.Context, session *UploadSession, invoice []byte) error {
	if session == nil {
		return fmt.Errorf("upload session is nil")
	}
	if session.ReferenceNumber == "" {
		return fmt.Errorf("upload session missing reference number")
	}

	token, err := c.AccessTokenValue(ctx)
	if err != nil {
		return err
	}

	request, err := buildUploadInvoiceRequest(session, invoice)
	if err != nil {
		return err
	}

	resp, err := c.Client.R().
		SetBody(request).
		SetContext(ctx).
		SetAuthToken(token).
		Post(c.URL + "/sessions/online/" + session.ReferenceNumber + "/invoices")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return newErrorResponse(resp)
	}

	return nil
}
