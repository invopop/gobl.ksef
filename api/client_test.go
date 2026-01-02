package api_test

import (
	"context"
	"testing"

	ksef_api "github.com/invopop/gobl.ksef/api"
	"github.com/stretchr/testify/assert"
)

func TestFetchSessionToken(t *testing.T) {
	t.Run("should get session token", func(t *testing.T) {
		client := ksef_api.NewClient(
			&ksef_api.ContextIdentifier{Nip: "8126178616"},
			"./test/cert-20260102-131809.pfx",
			ksef_api.WithDebugClient(),
		)
		// defer httpmock.DeactivateAndReset()
		ctx := context.Background()
		err := client.Authenticate(ctx)
		assert.NoError(t, err)

		// assert.Equal(t, client.SessionToken, "exampleSessionToken")
	})
}
