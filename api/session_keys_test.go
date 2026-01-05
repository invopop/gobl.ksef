package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSelectSymmetricKeyCertificate(t *testing.T) {
	now := time.Date(2024, time.July, 11, 12, 0, 0, 0, time.UTC)

	certificates := []PublicKeyCertificate{
		{
			Certificate: "inactive",
			ValidFrom:   now.Add(-24 * time.Hour),
			ValidTo:     now.Add(-time.Hour),
			Usage:       []string{symmetricKeyUsage},
		},
		{
			Certificate: "token-only",
			ValidFrom:   now.Add(-time.Hour),
			ValidTo:     now.Add(time.Hour),
			Usage:       []string{"KsefTokenEncryption"},
		},
		{
			Certificate: "valid",
			ValidFrom:   now.Add(-time.Hour),
			ValidTo:     now.Add(time.Hour),
			Usage: []string{
				"KsefTokenEncryption",
				symmetricKeyUsage,
			},
		},
		{
			Certificate: "future",
			ValidFrom:   now.Add(time.Hour),
			ValidTo:     now.Add(2 * time.Hour),
			Usage:       []string{symmetricKeyUsage},
		},
	}

	cert, err := selectSymmetricKeyCertificate(certificates, now)
	require.NoError(t, err)
	assert.Equal(t, "valid", cert.Certificate)
}

func TestSelectSymmetricKeyCertificateNoMatch(t *testing.T) {
	now := time.Date(2024, time.July, 11, 12, 0, 0, 0, time.UTC)

	certificates := []PublicKeyCertificate{
		{
			Certificate: "missing-usage",
			ValidFrom:   now.Add(-time.Hour),
			ValidTo:     now.Add(time.Hour),
			Usage:       []string{"KsefTokenEncryption"},
		},
	}

	_, err := selectSymmetricKeyCertificate(certificates, now)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no suitable RSA public key found")
}
