package ksef_api

import (
	"crypto/sha256"
	"encoding/base64"
	"os"
)

const (
	// Requests
	KSeFProductionBaseURL = "https://ksef.mf.gov.pl"
	KSeFDemoBaseURL       = "https://ksef-demo.mf.gov.pl"
	KSeFTestingBaseURL    = "https://ksef-test.mf.gov.pl"
	KsefProductionKeyPath = "../keys/prod.pem"
	KsefTestingKeyPath    = "../keys/test.pem"
	KsefDemoKeyPath       = "../keys/demo.pem"
)

type KSeFEnv struct {
	Url     string
	KeyPath string
}

func SendInvoices(env KSeFEnv, nip string, token string, invoices []string) (string, error) {
	client := NewClient(env.Url)
	session, err := client.NewSession(nip, token, env.KeyPath)
	if err != nil {
		return "", err
	}

	for _, invoice := range invoices {
		content, err := os.ReadFile(invoice)
		if err != nil {
			print(err)
			continue
		}
		digest := sha256.Sum256(content)
		digestBase64 := base64.StdEncoding.EncodeToString(digest[:])
		sendInvoiceResponse, err := session.SendInvoice(content, digestBase64)
		if err != nil {
			print(err)
			continue
		}

		invoiceStatusResponse, err := session.WaitUntilInvoiceIsProcessed(sendInvoiceResponse.ElementReferenceNumber)
		if err != nil {
			print(err)
			continue
		}
		generateQRCode(env, digestBase64, invoiceStatusResponse.InvoiceStatus.KsefReferenceNumber, invoice[:len(invoice)-4]+".png")
	}
	res, err := session.WaitUntilSessionIsTerminated()
	if err != nil {
		return "", err
	}
	upoBytes, err := base64.StdEncoding.DecodeString(res.Upo)
	if err != nil {
		return "", err
	}
	println(string(upoBytes))
	file, err := os.Create(res.ReferenceNumber + ".xml")
	if err != nil {
		return "", err
	}
	defer file.Close()
	_, err = file.Write(upoBytes)
	if err != nil {
		return "", err
	}

	return string(upoBytes), nil
}
