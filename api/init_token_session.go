package api

import (
	"context"
	"encoding/base64"
	"encoding/xml"
)

// InitSessionTokenResponse defines the token session initialization response structure
type InitSessionTokenResponse struct {
	Timestamp       string        `json:"timestamp"`
	ReferenceNumber string        `json:"referenceNumber"`
	SessionToken    *SessionToken `json:"sessionToken"`
}

// SessionToken defines the session token part of the  token session initialization response
type SessionToken struct {
	Token   string `json:"token"`
	Context struct {
		ContextIdentifier struct {
			Type       string `json:"type"`
			Identifier string `json:"identifier"`
		} `json:"contextIdentifier"`
		ContextName struct {
			Type      string `json:"type"`
			TradeName string `json:"tradeName"`
			FullName  string `json:"fullName"`
		} `json:"contextName"`
		CredentialsRoleList []struct {
			Type            string `json:"type"`
			RoleType        string `json:"roleType"`
			RoleDescription string `json:"roleDescription"`
		} `json:"credentialsRoleList"`
	} `json:"context"`
}

// InitSessionTokenRequest defines the structure of the token session initialization
type InitSessionTokenRequest struct {
	Context       *InitSessionTokenContext `xml:"ns3:Context"`
	XMLName       xml.Name
	XMLNamespace  string `xml:"xmlns,attr"`
	XMLNamespace2 string `xml:"xmlns:ns2,attr"`
	XMLNamespace3 string `xml:"xmlns:ns3,attr"`
}

// InitSessionTokenContext defines the Context part of the token session initialization
type InitSessionTokenContext struct {
	Challenge    string                        `xml:"Challenge"`
	Identifier   *InitSessionTokenIdentifier   `xml:"Identifier"`
	DocumentType *InitSessionTokenDocumentType `xml:"DocumentType"`
	Token        string                        `xml:"Token"`
}

// InitSessionTokenIdentifier defines the Identifier part of the token session initialization
type InitSessionTokenIdentifier struct {
	Identifier string `xml:"ns2:Identifier"`
	Type       string `xml:"xsi:type,attr"`
	Namespace  string `xml:"xmlns:xsi,attr"`
}

// InitSessionTokenDocumentType defines the DocumentType part of the token session initialization
type InitSessionTokenDocumentType struct {
	Service  string                    `xml:"ns2:Service"`
	FormCode *InitSessionTokenFormCode `xml:"ns2:FormCode"`
}

// InitSessionTokenFormCode defines the FormCode part of the token session initialization
type InitSessionTokenFormCode struct {
	SystemCode      string `xml:"ns2:SystemCode"`
	SchemaVersion   string `xml:"ns2:SchemaVersion"`
	TargetNamespace string `xml:"ns2:TargetNamespace"`
	Value           string `xml:"ns2:Value"`
}

const (
	// XMLNamespace namespace setting for token initialization XML
	XMLNamespace = "http://ksef.mf.gov.pl/schema/gtw/svc/online/types/2021/10/01/0001"
	// XMLNamespace2 namespace setting for token initialization XML
	XMLNamespace2 = "http://ksef.mf.gov.pl/schema/gtw/svc/types/2021/10/01/0001"
	// XMLNamespace3 namespace setting for token initialization XML
	XMLNamespace3 = "http://ksef.mf.gov.pl/schema/gtw/svc/online/auth/request/2021/10/01/0001"
	// XSIType namespace setting for token initialization XML
	XSIType = "ns2:SubjectIdentifierByCompanyType"
	// XSINamespace namespace setting for token initialization XML
	XSINamespace = "http://www.w3.org/2001/XMLSchema-instance"
	// RootElementName root element name for token initialization XML
	RootElementName = "ns3:InitSessionTokenRequest"
)

func initTokenSession(ctx context.Context, c *Client, token []byte, challenge string) (*InitSessionTokenResponse, error) {
	response := &InitSessionTokenResponse{}

	request := &InitSessionTokenRequest{
		XMLNamespace:  XMLNamespace,
		XMLNamespace2: XMLNamespace2,
		XMLNamespace3: XMLNamespace3,
		XMLName:       xml.Name{Local: RootElementName},
		Context: &InitSessionTokenContext{
			Identifier: &InitSessionTokenIdentifier{
				Namespace:  XSINamespace,
				Type:       XSIType,
				Identifier: c.ID,
			},
			Challenge: challenge,
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
	bytes, _ := bytes(*request)

	resp, err := c.Client.R().
		SetResult(response).
		SetBody(bytes).
		SetContext(ctx).
		SetHeader("Content-Type", "application/octet-stream; charset=utf-8").
		Post(c.URL + "/api/online/Session/InitToken")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, newErrorResponse(resp)
	}

	return response, nil
}
