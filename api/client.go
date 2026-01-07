// Package api used for communication with the KSeF API
package api

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
)

// ClientOptFunc defines function for customizing the KSeF client
type ClientOptFunc func(*clientOpts)

// clientOpts defines the client parameters
type clientOpts struct {
	client              *resty.Client      // Resty client used for making the requests
	url                 string             // Base API URL for the requests
	contextIdentifier   *ContextIdentifier // Identifies the business entity the requests are made for
	certificatePath     string             // Path to the .p12 / .pfx certificate for KSeF API authorization
	certificatePassword string             // Password to certificate above
	accessToken         *ApiToken          // Access token used for making most of the requests
	refeshToken         *ApiToken          // Refresh token used for refreshing the access token
}

func defaultClientOpts(contextIdentifier *ContextIdentifier, certificatePath string) clientOpts {
	return clientOpts{
		client:              resty.New(),
		url:                 "https://ksef-test.mf.gov.pl/api/v2",
		contextIdentifier:   contextIdentifier,
		certificatePath:     certificatePath,
		certificatePassword: "",
	}
}

// Client defines KSeF client
type Client struct {
	clientOpts
}

// WithClient allows to customize the http client used for making the requests
func WithClient(client *resty.Client) ClientOptFunc {
	return func(o *clientOpts) {
		o.client = client
	}
}

// WithDebugClient uses a more verbose client
func WithDebugClient() ClientOptFunc {
	c := resty.New()
	c.SetDebug(true)
	return func(o *clientOpts) {
		o.client = c
	}
}

// WithCertificatePassword allows passing the password to the certificate above
func WithCertificatePassword(password string) ClientOptFunc {
	return func(o *clientOpts) {
		o.certificatePassword = password
	}
}

// WithProductionURL sets the client url to KSeF production
func WithProductionURL(o *clientOpts) {
	o.url = "https://ksef.mf.gov.pl/api/v2"
}

// WithDemoURL sets the client url to KSeF demo
func WithDemoURL(o *clientOpts) {
	o.url = "https://ksef-demo.mf.gov.pl/api/v2"
}

// NewClient returns a KSeF API client
func NewClient(contextIdentifier *ContextIdentifier, certificatePath string, opts ...ClientOptFunc) *Client {
	o := defaultClientOpts(contextIdentifier, certificatePath)
	for _, fn := range opts {
		fn(&o)
	}
	return &Client{
		clientOpts: o,
	}
}

// Performs the complete authentication flow
func (c *Client) Authenticate(ctx context.Context) error {
	challenge, err := fetchChallenge(ctx, c)
	if err != nil {
		return err
	}

	authResp, err := authorizeWithCertificate(ctx, c, challenge, c.contextIdentifier)
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

	c.accessToken = exchResp.AccessToken
	c.refeshToken = exchResp.RefreshToken

	return nil
}
