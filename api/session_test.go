package api_test

import (
	"context"
	"testing"

	ksef_api "github.com/invopop/gobl.ksef/api"
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

		resp, err := ksef_api.CreateSession(ctx, client)
		require.NoError(t, err)
		assert.NotEmpty(t, resp.ReferenceNumber)
		assert.NotEmpty(t, resp.ValidUntil)
	})
}
