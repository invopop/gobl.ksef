package ksef_test

import (
	"testing"

	ksef "github.com/invopop/gobl.ksef"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/pl"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNewInv(t *testing.T) {
	t.Run("sets preceding invoice", func(t *testing.T) {
		inv := &bill.Invoice{
			Supplier: &org.Party{
				TaxID: &tax.Identity{
					Country: l10n.PL,
				},
			},
			Totals: &bill.Totals{
				Taxes: &tax.Total{},
			},
			Preceding: []*bill.Preceding{
				{},
			},
		}

		invoice := ksef.NewInv(inv)

		assert.NotNil(t, invoice.CorrectedInv)
	})

	t.Run("sets correction reason", func(t *testing.T) {
		reason := "example reason"

		inv := &bill.Invoice{
			Supplier: &org.Party{
				TaxID: &tax.Identity{
					Country: l10n.PL,
				},
			},
			Totals: &bill.Totals{
				Taxes: &tax.Total{},
			},
			Preceding: []*bill.Preceding{
				{
					Reason: reason,
				},
			},
		}

		invoice := ksef.NewInv(inv)

		assert.Equal(t, reason, invoice.CorrectionReason)
	})

	t.Run("sets correction type", func(t *testing.T) {
		inv := &bill.Invoice{
			Supplier: &org.Party{
				TaxID: &tax.Identity{
					Country: l10n.PL,
				},
			},
			Totals: &bill.Totals{
				Taxes: &tax.Total{},
			},
			Preceding: []*bill.Preceding{
				{
					Ext: tax.Extensions{
						pl.ExtKeyKSeFEffectiveDate: "1",
					},
				},
			},
		}

		invoice := ksef.NewInv(inv)

		assert.Equal(t, "1", invoice.CorrectionType)
	})
}
