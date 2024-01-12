package ksef_api

import (
	"encoding/xml"
	"time"
)

type ErrorResponse struct {
	Exception struct {
		ServiceCtx          string `json:"serviceCtx"`
		ServiceCode         string `json:"serviceCode"`
		ServiceName         string `json:"serviceName"`
		Timestamp           string `json:"timestamp"`
		ReferenceNumber     string `json:"referenceNumber"`
		ExceptionDetailList []struct {
			ExceptionCode        int    `json:"exceptionCode"`
			ExceptionDescription string `json:"exceptionDescription"`
		} `json:"exceptionDetailList"`
	} `json:"exception"`
}

type ContextIdentifier struct {
	Identifier string `json:"identifier"`
	Type       string `json:"type"`
}

type AuthorisationChallengeRequest struct {
	ContextIdentifier *ContextIdentifier `json:"contextIdentifier"`
}

type AuthorisationChallengeResponse struct {
	Timestamp time.Time `json:"timestamp"`
	Challenge string    `json:"challenge"`
}

const (
	XMLNamespace    = "http://ksef.mf.gov.pl/schema/gtw/svc/online/types/2021/10/01/0001"
	XMLNamespace2   = "http://ksef.mf.gov.pl/schema/gtw/svc/types/2021/10/01/0001"
	XMLNamespace3   = "http://ksef.mf.gov.pl/schema/gtw/svc/online/auth/request/2021/10/01/0001"
	XSIType         = "ns2:SubjectIdentifierByCompanyType"
	XSINamespace    = "http://www.w3.org/2001/XMLSchema-instance"
	RootElementName = "ns3:InitSessionTokenRequest"
)

type InitSessionTokenRequest struct {
	Context       *InitSessionTokenContext `xml:"ns3:Context"`
	XMLName       xml.Name
	XMLNamespace  string `xml:"xmlns,attr"`
	XMLNamespace2 string `xml:"xmlns:ns2,attr"`
	XMLNamespace3 string `xml:"xmlns:ns3,attr"`
}

type InitSessionTokenIdentifier struct {
	Identifier string `xml:"ns2:Identifier"`
	Type       string `xml:"xsi:type,attr"`
	Namespace  string `xml:"xmlns:xsi,attr"`
}

type InitSessionTokenDocumentType struct {
	Service  string                    `xml:"ns2:Service"`
	FormCode *InitSessionTokenFormCode `xml:"ns2:FormCode"`
}

type InitSessionTokenFormCode struct {
	SystemCode      string `xml:"ns2:SystemCode"`
	SchemaVersion   string `xml:"ns2:SchemaVersion"`
	TargetNamespace string `xml:"ns2:TargetNamespace"`
	Value           string `xml:"ns2:Value"`
}

type InitSessionTokenContext struct {
	Challenge    string                        `xml:"Challenge"`
	Identifier   *InitSessionTokenIdentifier   `xml:"Identifier"`
	DocumentType *InitSessionTokenDocumentType `xml:"DocumentType"`
	Token        string                        `xml:"Token"`
}

type InitSessionTokenResponse struct {
	Timestamp       string `json:"timestamp"`
	ReferenceNumber string `json:"referenceNumber"`
	SessionToken    struct {
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
	} `json:"sessionToken"`
}

type TerminateSessionResponse struct {
	Timestamp             string `json:"timestamp"`
	ReferenceNumber       string `json:"referenceNumber"`
	ProcessingCode        int    `json:"processingCode"`
	ProcessingDescription string `json:"processingDescription"`
}

type SendInvoiceRequest struct {
	InvoiceHash    *InvoiceHash    `json:"invoiceHash"`
	InvoicePayload *InvoicePayload `json:"invoicePayload"`
}

type InvoiceHash struct {
	HashSHA  *HashSHA `json:"hashSHA"`
	FileSize int      `json:"fileSize"`
}

type HashSHA struct {
	Algorithm string `json:"algorithm"`
	Encoding  string `json:"encoding"`
	Value     string `json:"value"`
}

type InvoicePayload struct {
	Type        string `json:"type"`
	InvoiceBody string `json:"invoiceBody"`
}

type SendInvoiceResponse struct {
	Timestamp              string `json:"timestamp"`
	ReferenceNumber        string `json:"referenceNumber"`
	ProcessingCode         int    `json:"processingCode"`
	ProcessingDescription  string `json:"processingDescription"`
	ElementReferenceNumber string `json:"elementReferenceNumber"`
}

type InvoiceStatusResponse struct {
	Timestamp              string `json:"timestamp"`
	ReferenceNumber        string `json:"referenceNumber"`
	ProcessingCode         int    `json:"processingCode"`
	ProcessingDescription  string `json:"processingDescription"`
	ElementReferenceNumber string `json:"elementReferenceNumber"`
	InvoiceStatus          struct {
		InvoiceNumber        string `json:"invoiceNumber"`
		KsefReferenceNumber  string `json:"ksefReferenceNumber"`
		AcquisitionTimestamp string `json:"acquisitionTimestamp"`
	} `json:"invoiceStatus"`
}

type SessionStatusResponse struct {
	Timestamp             string `json:"timestamp"`
	ReferenceNumber       string `json:"referenceNumber"`
	ProcessingCode        int    `json:"processingCode"`
	ProcessingDescription string `json:"processingDescription"`
	NumberOfElements      int    `json:"numberOfElements"`
	PageSize              int    `json:"pageSize"`
	PageOffset            int    `json:"pageOffset"`
	InvoiceStatusList     []struct {
	} `json:"invoiceStatusList"`
}

type SessionStatusByReferenceResponse struct {
	ProcessingCode        int    `json:"processingCode"`
	ProcessingDescription string `json:"processingDescription"`
	ReferenceNumber       string `json:"referenceNumber"`
	Timestamp             string `json:"timestamp"`
	Upo                   string `json:"upo"`
}
