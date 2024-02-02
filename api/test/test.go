// Package api_test provides tools for testing the api
package api_test

import (
	"context"
	"time"

	"github.com/go-resty/resty/v2"
	ksef_api "github.com/invopop/gobl.ksef/api"
	"github.com/jarcoal/httpmock"
)

// Client creates authorized client for testing
func Client() (*ksef_api.Client, error) {
	mockClient := resty.New()

	httpmock.ActivateNonDefault(mockClient.GetClient())

	reqT, err := time.Parse("2006-01-02T15:04:05.000Z", "2024-01-26T16:18:51.701Z")
	if err != nil {
		return nil, err
	}

	httpmock.RegisterResponder("POST", "https://ksef-test.mf.gov.pl/api/online/Session/AuthorisationChallenge",
		httpmock.NewJsonResponderOrPanic(200, &ksef_api.AuthorisationChallengeResponse{Timestamp: reqT, Challenge: "20240126-CR-077CAFEC31-83ACAC25E4-64"}))

	sessionToken := "exampleSessionToken"
	httpmock.RegisterResponder("POST", "https://ksef-test.mf.gov.pl/api/online/Session/InitToken",
		httpmock.NewJsonResponderOrPanic(200, &ksef_api.InitSessionTokenResponse{ReferenceNumber: "ExampleReferenceNumber", SessionToken: &ksef_api.SessionToken{Token: sessionToken}}))

	client := ksef_api.NewClient(
		ksef_api.WithClient(mockClient),
		ksef_api.WithID("1234567788"),
		ksef_api.WithToken("624A48824F01935DADE66C83D4874C0EF7AF0529CB5F0F412E6932F189D3864A"),
		ksef_api.WithKeyPath("./keys/test.pem"),
	)

	ctx := context.Background()
	err = ksef_api.FetchSessionToken(ctx, client)
	if err != nil {
		return nil, err
	}

	return client, nil
}
