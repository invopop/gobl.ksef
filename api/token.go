package api

import (
	"context"
	"fmt"
	"time"
)

func (t *ApiToken) isExpired(now time.Time) bool {
	if t == nil {
		return true
	}
	if t.ValidUntil == "" {
		return true
	}
	expiry, err := time.Parse(time.RFC3339Nano, t.ValidUntil)
	if err != nil {
		return true
	}
	return !now.Before(expiry)
}

// AccessTokenValue returns a valid access token, refreshing it when needed
// This token needs to be inserted in Authorization: Bearer <token> header
func (c *Client) AccessTokenValue(ctx context.Context) (string, error) {
	if c.accessToken != nil && !c.accessToken.isExpired(time.Now()) {
		return c.accessToken.Token, nil
	}

	if err := c.refreshAccessToken(ctx); err != nil {
		return "", err
	}

	if c.accessToken == nil {
		return "", fmt.Errorf("missing access token after refresh")
	}

	return c.accessToken.Token, nil
}

func (c *Client) refreshAccessToken(ctx context.Context) error {
	if c.refeshToken == nil {
		return fmt.Errorf("refresh token not available")
	}
	if c.refeshToken.isExpired(time.Now()) {
		// LATER: Re-authenticate
		return fmt.Errorf("refresh token expired")
	}

	response := &refreshAccessTokenResponse{}
	resp, err := c.client.R().
		SetHeader("Authorization", "Bearer "+c.refeshToken.Token).
		SetResult(response).
		SetContext(ctx).
		Post(c.url + "/auth/token/refresh")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return newErrorResponse(resp)
	}
	if response.AccessToken == nil {
		return fmt.Errorf("refresh response missing access token")
	}

	c.accessToken = response.AccessToken

	return nil
}

type refreshAccessTokenResponse struct {
	AccessToken *ApiToken `json:"accessToken"`
}
