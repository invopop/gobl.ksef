package api

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
)

// InvoiceStatusResponse defines the post invocie response structure
type InvoiceStatusResponse struct {
	Timestamp              string `json:"timestamp"`
	ReferenceNumber        string `json:"referenceNumber"`
	ProcessingCode         int    `json:"processingCode"`
	ProcessingDescription  string `json:"processingDescription"`
	ElementReferenceNumber string `json:"elementReferenceNumber"`
	InvoiceStatus          struct {
		InvoiceNumber        string `json:"invoiceNumber"`
		KsefReferenceNumber  string `json:"ksefReferenceNumber"`
		AcquisitionTimestamp string `json:"acquisitionTimestamp"`
	} `json:"invoiceStatus"`
}

// SendInvoiceResponse defines the post invocie response structure
type SendInvoiceResponse struct {
	Timestamp              string `json:"timestamp"`
	ReferenceNumber        string `json:"referenceNumber"`
	ProcessingCode         int    `json:"processingCode"`
	ProcessingDescription  string `json:"processingDescription"`
	ElementReferenceNumber string `json:"elementReferenceNumber"`
}

// SendInvoiceRequest defines the post invocie request structure
type SendInvoiceRequest struct {
	InvoiceHash    *InvoiceHash    `json:"invoiceHash"`
	InvoicePayload *InvoicePayload `json:"invoicePayload"`
}

// InvoicePayload defines the InvoicePayload part of the post invocie request
type InvoicePayload struct {
	Type        string `json:"type"`
	InvoiceBody string `json:"invoiceBody"`
}

// InvoiceHash defines the InvoiceHash part of the post invocie request
type InvoiceHash struct {
	HashSHA  *HashSHA `json:"hashSHA"`
	FileSize int      `json:"fileSize"`
}

// HashSHA defines the HashSHA part of the post invocie request
type HashSHA struct {
	Algorithm string `json:"algorithm"`
	Encoding  string `json:"encoding"`
	Value     string `json:"value"`
}

// SendInvoice puts the invoice to the KSeF API
func SendInvoice(ctx context.Context, c *Client, data []byte) (*SendInvoiceResponse, error) {
	contentBase64 := base64.StdEncoding.EncodeToString(data)

	request := SendInvoiceRequest{
		InvoiceHash: &InvoiceHash{
			HashSHA: &HashSHA{
				Algorithm: "SHA-256",
				Encoding:  "Base64",
				Value:     digestBase64(data),
			},
			FileSize: len(data),
		},
		InvoicePayload: &InvoicePayload{
			Type:        "plain",
			InvoiceBody: contentBase64,
		},
	}
	response := &SendInvoiceResponse{}
	var errorResponse ErrorResponse
	resp, err := c.Client.R().
		SetResult(&response).
		SetError(&errorResponse).
		SetBody(request).
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("SessionToken", c.SessionToken).
		Put(c.URL + "/api/online/Invoice/Send")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(errorResponse.Exception.ExceptionDetailList[0].ExceptionDescription)
	}
	return response, nil
}

// FetchInvoiceStatus gets the status of the invoice being processed
func FetchInvoiceStatus(ctx context.Context, c *Client, referenceNumber string) (*InvoiceStatusResponse, error) {
	response := &InvoiceStatusResponse{}
	var errorResponse ErrorResponse
	resp, err := c.Client.R().
		SetResult(response).
		SetError(&errorResponse).
		SetHeader("SessionToken", c.SessionToken).
		SetContext(ctx).
		Get(c.URL + "/api/online/Invoice/Status/" + referenceNumber)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(errorResponse.Exception.ExceptionDetailList[0].ExceptionDescription)
	}

	return response, nil
}

func digestBase64(content []byte) string {
	digest := sha256.Sum256(content)
	return base64.StdEncoding.EncodeToString(digest[:])
}
