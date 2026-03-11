package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsNonPolishSubjectSerialNumber(t *testing.T) {
	tests := []struct {
		name          string
		serialNumber  string
		expectForeign bool
	}{
		// Polish prefixes → not foreign
		{name: "TINPL prefix", serialNumber: "TINPL-1234567890", expectForeign: false},
		{name: "PNOPL prefix", serialNumber: "PNOPL-88102341294", expectForeign: false},
		{name: "PESEL prefix", serialNumber: "PESEL-88102341294", expectForeign: false},
		{name: "NIP prefix", serialNumber: "NIP-1234567890", expectForeign: false},
		// Foreign prefixes → foreign
		{name: "Lithuanian PASLT", serialNumber: "PASLT-25428602", expectForeign: true},
		{name: "German TINDE", serialNumber: "TINDE-123456789", expectForeign: true},
		{name: "unrelated value", serialNumber: "CN=Test", expectForeign: true},
		// Empty → not foreign (safe default for Polish CCK KSeF certs
		// that carry the NIP in OID 2.5.4.97 instead of SerialNumber)
		{name: "empty string", serialNumber: "", expectForeign: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isNonPolishSubjectSerialNumber(tt.serialNumber)
			assert.Equal(t, tt.expectForeign, got)
		})
	}
}

func TestSubjectIdentifierType(t *testing.T) {
	tests := []struct {
		name     string
		id       *ContextIdentifier
		expected string
	}{
		{
			name:     "nil cert with Nip defaults to certificateSubject",
			id:       &ContextIdentifier{Nip: "5213967846"},
			expected: "certificateSubject",
		},
		{
			name:     "NipVatUe forces certificateFingerprint",
			id:       &ContextIdentifier{NipVatUe: "PL5213967846"},
			expected: "certificateFingerprint",
		},
		{
			name:     "nil identifier defaults to certificateSubject",
			id:       nil,
			expected: "certificateSubject",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// nil certificate is treated as non-foreign (safe default)
			got := subjectIdentifierType(nil, tt.id)
			assert.Equal(t, tt.expected, got)
		})
	}
}
