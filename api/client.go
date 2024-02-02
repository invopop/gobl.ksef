// Package api used for communication with the KSeF API
package api

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
)

// ClientOptFunc defines function for customizing the KSeF client
type ClientOptFunc func(*ClientOpts)

// ClientOpts defines the client parameters
type ClientOpts struct {
	Client           *resty.Client
	URL              string
	ID               string
	Token            string
	SessionToken     string
	SessionReference string
	KeyPath          string
}

func defaultClientOpts() ClientOpts {
	return ClientOpts{
		Client:           resty.New(),
		URL:              "https://ksef-test.mf.gov.pl",
		ID:               "",
		Token:            "",
		SessionToken:     "",
		SessionReference: "",
		KeyPath:          "",
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

// WithID allows customizing the Polish tax id number (NIP)
func WithID(id string) ClientOptFunc {
	return func(o *ClientOpts) {
		o.ID = id
	}
}

// WithToken allows customizing the KSeF authorization token
func WithToken(token string) ClientOptFunc {
	return func(o *ClientOpts) {
		o.Token = token
	}
}

// WithKeyPath allows customizing the public key for KSeF API authorization
func WithKeyPath(keyPath string) ClientOptFunc {
	return func(o *ClientOpts) {
		o.KeyPath = keyPath
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
func NewClient(opts ...ClientOptFunc) *Client {
	o := defaultClientOpts()
	for _, fn := range opts {
		fn(&o)
	}
	o.Client.SetDebug(true)
	return &Client{
		ClientOpts: o,
	}
}

// FetchSessionToken requests new session token
func FetchSessionToken(ctx context.Context, c *Client) error {
	challenge, err := fetchChallenge(ctx, c)
	if err != nil {
		return err
	}

	encryptedToken, err := encryptToken(c, challenge)
	if err != nil {
		return fmt.Errorf("cannot encrypt token: %v", err)
	}

	sessionToken, err := initTokenSession(ctx, c, encryptedToken, challenge.Challenge)
	if err != nil {
		return fmt.Errorf("cannot init session token: %v", err)
	}

	c.SessionToken = sessionToken.SessionToken.Token
	c.SessionReference = sessionToken.ReferenceNumber

	return nil
}

func publicKey(keyPath string) (*rsa.PublicKey, error) {
	key, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("cannot read key file %s: %v", keyPath, err)
	}
	block, _ := pem.Decode(key)
	parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("cannot parse public key: %v", err)
	}
	return parsedKey.(*rsa.PublicKey), nil
}

func encryptToken(c *Client, challenge *AuthorisationChallengeResponse) ([]byte, error) {
	rawToken := fmt.Sprintf("%s|%d", c.Token, challenge.Timestamp.UnixMilli())

	publicKey, err := publicKey(c.KeyPath)
	if err != nil {
		return nil, err
	}

	return rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(rawToken))
}
