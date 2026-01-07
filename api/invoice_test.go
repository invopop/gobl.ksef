package api_test

import (
	"context"
	"fmt"
	"testing"

	ksef_api "github.com/invopop/gobl.ksef/api"
	"github.com/invopop/gobl.ksef/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUploadInvoice(t *testing.T) {
	t.Run("uploads invoice during session", func(t *testing.T) {
		client := ksef_api.NewClient(
			&ksef_api.ContextIdentifier{Nip: "8126178616"},
			"./test/cert-20260102-131809.pfx",
			ksef_api.WithDebugClient(),
		)

		ctx := context.Background()
		err := client.Authenticate(ctx)
		require.NoError(t, err)

		uploadSession, err := ksef_api.CreateSession(ctx, client)
		require.NoError(t, err)

		doc, err := test.NewDocumentFrom("invoice-pl-pl.json")
		require.NoError(t, err)
		invoiceBytes, err := doc.Bytes()
		require.NoError(t, err)

		err = ksef_api.UploadInvoice(ctx, uploadSession, invoiceBytes, client)
		require.NoError(t, err)

		err = ksef_api.TerminateSession(uploadSession, ctx, client)
		assert.NoError(t, err)

		_, err = ksef_api.PollSessionStatus(ctx, uploadSession, client)
		assert.NoError(t, err)

		if err != nil {
			failedUploads, err := ksef_api.GetFailedUploadData(ctx, uploadSession, client)
			assert.NoError(t, err)
			for _, inv := range failedUploads.Invoices {
				fmt.Printf("Failed invoice %s (ordinal %d): %+v\n", inv.ReferenceNumber, inv.OrdinalNumber, inv.Status)
			}
		}
	})
}
