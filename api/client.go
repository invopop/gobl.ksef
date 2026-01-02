// Package api used for communication with the KSeF API
package api

import (
	"context"

	"github.com/go-resty/resty/v2"
)

// ClientOptFunc defines function for customizing the KSeF client
type ClientOptFunc func(*ClientOpts)

// ClientOpts defines the client parameters
type ClientOpts struct {
	Client              *resty.Client      // Resty client used for making the requests
	URL                 string             // Base API URL for the requests
	ContextIdentifier   *ContextIdentifier // Identifies the business entity the requests are made for
	CertificatePath     string             // Path to the .p12 / .pfx certificate for KSeF API authorization
	CertificatePassword string             // Password to certificate above
}

func defaultClientOpts(contextIdentifier *ContextIdentifier, certificatePath string) ClientOpts {
	return ClientOpts{
		Client:              resty.New(),
		URL:                 "https://ksef-test.mf.gov.pl",
		ContextIdentifier:   contextIdentifier,
		CertificatePath:     certificatePath,
		CertificatePassword: "",
	}
}

// Client defines KSeF client
type Client struct {
	ClientOpts
}

// WithClient allows to customize the http client used for making the requests
func WithClient(client *resty.Client) ClientOptFunc {
	return func(o *ClientOpts) {
		o.Client = client
	}
}

// WithDebugClient uses a more verbose client
func WithDebugClient() ClientOptFunc {
	c := resty.New()
	c.SetDebug(true)
	return func(o *ClientOpts) {
		o.Client = c
	}
}

// WithCertificatePassword allows passing the password to the certificate above
func WithCertificatePassword(password string) ClientOptFunc {
	return func(o *ClientOpts) {
		o.CertificatePassword = password
	}
}

// WithProductionURL sets the client url to KSeF production
func WithProductionURL(o *ClientOpts) {
	o.URL = "https://ksef.mf.gov.pl"
}

// WithDemoURL sets the client url to KSeF demo
func WithDemoURL(o *ClientOpts) {
	o.URL = "https://ksef-demo.mf.gov.pl"
}

// NewClient returns a KSeF API client
func NewClient(contextIdentifier *ContextIdentifier, certificatePath string, opts ...ClientOptFunc) *Client {
	o := defaultClientOpts(contextIdentifier, certificatePath)
	for _, fn := range opts {
		fn(&o)
	}
	return &Client{
		ClientOpts: o,
	}
}

// Performs the complete authentication flow
func (c *Client) Authenticate(ctx context.Context, contextIdentifier *ContextIdentifier) (*ExchangeResponse, error) {
	challenge, err := fetchChallenge(ctx, c)
	if err != nil {
		return nil, err
	}

	authResp, err := authorizeWithCertificate(ctx, c, challenge, contextIdentifier)
	if err != nil {
		return nil, err
	}

	err = pollAuthorizationStatus(ctx, c, authResp.ReferenceNumber, authResp.AuthenticationToken)
	if err != nil {
		return nil, err
	}

	exchResp, err := exchangeToken(ctx, c, authResp.AuthenticationToken)
	if err != nil {
		return nil, err
	}

	// TODO: save tokens to the client
	return exchResp, nil
}
