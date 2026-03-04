package api

import (
	"encoding/xml"
	"errors"

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
	if c.useCertificateFingerprint || (contextIdentifier != nil && contextIdentifier.NipVatUe != "") {
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
