package api

import (
	"context"
	"time"
)

// AuthorisationChallengeResponse defines the authorization challenge response
type AuthorisationChallengeResponse struct {
	Timestamp time.Time `json:"timestamp"`
	Challenge string    `json:"challenge"`
}

// AuthorisationChallengeRequest defines the structure of the token session initialization
type AuthorisationChallengeRequest struct {
	ContextIdentifier *ContextIdentifier `json:"contextIdentifier"`
}

// ContextIdentifier defines the ContextIdentifier part of the authorization challenge
type ContextIdentifier struct {
	Identifier string `json:"identifier"`
	Type       string `json:"type"`
}

func fetchChallenge(ctx context.Context, c *Client) (*AuthorisationChallengeResponse, error) {
	response := &AuthorisationChallengeResponse{}

	request := &AuthorisationChallengeRequest{
		ContextIdentifier: &ContextIdentifier{
			Identifier: c.ID,
			Type:       "onip",
		},
	}

	resp, err := c.Client.R().
		SetResult(response).
		SetBody(request).
		SetContext(ctx).
		Post(c.URL + "/api/online/Session/AuthorisationChallenge")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, newErrorResponse(resp)
	}

	return response, nil
}
