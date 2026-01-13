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
		invoiceBytes, err := doc.Bytes()
		require.NoError(t, err)

		err = uploadSession.UploadInvoice(ctx, invoiceBytes)
		require.NoError(t, err)

		err = uploadSession.FinishUpload(ctx)
		assert.NoError(t, err)

		_, err = uploadSession.PollStatus(ctx)
		assert.NoError(t, err)

		failedUploads, err := uploadSession.GetFailedUploadData(ctx)
		assert.NoError(t, err)
		for _, inv := range failedUploads {
			fmt.Printf("Failed invoice %s (ordinal %d): %+v\n", inv.ReferenceNumber, inv.OrdinalNumber, inv.Status)
		}
	})
}
