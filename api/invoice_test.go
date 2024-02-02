package api_test

import (
	"context"
	"os"
	"testing"

	ksef_api "github.com/invopop/gobl.ksef/api"
	api_test "github.com/invopop/gobl.ksef/api/test"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestSendInvoice(t *testing.T) {
	t.Run("should post invoice", func(t *testing.T) {
		client, err := api_test.Client()
		defer httpmock.DeactivateAndReset()
		assert.NoError(t, err)

		elementReferenceNumber := "ExampleReferenceNumber"
		httpmock.RegisterResponder("PUT", "https://ksef-test.mf.gov.pl/api/online/Invoice/Send",
			httpmock.NewJsonResponderOrPanic(200, &ksef_api.SendInvoiceResponse{ElementReferenceNumber: elementReferenceNumber}))

		content, err := os.ReadFile("../test/data/out/invoice-pl-pl.xml")
		assert.NoError(t, err)

		ctx := context.Background()
		sendInvoiceResponse, err := ksef_api.SendInvoice(ctx, client, content)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, sendInvoiceResponse.ElementReferenceNumber, elementReferenceNumber)
	})
}

func TestGetInvoiceStatus(t *testing.T) {
	t.Run("should get invoice status", func(t *testing.T) {
		client, err := api_test.Client()
		defer httpmock.DeactivateAndReset()
		assert.NoError(t, err)

		httpmock.RegisterResponder("GET", "https://ksef-test.mf.gov.pl/api/online/Invoice/Status/exampleReferenceNumber",
			httpmock.NewJsonResponderOrPanic(200, &ksef_api.InvoiceStatusResponse{ProcessingCode: 200}))

		ctx := context.Background()
		invoiceStatusResponse, err := ksef_api.FetchInvoiceStatus(ctx, client, "exampleReferenceNumber")
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, invoiceStatusResponse.ProcessingCode, 200)
	})
}
