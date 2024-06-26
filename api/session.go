package api

import (
	"context"
	"encoding/xml"
)

// SessionStatusByReferenceResponse defines the response of the session status
type SessionStatusByReferenceResponse struct {
	ProcessingCode        int    `json:"processingCode"`
	ProcessingDescription string `json:"processingDescription"`
	ReferenceNumber       string `json:"referenceNumber"`
	Timestamp             string `json:"timestamp"`
	Upo                   string `json:"upo"`
}

// SessionStatusResponse defines the response of the session status
type SessionStatusResponse struct {
	Timestamp             string `json:"timestamp"`
	ReferenceNumber       string `json:"referenceNumber"`
	ProcessingCode        int    `json:"processingCode"`
	ProcessingDescription string `json:"processingDescription"`
	NumberOfElements      int    `json:"numberOfElements"`
	PageSize              int    `json:"pageSize"`
	PageOffset            int    `json:"pageOffset"`
	InvoiceStatusList     []struct {
		AcquisitionTimestamp   string `json:"acquisitionTimestamp"`
		ElementReferenceNumber string `json:"elementReferenceNumber"`
		InvoiceNumber          string `json:"invoiceNumber"`
		KSefReferenceNumber    string `json:"ksefReferenceNumber"`
		ProcessingCode         int    `json:"processingCode"`
		ProcessingDescription  string `json:"processingDescription"`
	} `json:"invoiceStatusList"`
}

// TerminateSessionResponse defines the response of the session termination
type TerminateSessionResponse struct {
	Timestamp             string `json:"timestamp"`
	ReferenceNumber       string `json:"referenceNumber"`
	ProcessingCode        int    `json:"processingCode"`
	ProcessingDescription string `json:"processingDescription"`
}

// TerminateSession ends the current session
func TerminateSession(ctx context.Context, s *Client) (*TerminateSessionResponse, error) {
	response := &TerminateSessionResponse{}
	resp, err := s.Client.R().
		SetResult(response).
		SetContext(ctx).
		SetHeader("SessionToken", s.SessionToken).
		Get(s.URL + "/api/online/Session/Terminate")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, newErrorResponse(resp)
	}

	return response, nil
}

// GetSessionStatus gets the session status of the current session
func GetSessionStatus(ctx context.Context, c *Client) (*SessionStatusResponse, error) {
	response := &SessionStatusResponse{}
	resp, err := c.Client.R().
		SetResult(response).
		SetContext(ctx).
		SetHeader("SessionToken", c.SessionToken).
		Get(c.URL + "/api/online/Session/Status")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, newErrorResponse(resp)
	}

	return response, nil
}

// GetSessionStatusByReference gets the session status by reference number
func GetSessionStatusByReference(ctx context.Context, c *Client) (*SessionStatusByReferenceResponse, error) {
	response := &SessionStatusByReferenceResponse{}
	resp, err := c.Client.R().
		SetResult(response).
		SetContext(ctx).
		Get(c.URL + "/api/common/Status/" + c.SessionReference)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, newErrorResponse(resp)
	}

	return response, nil
}

func bytes(d InitSessionTokenRequest) ([]byte, error) {
	bytes, err := xml.MarshalIndent(d, "", "  ")
	if err != nil {
		return nil, err
	}
	return append([]byte(`<?xml version="1.0" encoding="utf-8" standalone="yes"?>`+"\n"), bytes...), nil
}
