package api_test

import (
	"context"
	"testing"
	"time"

	ksef_api "github.com/invopop/gobl.ksef/api"
	"github.com/stretchr/testify/require"
)

func TestListInvoices(t *testing.T) {
	t.Run("lists uploaded invoices in the last 14 days", func(t *testing.T) {
		client := ksef_api.NewClient(
			&ksef_api.ContextIdentifier{Nip: "8126178616"},
			"./test/cert-20260102-131809.pfx",
			ksef_api.WithDebugClient(),
		)

		ctx := context.Background()
		require.NoError(t, client.Authenticate(ctx))

		today := time.Now().UTC()
		params := ksef_api.ListInvoicesParams{
			SubjectType: ksef_api.InvoiceSubjectTypeSupplier,
			From:        today.AddDate(0, 0, -14).Format(time.RFC3339),
			To:          today.Format(time.RFC3339),
		}

		_, err := client.ListInvoices(ctx, params)
		require.NoError(t, err)
	})
}
