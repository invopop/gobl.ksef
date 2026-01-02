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
func (c *Client) AccessTokenValue(ctx context.Context) (string, error) {
	if c.AccessToken != nil && !c.AccessToken.isExpired(time.Now()) {
		return c.AccessToken.Token, nil
	}

	if err := c.refreshAccessToken(ctx); err != nil {
		return "", err
	}

	if c.AccessToken == nil {
		return "", fmt.Errorf("missing access token after refresh")
	}

	return c.AccessToken.Token, nil
}

func (c *Client) refreshAccessToken(ctx context.Context) error {
	if c.RefreshToken == nil {
		return fmt.Errorf("refresh token not available")
	}
	if c.RefreshToken.isExpired(time.Now()) {
		// LATER: Re-authenticate
		return fmt.Errorf("refresh token expired")
	}

	response := &refreshAccessTokenResponse{}
	resp, err := c.Client.R().
		SetHeader("Authorization", "Bearer "+c.RefreshToken.Token).
		SetResult(response).
		SetContext(ctx).
		Post(c.URL + "/v2/auth/token/refresh")
	if err != nil {
		return err
	}
	if resp.IsError() {
		return newErrorResponse(resp)
	}
	if response.AccessToken == nil {
		return fmt.Errorf("refresh response missing access token")
	}

	c.AccessToken = response.AccessToken

	return nil
}

type refreshAccessTokenResponse struct {
	AccessToken *ApiToken `json:"accessToken"`
}
