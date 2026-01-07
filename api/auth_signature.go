package api

import (
	"crypto"
	"crypto/rsa"
	"fmt"
	"os"

	xades "github.com/MieszkoGulinski/goxades"
	"github.com/beevik/etree"
	dsig "github.com/russellhaering/goxmldsig"
	"software.sslmate.com/src/go-pkcs12"
)

func buildSignedAuthorizationRequest(c *Client, challenge *AuthorizationChallengeResponse, contextIdentifier *ContextIdentifier) ([]byte, error) {
	// I tried to use the github.com/invopop/xmldsig library, but it doesn't work, as it has many options hardcoded that aren't compatible with the KSEF API

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
	if contextIdentifier != nil && contextIdentifier.NipVatUe != "" {
		subjectIdentifierType = "certificateFingerprint"
	}
	root.CreateElement("SubjectIdentifierType").SetText(subjectIdentifierType)

	// 2. Read the certificate from file (.p12 / .pfx) and extract private key and certificate
	p12Bytes, err := os.ReadFile(c.certificatePath)
	if err != nil {
		return nil, err
	}

	privateKey, cert, _, err :=
		pkcs12.DecodeChain(p12Bytes, c.certificatePassword)
	if err != nil {
		return nil, err
	}

	rsaKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("certificate private key is not RSA, goxades only supports RSA")
	}

	store := xades.MemoryX509KeyStore{
		PrivateKey: rsaKey,
		Cert:       cert,
		CertBinary: cert.Raw,
	}

	// 3. Sign the XML request
	canonicalizerSignedInfo := dsig.MakeC14N10ExclusiveCanonicalizerWithPrefixList("") // http://www.w3.org/TR/2001/REC-xml-c14n-20010315
	// Using exclusive canonicalizer resulted in xsi and xsd attributes disappearing from AuthTokenRequest
	canonicalizerData := dsig.MakeC14N10RecCanonicalizer()                              // http://www.w3.org/TR/2001/REC-xml-c14n-20010315
	canonicalizerSignedProps := dsig.MakeC14N10ExclusiveCanonicalizerWithPrefixList("") // http://www.w3.org/2001/10/xml-exc-c14n#

	// Taken from example in library docs
	signContext := xades.SigningContext{
		DataContext: xades.SignedDataContext{
			Canonicalizer: canonicalizerData,
			Hash:          crypto.SHA256,
			ReferenceURI:  "",
			IsEnveloped:   true,
		},
		PropertiesContext: xades.SignedPropertiesContext{
			Canonicalizer: canonicalizerSignedProps,
			Hash:          crypto.SHA256,
		},
		Canonicalizer:     canonicalizerSignedInfo,
		Hash:              crypto.SHA256,
		KeyStore:          store,
		IssuerSerializer:  xades.IssuerSerializerKSeF,
		SigningTimeFormat: xades.SigningTimeFormatKSeF,
	}
	signature, err := xades.CreateSignature(root, &signContext)
	if err != nil {
		return nil, err
	}
	root.AddChild(signature)

	signedXML, err := doc.WriteToString()
	if err != nil {
		return nil, err
	}

	return []byte(signedXML), nil
}
