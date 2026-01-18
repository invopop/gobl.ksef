package api_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	ksef_api "github.com/invopop/gobl.ksef/api"
	"github.com/invopop/gobl.ksef/test"
	"github.com/invopop/gobl/regimes/pl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateSession(t *testing.T) {
	t.Run("creates session", func(t *testing.T) {
		client := ksef_api.NewClient(
			&ksef_api.ContextIdentifier{Nip: "8126178616"},
			"./test/cert-20260102-131809.pfx",
			ksef_api.WithDebugClient(),
		)

		ctx := context.Background()
		err := client.Authenticate(ctx)
		require.NoError(t, err)

		uploadSession, err := client.CreateSession(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, uploadSession.ReferenceNumber)
		assert.NotEmpty(t, uploadSession.ValidUntil)
		assert.Len(t, uploadSession.SymmetricKey, 32)
		assert.Len(t, uploadSession.InitializationVector, 16)

		err = uploadSession.FinishUpload(ctx)
		assert.NoError(t, err)
	})
}

func TestUploadInvoice(t *testing.T) {
	t.Run("uploads invoice during session", func(t *testing.T) {
		fmt.Println(1)
		client := ksef_api.NewClient(
			&ksef_api.ContextIdentifier{Nip: "8126178616"},
			"./test/cert-20260102-131809.pfx",
			ksef_api.WithDebugClient(),
		)

		ctx := context.Background()
		err := client.Authenticate(ctx)
		require.NoError(t, err)

		uploadSession, err := client.CreateSession(ctx)
		require.NoError(t, err)

		doc, err := test.NewDocumentFrom("invoice-pl-pl.json")
		require.NoError(t, err)

		// Generate unique identifier for the invoice.
		// Without it, uploading will result in error because of a duplicate.
		now := time.Now().UTC()
		doc.Inv.IssueDate = now.Format("2006-01-02")             // current date
		doc.Inv.SequentialNumber = fmt.Sprintf("%d", now.Unix()) // Unix timestamp in seconds

		invoiceBytes, err := doc.Bytes()
		require.NoError(t, err)

		err = uploadSession.UploadInvoice(ctx, invoiceBytes)
		require.NoError(t, err)

		err = uploadSession.FinishUpload(ctx)
		assert.NoError(t, err)

		_, err = uploadSession.PollStatus(ctx)
		assert.NoError(t, err)

		uploadedInvoices, err := uploadSession.ListUploadedInvoices(ctx)
		assert.NoError(t, err)
		assert.Len(t, uploadedInvoices, 1)

		// For debugging - we should not get any failed uploads, but if an upload fails, we should get more information about what exactly went wrong
		failedUploads, err := uploadSession.GetFailedUploadData(ctx)
		assert.NoError(t, err)
		for _, inv := range failedUploads {
			fmt.Printf("Failed invoice %s (ordinal %d): %+v\n", inv.ReferenceNumber, inv.OrdinalNumber, inv.Status)
		}

		envelope, err := test.LoadTestEnvelope("invoice-pl-pl.json")
		require.NoError(t, err)

		// Create QR code for the uploaded invoice and check if it's properly generated
		err = client.Sign(envelope, "8126178616", &(uploadedInvoices[0]))
		assert.NoError(t, err)

		require.NotNil(t, envelope.Head.Stamps)
		assert.GreaterOrEqual(t, len(envelope.Head.Stamps), 3)

		var qrStampValue string
		for _, stamp := range envelope.Head.Stamps {
			if stamp != nil && stamp.Provider == pl.StampProviderKSeFQR {
				qrStampValue = stamp.Value
				break
			}
		}
		require.NotEmpty(t, qrStampValue)

		// Check if the URL is correctly formed
		// IMPORTANT: when the URL contains invalid parameters (e.g. NIP is different), the response is still 200,
		// but the website content says that "no invoice found".
		// To check if the URL is actually valid, we need to check the returned HTML, and this is very fragile
		resp, err := http.Get(qrStampValue)
		require.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
