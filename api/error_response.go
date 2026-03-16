package api

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
)

// ErrorResponse parses error responses
type ErrorResponse struct {
	Exception struct {
		ServiceCtx          string `json:"serviceCtx"`
		ServiceCode         string `json:"serviceCode"`
		ServiceName         string `json:"serviceName"`
		Timestamp           string `json:"timestamp"`
		ReferenceNumber     string `json:"referenceNumber"`
		ExceptionDetailList []struct {
			ExceptionCode        int    `json:"exceptionCode"`
			ExceptionDescription string `json:"exceptionDescription"`
		} `json:"exceptionDetailList"`
	} `json:"exception"`
}

// Error implements the error interface
func (e ErrorResponse) Error() string {
	msgs := make([]string, len(e.Exception.ExceptionDetailList))
	for i, detail := range e.Exception.ExceptionDetailList {
		msgs[i] = fmt.Sprintf("Code %d: %s", detail.ExceptionCode, detail.ExceptionDescription)
	}

	return strings.Join(msgs, ", ")
}

// ServiceError wraps an HTTP error response from the KSeF API,
// preserving the status code and Retry-After header for callers to inspect.
type ServiceError struct {
	StatusCode int
	Status     string
	RetryAfter int32 // seconds, from Retry-After header (0 if not present)
	Err        error // underlying ErrorResponse or raw body message
}

// Error implements the error interface.
func (e *ServiceError) Error() string {
	msg := fmt.Sprintf("KSeF service error response (Status %s)", e.Status)
	if e.Err != nil {
		return fmt.Sprintf("%s: %s", msg, e.Err.Error())
	}
	return msg
}

// Unwrap returns the underlying error for use with errors.Is/As.
func (e *ServiceError) Unwrap() error {
	return e.Err
}

// IsTransient returns true if this is a transient error (429 rate limit or 5xx server error).
func (e *ServiceError) IsTransient() bool {
	return e.StatusCode == 429 || e.StatusCode >= 500
}

func newErrorResponse(resp *resty.Response) error {
	se := &ServiceError{
		StatusCode: resp.StatusCode(),
		Status:     resp.Status(),
	}

	// Parse Retry-After header (delta-seconds as per KSeF docs).
	if ra := resp.Header().Get("Retry-After"); ra != "" {
		if v, err := strconv.ParseInt(ra, 10, 32); err == nil && v > 0 {
			se.RetryAfter = int32(v)
		}
	}

	if resp.StatusCode() >= 500 {
		// 5xx errors don't include an ErrorResponse body
		return se
	}

	er := new(ErrorResponse)
	if err := json.Unmarshal(resp.Body(), er); err != nil {
		se.Err = fmt.Errorf("invalid JSON response (%w): %s", err, resp.Body())
		return se
	}

	se.Err = er
	return se
}

