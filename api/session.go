package api

import (
	"context"
	"encoding/xml"
	"errors"
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
	var errorResponse ErrorResponse
	resp, err := s.Client.R().
		SetResult(response).
		SetError(&errorResponse).
		SetContext(ctx).
		SetHeader("SessionToken", s.SessionToken).
		Get(s.URL + "/api/online/Session/Terminate")

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(errorResponse.Exception.ExceptionDetailList[0].ExceptionDescription)
	}

	return response, nil
}

// GetSessionStatus gets the session status of the current session
func GetSessionStatus(ctx context.Context, c *Client) (*SessionStatusResponse, error) {
	response := &SessionStatusResponse{}
	var errorResponse ErrorResponse
	resp, err := c.Client.R().
		SetResult(response).
		SetError(&errorResponse).
		SetContext(ctx).
		SetHeader("SessionToken", c.SessionToken).
		Get(c.URL + "/api/online/Session/Status")

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(errorResponse.Exception.ExceptionDetailList[0].ExceptionDescription)
	}

	return response, nil
}

// GetSessionStatusByReference gets the session status by reference number
func GetSessionStatusByReference(ctx context.Context, c *Client) (*SessionStatusByReferenceResponse, error) {
	response := &SessionStatusByReferenceResponse{}
	var errorResponse ErrorResponse
	resp, err := c.Client.R().
		SetResult(response).
		SetError(&errorResponse).
		SetContext(ctx).
		Get(c.URL + "/api/common/Status/" + c.SessionReference)

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(errorResponse.Exception.ExceptionDetailList[0].ExceptionDescription)
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
