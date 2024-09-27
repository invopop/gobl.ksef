package ksef_test

import (
	"testing"

	ksef "github.com/invopop/gobl.ksef"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/regimes/pl"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNewInv(t *testing.T) {
	t.Run("sets preceding invoice", func(t *testing.T) {
		inv := &bill.Invoice{
			Currency: currency.PLN,
			Supplier: &org.Party{
				TaxID: &tax.Identity{
					Country: l10n.PL.Tax(),
				},
			},
			Totals: &bill.Totals{
				Taxes: &tax.Total{},
			},
			Preceding: []*org.DocumentRef{
				{},
			},
		}

		invoice := ksef.NewInv(inv)

		assert.NotNil(t, invoice.CorrectedInv)
	})

	t.Run("sets correction reason", func(t *testing.T) {
		reason := "example reason"

		inv := &bill.Invoice{
			Currency: currency.PLN,
			Supplier: &org.Party{
				TaxID: &tax.Identity{
					Country: l10n.PL.Tax(),
				},
			},
			Totals: &bill.Totals{
				Taxes: &tax.Total{},
			},
			Preceding: []*org.DocumentRef{
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
			Currency: currency.PLN,
			Supplier: &org.Party{
				TaxID: &tax.Identity{
					Country: l10n.PL.Tax(),
				},
			},
			Totals: &bill.Totals{
				Taxes: &tax.Total{},
			},
			Preceding: []*org.DocumentRef{
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

	t.Run("sets the self-billing annotation to false in non-self-billed invoices", func(t *testing.T) {
		inv := &bill.Invoice{
			Currency: currency.PLN,
			Supplier: &org.Party{
				TaxID: &tax.Identity{
					Country: l10n.PL.Tax(),
				},
			},
			Totals: &bill.Totals{
				Taxes: &tax.Total{},
			},
		}

		invoice := ksef.NewInv(inv)

		assert.Equal(t, 2, invoice.Annotations.SelfBilling)
	})
	t.Run("sets the self-billing annotation to true in self-billed invoices", func(t *testing.T) {
		inv := &bill.Invoice{
			Currency: currency.PLN,
			Supplier: &org.Party{
				TaxID: &tax.Identity{
					Country: l10n.PL.Tax(),
				},
			},
			Tags: tax.WithTags(tax.TagSelfBilled),
			Totals: &bill.Totals{
				Taxes: &tax.Total{},
			},
		}

		invoice := ksef.NewInv(inv)

		assert.Equal(t, 1, invoice.Annotations.SelfBilling)
	})
}
