// Package api used for communication with the KSeF API
package api

import (
	"context"
	"fmt"

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
	AccessToken         *ApiToken          // Access token used for making most of the requests
	RefreshToken        *ApiToken          // Refresh token used for refreshing the access token
}

func defaultClientOpts(contextIdentifier *ContextIdentifier, certificatePath string) ClientOpts {
	return ClientOpts{
		Client:              resty.New(),
		URL:                 "https://ksef-test.mf.gov.pl/api/v2",
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
	o.URL = "https://ksef.mf.gov.pl/api/v2"
}

// WithDemoURL sets the client url to KSeF demo
func WithDemoURL(o *ClientOpts) {
	o.URL = "https://ksef-demo.mf.gov.pl/api/v2"
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
func (c *Client) Authenticate(ctx context.Context) error {
	challenge, err := fetchChallenge(ctx, c)
	if err != nil {
		return err
	}

	authResp, err := authorizeWithCertificate(ctx, c, challenge, c.ContextIdentifier)
	if err != nil {
		return err
	}
	if authResp.AuthenticationToken == nil {
		return fmt.Errorf("authorization response missing authentication token")
	}

	err = pollAuthorizationStatus(ctx, c, authResp.ReferenceNumber, authResp.AuthenticationToken.Token)
	if err != nil {
		return err
	}

	exchResp, err := exchangeToken(ctx, c, authResp.AuthenticationToken.Token)
	if err != nil {
		return err
	}

	c.AccessToken = exchResp.AccessToken
	c.RefreshToken = exchResp.RefreshToken

	return nil
}
