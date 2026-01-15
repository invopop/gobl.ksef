package api

import (
	"context"
	"errors"
	"time"
)

var ErrNoSuitableRSAPublicKey = errors.New("no suitable RSA public key found")

// payload returned from /security/public-key-certificates
type publicKeyCertificate struct {
	Certificate string    `json:"certificate"`
	ValidFrom   time.Time `json:"validFrom"`
	ValidTo     time.Time `json:"validTo"`
	Usage       []string  `json:"usage"`
}

const symmetricKeyUsage = "SymmetricKeyEncryption"

// returns the RSA public key used to encrypt the per-session symmetric key.
func (c *Client) getRSAPublicKey(ctx context.Context) (*publicKeyCertificate, error) {
	var certificates []publicKeyCertificate
	resp, err := c.client.R().
		SetContext(ctx).
		SetResult(&certificates).
		Get(c.url + "/security/public-key-certificates")
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, newErrorResponse(resp)
	}

	return selectSymmetricKeyCertificate(certificates, time.Now().UTC())
}

func selectSymmetricKeyCertificate(certificates []publicKeyCertificate, now time.Time) (*publicKeyCertificate, error) {
	for i := range certificates {
		cert := &certificates[i]
		if cert.ValidFrom.After(now) {
			continue
		}
		if !cert.ValidTo.After(now) {
			continue
		}
		if !containsUsage(cert.Usage, symmetricKeyUsage) {
			continue
		}
		return cert, nil
	}

	return nil, ErrNoSuitableRSAPublicKey
}

func containsUsage(usages []string, desired string) bool {
	for _, usage := range usages {
		if usage == desired {
			return true
		}
	}
	return false
}
