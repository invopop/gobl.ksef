package api_test

import (
	"context"
	"testing"

	ksef_api "github.com/invopop/gobl.ksef/api"
	api_test "github.com/invopop/gobl.ksef/api/test"
	"github.com/stretchr/testify/assert"
)

func TestTerminateSession(t *testing.T) {
	t.Run("terminates the session", func(t *testing.T) {
		client, err := api_test.Client()
		// defer httpmock.DeactivateAndReset()
		assert.NoError(t, err)

		ctx := context.Background()
		err = ksef_api.TerminateSession("12345", ctx, client)
		assert.NoError(t, err)
	})
}
