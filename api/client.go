// Package api used for communication with the KSeF API
package api

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/invopop/xmldsig"
)

// ErrCertificateNotConfigured indicates that the client certificate is missing.
var ErrCertificateNotConfigured = errors.New("client certificate is not configured; pass it to api.NewClient")

// ClientOptFunc defines function for customizing the KSeF client
type ClientOptFunc func(*clientOpts)

type Environment string

const (
	EnvironmentProduction Environment = "prod"
	EnvironmentDemo       Environment = "demo"
	EnvironmentTest       Environment = "test"
)

// clientOpts defines the client parameters
type clientOpts struct {
	client            *resty.Client        // Resty client used for making the requests
	url               string               // Base API URL for the requests
	qrUrl             string               // Base API URL for QR code verification
	contextIdentifier *ContextIdentifier   // Identifies the business entity the requests are made for
	certificate       *xmldsig.Certificate // Certificate used for authorization
	accessToken       *apiToken            // Access token used for making most of the requests
	refeshToken       *apiToken            // Refresh token used for refreshing the access token
}

func defaultClientOpts(contextIdentifier *ContextIdentifier, certificate *xmldsig.Certificate) clientOpts {
	return clientOpts{
		client:            resty.New(),
		url:               "https://api-test.ksef.mf.gov.pl/v2",
		qrUrl:             EnvironmentTestQrUrl,
		contextIdentifier: contextIdentifier,
		certificate:       certificate,
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

// WithProductionURL sets the client url to KSeF production
func WithProductionURL(o *clientOpts) {
	o.url = "https://api.ksef.mf.gov.pl/v2"
	o.qrUrl = EnvironmentProductionQrUrl
}

// WithDemoURL sets the client url to KSeF demo
func WithDemoURL(o *clientOpts) {
	o.url = "https://api-demo.ksef.mf.gov.pl/v2"
	o.qrUrl = EnvironmentDemoQrUrl
}

// NewClient returns a KSeF API client
func NewClient(contextIdentifier *ContextIdentifier, certificate *xmldsig.Certificate, opts ...ClientOptFunc) *Client {
	o := defaultClientOpts(contextIdentifier, certificate)
	for _, fn := range opts {
		fn(&o)
	}
	return &Client{
		clientOpts: o,
	}
}

// Authenticate performs the full authorization exchange and stores the resulting tokens on the client.
func (c *Client) Authenticate(ctx context.Context) error {
	if c.certificate == nil {
		return ErrCertificateNotConfigured
	}

	challenge, err := c.fetchChallenge(ctx)
	if err != nil {
		return err
	}

	authResp, err := c.authorizeWithCertificate(ctx, challenge, c.contextIdentifier)
	if err != nil {
		return err
	}
	if authResp.AuthenticationToken == nil {
		return fmt.Errorf("authorization response missing authentication token")
	}

	err = c.pollAuthorizationStatus(ctx, authResp.ReferenceNumber, authResp.AuthenticationToken.Token)
	if err != nil {
		return err
	}

	exchResp, err := c.exchangeToken(ctx, authResp.AuthenticationToken.Token)
	if err != nil {
		return err
	}

	c.accessToken = exchResp.AccessToken
	c.refeshToken = exchResp.RefreshToken

	return nil
}
