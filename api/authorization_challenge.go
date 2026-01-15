package api

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var (
	ErrAuthorizationPollingCountExceeded = errors.New("authorization polling count exceeded")
	ErrAuthorizationFailed               = errors.New("authorization failed")
)

// ContextIdentifier defines the context of the authorization (what business entity we're making the request for)
type ContextIdentifier struct {
	Nip        string `xml:"Nip,omitempty"`
	InternalId string `xml:"InternalId,omitempty"`
	NipVatUe   string `xml:"NipVatUe,omitempty"`
	PeppolId   string `xml:"PeppolId,omitempty"`
}

type authorizationChallengeResponse struct {
	Timestamp time.Time `json:"timestamp"`
	Challenge string    `json:"challenge"`
}

type authorizationResponse struct {
	ReferenceNumber     string    `json:"referenceNumber"`
	AuthenticationToken *apiToken `json:"authenticationToken"`
}

type authorizationPollResponse struct {
	Status *authorizationPollStatus `json:"status"`
}

type authorizationPollStatus struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

type exchangeResponse struct {
	AccessToken  *apiToken `json:"accessToken"`
	RefreshToken *apiToken `json:"refreshToken"`
}

type apiToken struct {
	Token      string `json:"token"`
	ValidUntil string `json:"validUntil"`
}

func (c *Client) fetchChallenge(ctx context.Context) (*authorizationChallengeResponse, error) {
	response := &authorizationChallengeResponse{}

	resp, err := c.client.R().
		SetResult(response).
		SetContext(ctx).
		Post(c.url + "/auth/challenge")

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, newErrorResponse(resp)
	}

	return response, nil
}

func (c *Client) authorizeWithCertificate(ctx context.Context, challenge *authorizationChallengeResponse, contextIdentifier *ContextIdentifier) (*authorizationResponse, error) {
	signedRequestStr, err := c.buildSignedAuthorizationRequest(challenge, contextIdentifier)
	if err != nil {
		return nil, err
	}

	response := &authorizationResponse{}
	resp, err := c.client.R().
		SetHeader("Content-Type", "application/xml").
		SetHeader("Accept", "application/json"). // request body is in XML, but response is in JSON
		SetBody(signedRequestStr).
		SetResult(response).
		SetContext(ctx).
		Post(c.url + "/auth/xades-signature")

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, newErrorResponse(resp)
	}

	return response, nil
}

func (c *Client) pollAuthorizationStatus(ctx context.Context, referenceNumber string, authorizationToken string) error {
	attempt := 0
	for {
		attempt++
		if attempt > 30 {
			return ErrAuthorizationPollingCountExceeded
		}

		response := &authorizationPollResponse{}
		resp, err := c.client.R().
			SetHeader("Authorization", "Bearer "+authorizationToken).
			SetResult(response).
			SetContext(ctx).
			Get(c.url + "/auth/" + referenceNumber)
		if err != nil {
			return err
		}
		if resp.IsError() {
			return newErrorResponse(resp)
		}

		if response.Status.Code == 200 {
			return nil
		}
		if response.Status.Code == 100 {
			time.Sleep(2 * time.Second) // TODO add exponential backoff
			continue
		}
		// any other status means that the authorization failed
		return fmt.Errorf("%w: %s", ErrAuthorizationFailed, response.Status.Description)
	}
}

func (c *Client) exchangeToken(ctx context.Context, authorizationToken string) (*exchangeResponse, error) {
	response := &exchangeResponse{}
	resp, err := c.client.R().
		SetHeader("Authorization", "Bearer "+authorizationToken).
		SetResult(response).
		SetContext(ctx).
		Post(c.url + "/auth/token/redeem")

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, newErrorResponse(resp)
	}

	return response, nil
}
