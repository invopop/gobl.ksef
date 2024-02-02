package api_test

import (
	"testing"

	api_test "github.com/invopop/gobl.ksef/api/test"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestFetchSessionToken(t *testing.T) {
	t.Run("should get session token", func(t *testing.T) {
		client, err := api_test.Client()
		defer httpmock.DeactivateAndReset()
		assert.NoError(t, err)

		assert.Equal(t, client.SessionToken, "exampleSessionToken")
	})
}
