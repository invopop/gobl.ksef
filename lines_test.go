package ksef_test

import (
	"testing"

	ksef "github.com/invopop/gobl.ksef"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
)

func TestNewLines(t *testing.T) {
	t.Run("calculates unitDiscount", func(t *testing.T) {
		lines := []*bill.Line{
			{
				Item: &org.Item{},
				Taxes: tax.Set{
					{
						Percent: &num.Percentage{},
					},
				},
				Quantity: num.MakeAmount(1, 0),
				Discounts: []*bill.LineDiscount{
					{
						Amount: num.MakeAmount(10000, 2),
					},
				},
			},
		}

		ln := ksef.NewLines(lines)

		assert.Equal(t, "100.00", ln[0].UnitDiscount)
	})

	t.Run("calculates unitDiscount per unit", func(t *testing.T) {
		lines := []*bill.Line{
			{
				Item: &org.Item{},
				Taxes: tax.Set{
					{
						Percent: &num.Percentage{},
					},
				},
				Quantity: num.MakeAmount(2, 0),
				Discounts: []*bill.LineDiscount{
					{
						Amount: num.MakeAmount(10000, 2),
					},
				},
			},
		}

		ln := ksef.NewLines(lines)

		assert.Equal(t, "50.00", ln[0].UnitDiscount)
	})

	t.Run("returns empty string for unitDiscount if no discount present", func(t *testing.T) {
		lines := []*bill.Line{
			{
				Item: &org.Item{},
				Taxes: tax.Set{
					{
						Percent: &num.Percentage{},
					},
				},
				Quantity:  num.MakeAmount(1, 0),
				Discounts: []*bill.LineDiscount{},
			},
		}

		ln := ksef.NewLines(lines)

		assert.Equal(t, "", ln[0].UnitDiscount)
	})

	t.Run("unitDiscount adds up multiple discounts", func(t *testing.T) {
		lines := []*bill.Line{
			{
				Item: &org.Item{},
				Taxes: tax.Set{
					{
						Percent: &num.Percentage{},
					},
				},
				Quantity: num.MakeAmount(1, 0),
				Discounts: []*bill.LineDiscount{
					{
						Amount: num.MakeAmount(10000, 2),
					},
					{
						Amount: num.MakeAmount(10000, 2),
					},
				},
			},
		}

		ln := ksef.NewLines(lines)

		assert.Equal(t, "200.00", ln[0].UnitDiscount)
	})
}
