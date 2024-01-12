package ksef_api

import (
	"crypto/sha256"
	"encoding/base64"
	"net/url"
	"os"
	"testing"
)

func TestLogin(t *testing.T) {
	env := KSeFEnv{
		Url:     KSeFTestingBaseURL,
		KeyPath: KsefTestingKeyPath,
	}
	// token to test env for organization with NIP 1234567788, not a security issue
	token := "624A48824F01935DADE66C83D4874C0EF7AF0529CB5F0F412E6932F189D3864A"

	client := NewClient(env.Url)
	println(KSeFTestingBaseURL)

	session, err := client.NewSession("1234567788", token, env.KeyPath)
	if err != nil {
		t.Fatal(err)
	}
	println(session.challenge)
	content, err := os.ReadFile("../../test/out/output.xml")
	if err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256(content)
	digestBase64 := base64.StdEncoding.EncodeToString(digest[:])
	sendInvoiceResponse, err := session.SendInvoice(content, digestBase64)
	if err != nil {
		t.Fatal(err)
	}

	invoiceStatusResponse, err := session.WaitUntilInvoiceIsProcessed(sendInvoiceResponse.ElementReferenceNumber)
	if err != nil {
		t.Fatal(err)
	}
	generateQRCode(env, url.QueryEscape(digestBase64), invoiceStatusResponse.InvoiceStatus.KsefReferenceNumber, "../../test/out/qr.png")
	res, err := session.WaitUntilSessionIsTerminated()
	if err != nil {
		t.Fatal(err)
	}
	str, err := base64.StdEncoding.DecodeString(res.Upo)
	if err != nil {
		t.Fatal(err)
	}
	println(string(str))

}
