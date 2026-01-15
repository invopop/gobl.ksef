package api

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
)

var (
	ErrInvoiceSubjectTypeRequired = errors.New("subject type is required")
	ErrInvoiceDateFromRequired    = errors.New("date from is required")
	ErrInvoiceInvalidSortOrder    = errors.New("invalid sort order")
	ErrInvoicesTruncated          = errors.New("invoice query truncated, reduce the date range")
)

// InvoiceSortOrder defines the order invoices are returned from the API.
type InvoiceSortOrder string

const (
	// InvoiceSortOrderAscending returns the oldest invoices first.
	InvoiceSortOrderAscending InvoiceSortOrder = "asc"
	// InvoiceSortOrderDescending returns the newest invoices first.
	InvoiceSortOrderDescending InvoiceSortOrder = "desc"
)

// InvoiceSubjectType identifies the party used for invoice queries.
type InvoiceSubjectType string

const (
	InvoiceSubjectTypeSupplier   InvoiceSubjectType = "Subject1"          // outgoing (seller)
	InvoiceSubjectTypeCustomer   InvoiceSubjectType = "Subject2"          // incoming (buyer)
	InvoiceSubjectTypeThirdParty InvoiceSubjectType = "Subject3"          // third party
	InvoiceSubjectTypeAuthorized InvoiceSubjectType = "SubjectAuthorized" // when acting on behalf of another party
)

// ListInvoicesParams describe how invoice metadata should be queried.
type ListInvoicesParams struct {
	SubjectType InvoiceSubjectType
	DateFrom    string
	DateTo      string
	SortOrder   InvoiceSortOrder
	PageOffset  int
	PageSize    int
}

type listInvoicesDateRange struct {
	DateType string `json:"dateType"`
	From     string `json:"from"`
	To       string `json:"to,omitempty"`
}

type listInvoicesRequest struct {
	SubjectType string                `json:"subjectType"`
	DateRange   listInvoicesDateRange `json:"dateRange"`
}

// InvoiceMetadata holds the subset of fields we care about from the metadata endpoint.
type InvoiceMetadata struct {
	KsefNumber           string `json:"ksefNumber"`
	PermanentStorageDate string `json:"permanentStorageDate"`
}

// ListInvoicesPageResponse stores the response returned by ListInvoicesPage.
type ListInvoicesPageResponse struct {
	HasMore     bool              `json:"hasMore"`
	IsTruncated bool              `json:"isTruncated"`
	Invoices    []InvoiceMetadata `json:"invoices"`
}

// ListInvoicesPage calls the metadata endpoint for a single page of results.
func (c *Client) ListInvoicesPage(ctx context.Context, params ListInvoicesParams) (*ListInvoicesPageResponse, error) {
	prepared, err := params.normalize()
	if err != nil {
		return nil, err
	}

	token, err := c.getAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	request := &listInvoicesRequest{
		SubjectType: string(prepared.SubjectType),
		DateRange: listInvoicesDateRange{
			DateType: "PermanentStorage",
			From:     prepared.DateFrom,
			To:       prepared.DateTo,
		},
	}

	response := &ListInvoicesPageResponse{}
	resp, err := c.client.R().
		SetContext(ctx).
		SetAuthToken(token).
		SetBody(request).
		SetResult(response).
		SetQueryParam("sortOrder", string(prepared.SortOrder)).
		SetQueryParam("pageOffset", strconv.Itoa(prepared.PageOffset)).
		SetQueryParam("pageSize", strconv.Itoa(prepared.PageSize)).
		Post(c.url + "/invoices/query/metadata")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, newErrorResponse(resp)
	}

	return response, nil
}

// ListInvoices fetches invoices page by page until the API indicates there is nothing more.
func (c *Client) ListInvoices(ctx context.Context, params ListInvoicesParams) ([]InvoiceMetadata, error) {
	params, err := params.normalize()
	if err != nil {
		return nil, err
	}

	var invoices []InvoiceMetadata
	for {
		response, err := c.ListInvoicesPage(ctx, params)
		if err != nil {
			return nil, err
		}
		if response.IsTruncated {
			return nil, ErrInvoicesTruncated
		}

		invoices = append(invoices, response.Invoices...)
		if !response.HasMore {
			break
		}
		params.PageOffset++
	}

	return invoices, nil
}

func (p ListInvoicesParams) normalize() (ListInvoicesParams, error) {
	if p.SubjectType == "" {
		return p, ErrInvoiceSubjectTypeRequired
	}
	if p.DateFrom == "" {
		return p, ErrInvoiceDateFromRequired
	}
	switch p.SortOrder {
	case "":
		p.SortOrder = InvoiceSortOrderDescending
	case InvoiceSortOrderAscending, InvoiceSortOrderDescending:
	default:
		return p, fmt.Errorf("%w: %s", ErrInvoiceInvalidSortOrder, p.SortOrder)
	}
	if p.PageOffset < 0 {
		p.PageOffset = 0
	}
	if p.PageSize <= 0 {
		p.PageSize = 100
	}
	return p, nil
}

// GetInvoice downloads the XML invoice body for the provided KSeF number.
func (c *Client) GetInvoice(ctx context.Context, ksefNumber string) ([]byte, error) {
	if ksefNumber == "" {
		return nil, fmt.Errorf("ksef number is required")
	}

	token, err := c.getAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.R().
		SetContext(ctx).
		SetAuthToken(token).
		Get(c.url + "/invoices/ksef/" + url.PathEscape(ksefNumber)) // PathEscape only for additional safety, KSeF numbers should be URL-safe
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, newErrorResponse(resp)
	}

	return resp.Body(), nil
}
