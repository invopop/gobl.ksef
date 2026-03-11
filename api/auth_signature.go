package api

import (
	"encoding/xml"
	"errors"
	"strings"

	"github.com/beevik/etree"
	"github.com/invopop/xmldsig"
)

var ErrCertificatePrivateKeyNotRSA = errors.New("certificate private key is not RSA, goxades only supports RSA")

func (c *Client) buildSignedAuthorizationRequest(challenge *authorizationChallengeResponse, contextIdentifier *ContextIdentifier) ([]byte, error) {
	// 1. Assembly the XML request - the signing library requires XML as an etree object

	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="utf-8"`)

	root := doc.CreateElement("AuthTokenRequest")
	root.CreateAttr("xmlns:xsi", "http://www.w3.org/2001/XMLSchema-instance")
	root.CreateAttr("xmlns:xsd", "http://www.w3.org/2001/XMLSchema")
	root.CreateAttr("xmlns", "http://ksef.mf.gov.pl/auth/token/2.0")

	root.CreateElement("Challenge").SetText(challenge.Challenge)

	ctx := root.CreateElement("ContextIdentifier")
	if contextIdentifier.Nip != "" {
		ctx.CreateElement("Nip").SetText(contextIdentifier.Nip)
	}
	if contextIdentifier.NipVatUe != "" {
		ctx.CreateElement("NipVatUe").SetText(contextIdentifier.NipVatUe)
	}
	if contextIdentifier.InternalId != "" {
		ctx.CreateElement("InternalId").SetText(contextIdentifier.InternalId)
	}
	if contextIdentifier.PeppolId != "" {
		ctx.CreateElement("PeppolId").SetText(contextIdentifier.PeppolId)
	}

	root.CreateElement("SubjectIdentifierType").SetText(
		subjectIdentifierType(c.certificate, contextIdentifier),
	)

	unsignedXML, err := doc.WriteToString()
	if err != nil {
		return nil, err
	}

	// Sign
	signature, err := xmldsig.Sign([]byte(unsignedXML),
		xmldsig.WithCertificate(c.certificate),
		xmldsig.WithKSeF(),
	)
	if err != nil {
		return nil, err
	}

	// attach signature to XML
	signatureXML, err := xml.Marshal(signature)
	if err != nil {
		return nil, err
	}
	sigDoc := etree.NewDocument()
	if err := sigDoc.ReadFromBytes(signatureXML); err != nil {
		return nil, err
	}
	root.AddChild(sigDoc.Root())

	signedXML, err := doc.WriteToString()
	if err != nil {
		return nil, err
	}

	return []byte(signedXML), nil
}

// polishCertPrefixes are the Subject.SerialNumber prefixes that identify
// a certificate as Polish (carrying a NIP or PESEL).
var polishCertPrefixes = []string{"TINPL", "PNOPL", "PESEL", "NIP"}

// subjectIdentifierType determines which SubjectIdentifierType to use in the
// KSeF auth request. It returns "certificateFingerprint" when the certificate
// carries a foreign identifier or when authenticating via NipVatUe
// (EU VAT context). Otherwise it returns "certificateSubject".
func subjectIdentifierType(cert *xmldsig.Certificate, id *ContextIdentifier) string {
	if isForeignCertificate(cert) || (id != nil && id.NipVatUe != "") {
		return "certificateFingerprint"
	}
	return "certificateSubject"
}

// isForeignCertificate returns true when the certificate's Subject serial
// number is present and does NOT start with a recognised Polish identifier
// prefix. This indicates a foreign certificate (e.g. Lithuanian "PASLT-…")
// that requires certificateFingerprint authentication.
//
// Returns false (non-foreign) when:
//   - cert is nil — safe default to avoid breaking callers without a certificate.
//   - SerialNumber is empty — e.g. Polish CCK KSeF certs that carry the NIP
//     in OID 2.5.4.97 instead of SerialNumber.
func isForeignCertificate(cert *xmldsig.Certificate) bool {
	if cert == nil {
		return false
	}
	return isNonPolishSubjectSerialNumber(cert.SubjectSerialNumber())
}

// isNonPolishSubjectSerialNumber returns true when the Subject serial number
// is non-empty and does not start with any recognised Polish prefix (TINPL,
// PNOPL, PESEL, NIP). An empty serial number returns false (safe default).
func isNonPolishSubjectSerialNumber(subjectSerialNumber string) bool {
	if subjectSerialNumber == "" {
		return false
	}
	for _, prefix := range polishCertPrefixes {
		if strings.HasPrefix(subjectSerialNumber, prefix) {
			return false
		}
	}
	return true
}
