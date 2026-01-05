package api

import (
	"context"
	"encoding/xml"
	"fmt"
	"time"

	"github.com/invopop/xmldsig"
)

// AuthorizationChallengeResponse defines the authorization challenge response - first step of the session initialization
// This request doesn't have a body
type AuthorizationChallengeResponse struct {
	Timestamp time.Time `json:"timestamp"`
	Challenge string    `json:"challenge"`
}

// Second step involves submitting the XAdES-signed XML with the challenge and context
type AuthorizationRequest struct {
	XMLName               xml.Name           `xml:"http://ksef.mf.gov.pl/auth/token/2.0 AuthTokenRequest"`
	XmlnsXsi              string             `xml:"xmlns:xsi,attr"`
	XmlnsXsd              string             `xml:"xmlns:xsd,attr"`
	Xmlns                 string             `xml:"xmlns,attr"` // note that it must be after the previous attributes
	Challenge             string             `xml:"Challenge"`
	ContextIdentifier     *ContextIdentifier `xml:"ContextIdentifier"`
	SubjectIdentifierType string             `xml:"SubjectIdentifierType"` // certificateSubject or certificateFingerprint
	Signature             *xmldsig.Signature `xml:"ds:Signature,omitempty"`
}

// ContextIdentifier defines the context of the authorization (what business entity we're making the request for)
type ContextIdentifier struct {
	Nip        string `xml:"Nip,omitempty"`
	InternalId string `xml:"InternalId,omitempty"`
	NipVatUe   string `xml:"NipVatUe,omitempty"`
	PeppolId   string `xml:"PeppolId,omitempty"`
}

// AuthorizationResponse defines the authorization response - second step of the session initialization
// While the request is in XML, the response is in JSON
type AuthorizationResponse struct {
	ReferenceNumber     string `json:"referenceNumber"`
	AuthenticationToken string `json:"authenticationToken"`
}

// AuthorizationPollResponse defines response when polling the status of authorization initialized above
type AuthorizationPollResponse struct {
	Status *AuthorizationPollStatus `json:"status"`
}

type AuthorizationPollStatus struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

// ExchangeResponse defines the exchange response, where authenticationToken is used to obtain accessToken and refreshToken
// This request doesn't have a body

type ExchangeResponse struct {
	AccessToken  *ApiToken `json:"accessToken"`
	RefreshToken *ApiToken `json:"refreshToken"`
}

type ApiToken struct {
	Token      string `json:"token"`
	ValidUntil string `json:"validUntil"`
}

func fetchChallenge(ctx context.Context, c *Client) (*AuthorizationChallengeResponse, error) {
	response := &AuthorizationChallengeResponse{}

	resp, err := c.Client.R().
		SetResult(response).
		SetContext(ctx).
		Post(c.URL + "/auth/challenge")

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, newErrorResponse(resp)
	}

	return response, nil
}

func authorizeWithCertificate(ctx context.Context, c *Client, challenge *AuthorizationChallengeResponse, contextIdentifier *ContextIdentifier) (*AuthorizationResponse, error) {
	signedRequestStr, err := buildSignedAuthorizationRequest(c, challenge, contextIdentifier)
	if err != nil {
		return nil, err
	}

	// Uncomment for debugging:

	// fmt.Println(string(signedRequestStr))
	// err = os.WriteFile("output.xml", signedRequestStr, 0644)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to write output.xml: %w", err)
	// }

	response := &AuthorizationResponse{}
	resp, err := c.Client.R().
		SetHeader("Content-Type", "application/xml").
		SetHeader("Accept", "application/json"). // request body is in XML, but response is in JSON
		SetBody(signedRequestStr).
		SetResult(response).
		SetContext(ctx).
		Post(c.URL + "/auth/xades-signature")

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, newErrorResponse(resp)
	}

	return response, nil
}

func pollAuthorizationStatus(ctx context.Context, c *Client, referenceNumber string, authorizationToken string) error {
	attempt := 0
	for {
		attempt++
		if attempt > 30 {
			return fmt.Errorf("authorization polling count exceeded")
		}

		response := &AuthorizationPollResponse{}
		resp, err := c.Client.R().
			SetHeader("Authorization", "Bearer "+authorizationToken).
			SetResult(response).
			SetContext(ctx).
			Get(c.URL + "/auth/" + referenceNumber)
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
		return fmt.Errorf("authorization failed: %s", response.Status.Description)
	}
}

func exchangeToken(ctx context.Context, c *Client, authorizationToken string) (*ExchangeResponse, error) {
	response := &ExchangeResponse{}
	resp, err := c.Client.R().
		SetHeader("Authorization", "Bearer "+authorizationToken).
		SetResult(response).
		SetContext(ctx).
		Post(c.URL + "/auth/token/redeem")

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, newErrorResponse(resp)
	}

	return response, nil
}
