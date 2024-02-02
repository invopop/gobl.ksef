package api

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
