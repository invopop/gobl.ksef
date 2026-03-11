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

	subjectIdentifierType := "certificateSubject"
	if !isPolishCertificate(c.certificate) || (contextIdentifier != nil && contextIdentifier.NipVatUe != "") {
		subjectIdentifierType = "certificateFingerprint"
	}
	root.CreateElement("SubjectIdentifierType").SetText(subjectIdentifierType)

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

// isPolishCertificate returns true when the certificate's Subject serial
// number starts with a recognised Polish identifier prefix. Foreign
// certificates (e.g. Lithuanian "PASLT-…") will return false, which
// signals that certificateFingerprint auth must be used instead of
// certificateSubject.
func isPolishCertificate(cert *xmldsig.Certificate) bool {
	if cert == nil {
		return true // safe default
	}
	sn := cert.SubjectSerialNumber()
	for _, prefix := range polishCertPrefixes {
		if strings.HasPrefix(sn, prefix) {
			return true
		}
	}
	return false
}
