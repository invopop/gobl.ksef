package ksef_api

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	// HTTP client used to communicate with the API.
	client *resty.Client
	url    string
}

type Session struct {
	token     InitSessionTokenResponse
	client    *Client
	challenge string
}

func (c *Client) GetChallenge(identifier string) (*AuthorisationChallengeResponse, error) {
	response := &AuthorisationChallengeResponse{}
	var errorResponse ErrorResponse

	request := &AuthorisationChallengeRequest{
		ContextIdentifier: &ContextIdentifier{
			Identifier: identifier,
			Type:       "onip",
		},
	}

	resp, err := c.client.R().
		SetResult(response).
		SetBody(request).
		Post(c.url + "/api/online/Session/AuthorisationChallenge")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(errorResponse.Exception.ExceptionDetailList[0].ExceptionDescription)
	}
	return response, nil

}

func (s *Session) InitSessionToken(identifier string, token []byte, keyPath string) (*InitSessionTokenResponse, error) {
	response := &InitSessionTokenResponse{}
	var errorResponse ErrorResponse

	request := &InitSessionTokenRequest{
		XMLNamespace:  XMLNamespace,
		XMLNamespace2: XMLNamespace2,
		XMLNamespace3: XMLNamespace3,
		XMLName:       xml.Name{Local: RootElementName},
		Context: &InitSessionTokenContext{
			Identifier: &InitSessionTokenIdentifier{
				Namespace:  XSINamespace,
				Type:       XSIType,
				Identifier: identifier,
			},
			Challenge: s.challenge,
			DocumentType: &InitSessionTokenDocumentType{
				Service: "KSeF",
				FormCode: &InitSessionTokenFormCode{
					SystemCode:      "FA (2)",
					SchemaVersion:   "1-0E",
					TargetNamespace: "http://crd.gov.pl/wzor/2023/06/29/12648",
					Value:           "FA",
				},
			},
			Token: base64.StdEncoding.EncodeToString(token),
		},
	}
	bytes, _ := Bytes(*request)
	println(string(bytes))

	resp, err := s.client.client.R().
		SetResult(response).
		SetError(&errorResponse).
		SetBody(bytes).
		SetHeader("Content-Type", "application/octet-stream; charset=utf-8").
		Post(s.client.url + "/api/online/Session/InitToken")

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(errorResponse.Exception.ExceptionDetailList[0].ExceptionDescription)
	}
	return response, nil
}

func NewClient(url string) *Client {
	return &Client{
		client: resty.New(),
		url:    url,
	}
}
func (c *Client) NewSession(identifier string, token string, keyPath string) (*Session, error) {
	session := &Session{
		client: c,
	}
	session.client.client.SetDebug(true)
	challenge, err := c.GetChallenge(identifier)
	if err != nil {
		return nil, err
	}
	println(challenge.Timestamp.UnixMilli())
	session.challenge = challenge.Challenge
	rawToken := fmt.Sprintf("%s|%d", token, challenge.Timestamp.UnixMilli())

	key, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("cannot read key file %s: %v", keyPath, err)
	}

	block, _ := pem.Decode(key)
	parsedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("cannot parse public key: %v", err)
	}
	publicKey := parsedKey.(*rsa.PublicKey)
	encryptedToken, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(rawToken))
	if err != nil {
		return nil, fmt.Errorf("cannot encrypt token: %v", err)
	}
	sessionToken, err := session.InitSessionToken(identifier, encryptedToken, keyPath)
	if err != nil {
		return nil, fmt.Errorf("cannot init session token: %v", err)
	}
	session.token = *sessionToken

	return session, nil
}

func (s *Session) TerminateSession() (*TerminateSessionResponse, error) {
	response := &TerminateSessionResponse{}
	var errorResponse ErrorResponse
	resp, err := s.client.client.R().
		SetResult(response).
		SetError(&errorResponse).
		SetHeader("SessionToken", s.token.SessionToken.Token).
		Get(s.client.url + "/api/online/Session/Terminate")

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(errorResponse.Exception.ExceptionDetailList[0].ExceptionDescription)
	}

	return response, nil
}

func (s *Session) GetSessionStatus() (*SessionStatusResponse, error) {
	response := &SessionStatusResponse{}
	var errorResponse ErrorResponse
	resp, err := s.client.client.R().
		SetResult(response).
		SetError(&errorResponse).
		SetHeader("SessionToken", s.token.SessionToken.Token).
		Get(s.client.url + "/api/online/Session/Status")

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(errorResponse.Exception.ExceptionDetailList[0].ExceptionDescription)
	}

	return response, nil
}

func (s *Session) GetSessionStatusByReference(referenceNumber string) (*SessionStatusByReferenceResponse, error) {
	response := &SessionStatusByReferenceResponse{}
	var errorResponse ErrorResponse
	resp, err := s.client.client.R().
		SetResult(response).
		SetError(&errorResponse).
		Get(s.client.url + "/api/common/Status/" + referenceNumber)

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, errors.New(errorResponse.Exception.ExceptionDetailList[0].ExceptionDescription)
	}

	return response, nil
}

func (s *Session) WaitUntilSessionIsTerminated() (*SessionStatusByReferenceResponse, error) {
	_, err := s.TerminateSession()
	if err != nil {
		return nil, err
	}
	for {
		status, err := s.GetSessionStatusByReference(s.token.ReferenceNumber)

		if err != nil {
			return nil, err
		}
		if status.ProcessingCode == 200 {
			return status, nil
		}
		time.Sleep(5 * time.Second)
	}
}

func Bytes(d InitSessionTokenRequest) ([]byte, error) {
	bytes, err := xml.MarshalIndent(d, "", "  ")
	if err != nil {
		return nil, err
	}
	return append([]byte(`<?xml version="1.0" encoding="utf-8" standalone="yes"?>`+"\n"), bytes...), nil
}
