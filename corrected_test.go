package ksef_test

import (
	"testing"
	"time"

	ksef "github.com/invopop/gobl.ksef"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/org"
	"github.com/stretchr/testify/assert"
)

func TestNewCorrectedInv(t *testing.T) {
	t.Run("sets invoice number", func(t *testing.T) {
		prc := &org.DocumentRef{
			Series: "SAMPLE",
			Code:   "001",
		}

		cor := ksef.NewCorrectedInv(prc)

		assert.Equal(t, "SAMPLE-001", cor.SequentialNumber)
	})

	t.Run("sets issue date", func(t *testing.T) {
		prc := &org.DocumentRef{
			IssueDate: cal.NewDate(2024, time.March, 14),
		}

		cor := ksef.NewCorrectedInv(prc)

		assert.Equal(t, "2024-03-14", cor.IssueDate)
	})

	t.Run("sets no ksef number flag", func(t *testing.T) {
		prc := &org.DocumentRef{}

		cor := ksef.NewCorrectedInv(prc)

		assert.Equal(t, 1, cor.NoKsefNumberPresent)
	})

	t.Run("sets ksef number", func(t *testing.T) {
		ksefID := "123"
		prc := &org.DocumentRef{
			Stamps: []*head.Stamp{
				{
					Provider: "ksef-id",
					Value:    ksefID,
				},
			},
		}

		cor := ksef.NewCorrectedInv(prc)

		assert.Equal(t, 1, cor.KsefNumberPresent)
		assert.Equal(t, ksefID, cor.KsefNumber)
	})
}
