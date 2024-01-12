package ksef_api

import (
	"encoding/base64"
	"errors"
	"time"
)

func (s *Session) SendInvoice(content []byte, digestBase64 string) (*SendInvoiceResponse, error) {
	contentBase64 := base64.StdEncoding.EncodeToString(content)

	request := SendInvoiceRequest{
		InvoiceHash: &InvoiceHash{
			HashSHA: &HashSHA{
				Algorithm: "SHA-256",
				Encoding:  "Base64",
				Value:     digestBase64,
			},
			FileSize: len(content),
		},
		InvoicePayload: &InvoicePayload{
			Type:        "plain",
			InvoiceBody: contentBase64,
		},
	}
	response := &SendInvoiceResponse{}
	var errorResponse ErrorResponse
	resp, err := s.client.client.R().
		SetResult(&response).
		SetError(&errorResponse).
		SetBody(request).
		SetHeader("Content-Type", "application/json").
		SetHeader("SessionToken", s.token.SessionToken.Token).
		Put(s.client.url + "/api/online/Invoice/Send")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(errorResponse.Exception.ExceptionDetailList[0].ExceptionDescription)
	}
	return response, nil
}

func (s *Session) GetInvoiceStatus(referenceNumber string) (*InvoiceStatusResponse, error) {
	response := &InvoiceStatusResponse{}
	var errorResponse ErrorResponse
	resp, err := s.client.client.R().
		SetResult(response).
		SetError(&errorResponse).
		SetHeader("SessionToken", s.token.SessionToken.Token).
		Get(s.client.url + "/api/online/Invoice/Status/" + referenceNumber)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(errorResponse.Exception.ExceptionDetailList[0].ExceptionDescription)
	}

	return response, nil
}

func (s *Session) WaitUntilInvoiceIsProcessed(referenceNumber string) (*InvoiceStatusResponse, error) {
	for {
		status, err := s.GetInvoiceStatus(referenceNumber)
		if err != nil {
			return nil, err
		}
		if status.ProcessingCode == 200 || status.ProcessingCode == 404 || status.ProcessingCode == 400 {
			return status, nil
		}
		time.Sleep(5 * time.Second)
	}
}
