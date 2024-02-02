package api_test

import (
	"context"
	"testing"

	ksef_api "github.com/invopop/gobl.ksef/api"
	api_test "github.com/invopop/gobl.ksef/api/test"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestTerminateSession(t *testing.T) {
	t.Run("terminates the session", func(t *testing.T) {
		client, err := api_test.Client()
		defer httpmock.DeactivateAndReset()
		assert.NoError(t, err)

		httpmock.RegisterResponder("GET", "https://ksef-test.mf.gov.pl/api/online/Session/Terminate",
			httpmock.NewJsonResponderOrPanic(200, &ksef_api.TerminateSessionResponse{ProcessingCode: 200}))

		ctx := context.Background()
		terminateSessionResponse, err := ksef_api.TerminateSession(ctx, client)
		assert.NoError(t, err)

		assert.Equal(t, terminateSessionResponse.ProcessingCode, 200)
	})
}

func TestGetSessionStatus(t *testing.T) {
	t.Run("returns session status", func(t *testing.T) {
		client, err := api_test.Client()
		defer httpmock.DeactivateAndReset()
		assert.NoError(t, err)

		httpmock.RegisterResponder("GET", "https://ksef-test.mf.gov.pl/api/online/Session/Status",
			httpmock.NewJsonResponderOrPanic(200, &ksef_api.SessionStatusResponse{ProcessingCode: 200}))

		ctx := context.Background()
		sessionStatusResponse, err := ksef_api.GetSessionStatus(ctx, client)
		assert.NoError(t, err)

		assert.Equal(t, sessionStatusResponse.ProcessingCode, 200)
	})
}

func TestGetSessionStatusByReference(t *testing.T) {
	t.Run("returns session status", func(t *testing.T) {
		client, err := api_test.Client()
		defer httpmock.DeactivateAndReset()
		assert.NoError(t, err)

		httpmock.RegisterResponder("GET", "https://ksef-test.mf.gov.pl/api/common/Status/ExampleReferenceNumber",
			httpmock.NewJsonResponderOrPanic(200, &ksef_api.SessionStatusByReferenceResponse{ProcessingCode: 200}))

		ctx := context.Background()
		sessionStatusResponse, err := ksef_api.GetSessionStatusByReference(ctx, client)
		assert.NoError(t, err)

		assert.Equal(t, sessionStatusResponse.ProcessingCode, 200)
	})
}
