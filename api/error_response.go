package api

import (
	"encoding/json"
	"errors"
	"fmt"
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

func newErrorResponse(resp *resty.Response) error {
	msg := fmt.Sprintf("KSeF service error response (Status %s)", resp.Status())

	if resp.StatusCode() >= 500 {
		// 5xx errors don't include an ErrorResponse body
		return errors.New(msg)
	}

	er := new(ErrorResponse)
	if err := json.Unmarshal(resp.Body(), er); err != nil {
		return fmt.Errorf("%s: %s", msg, resp.Body())
	}

	return fmt.Errorf("%s: %w", msg, er)
}
