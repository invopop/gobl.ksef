package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasPolishSubjectPrefix(t *testing.T) {
	tests := []struct {
		name           string
		serialNumber   string
		expectPolish   bool
	}{
		{name: "TINPL prefix", serialNumber: "TINPL-1234567890", expectPolish: true},
		{name: "PNOPL prefix", serialNumber: "PNOPL-88102341294", expectPolish: true},
		{name: "PESEL prefix", serialNumber: "PESEL-88102341294", expectPolish: true},
		{name: "NIP prefix", serialNumber: "NIP-1234567890", expectPolish: true},
		{name: "Lithuanian PASLT", serialNumber: "PASLT-25428602", expectPolish: false},
		{name: "German TINDE", serialNumber: "TINDE-123456789", expectPolish: false},
		{name: "empty string", serialNumber: "", expectPolish: false},
		{name: "unrelated value", serialNumber: "CN=Test", expectPolish: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasPolishSubjectPrefix(tt.serialNumber)
			assert.Equal(t, tt.expectPolish, got)
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
			// nil certificate is treated as Polish (safe default)
			got := subjectIdentifierType(nil, tt.id)
			assert.Equal(t, tt.expected, got)
		})
	}
}
