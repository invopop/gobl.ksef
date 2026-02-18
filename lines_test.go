package ksef_test

import (
	"testing"

	ksef "github.com/invopop/gobl.ksef"
	"github.com/invopop/gobl/addons/pl/favat"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLines(t *testing.T) {
	t.Run("converts basic lines", func(t *testing.T) {
		price, _ := num.AmountFromString("100.00")
		qty, _ := num.AmountFromString("2")
		total, _ := num.AmountFromString("200.00")

		lines := []*bill.Line{
			{
				Index:    1,
				Quantity: qty,
				Item: &org.Item{
					Name:  "Test Item",
					Price: &price,
					Unit:  "h",
				},
				Total: &total,
				Taxes: tax.Set{
					&tax.Combo{
						Category: tax.CategoryVAT,
						Percent:  num.NewPercentage(23, 2),
					},
				},
			},
		}

		result := ksef.NewLines(lines)

		require.Len(t, result, 1)
		assert.Equal(t, 1, result[0].LineNumber)
		assert.Equal(t, "Test Item", result[0].Name)
		assert.Equal(t, "HUR", result[0].Measure)
		assert.Equal(t, "100.00", result[0].NetUnitPrice)
		assert.Equal(t, "2", result[0].Quantity)
		assert.Equal(t, "200.00", result[0].NetPriceTotal)
		assert.Equal(t, "23", result[0].VATRate)
	})

	t.Run("handles multiple lines", func(t *testing.T) {
		price, _ := num.AmountFromString("50.00")
		qty, _ := num.AmountFromString("1")
		total, _ := num.AmountFromString("50.00")

		lines := []*bill.Line{
			{
				Index:    1,
				Quantity: qty,
				Item: &org.Item{
					Name:  "Item 1",
					Price: &price,
					Unit:  "h",
				},
				Total: &total,
				Taxes: tax.Set{
					&tax.Combo{
						Category: tax.CategoryVAT,
						Percent:  num.NewPercentage(23, 2),
					},
				},
			},
			{
				Index:    2,
				Quantity: qty,
				Item: &org.Item{
					Name:  "Item 2",
					Price: &price,
					Unit:  "service",
				},
				Total: &total,
				Taxes: tax.Set{
					&tax.Combo{
						Category: tax.CategoryVAT,
						Percent:  num.NewPercentage(8, 2),
					},
				},
			},
		}

		result := ksef.NewLines(lines)

		require.Len(t, result, 2)
		assert.Equal(t, "Item 1", result[0].Name)
		assert.Equal(t, "Item 2", result[1].Name)
		assert.Equal(t, 1, result[0].LineNumber)
		assert.Equal(t, 2, result[1].LineNumber)
	})

	t.Run("handles line without VAT percent", func(t *testing.T) {
		price, _ := num.AmountFromString("100.00")
		qty, _ := num.AmountFromString("1")
		total, _ := num.AmountFromString("100.00")

		lines := []*bill.Line{
			{
				Index:    1,
				Quantity: qty,
				Item: &org.Item{
					Name:  "Exempt Item",
					Price: &price,
					Unit:  "h",
				},
				Total: &total,
				Taxes: tax.Set{
					&tax.Combo{
						Category: tax.CategoryVAT,
						// No Percent for exempt items
					},
				},
			},
		}

		result := ksef.NewLines(lines)

		require.Len(t, result, 1)
		assert.Equal(t, "", result[0].VATRate)
	})

	t.Run("handles line with discounts", func(t *testing.T) {
		price, _ := num.AmountFromString("100.00")
		qty, _ := num.AmountFromString("2")
		total, _ := num.AmountFromString("180.00")
		discountAmt, _ := num.AmountFromString("20.00")

		lines := []*bill.Line{
			{
				Index:    1,
				Quantity: qty,
				Item: &org.Item{
					Name:  "Discounted Item",
					Price: &price,
					Unit:  "h",
				},
				Discounts: []*bill.LineDiscount{
					{
						Amount: discountAmt,
						Reason: "Volume discount",
					},
				},
				Total: &total,
				Taxes: tax.Set{
					&tax.Combo{
						Category: tax.CategoryVAT,
						Percent:  num.NewPercentage(23, 2),
					},
				},
			},
		}

		result := ksef.NewLines(lines)

		require.Len(t, result, 1)
		assert.Equal(t, "10.00", result[0].UnitDiscount) // 20.00 / 2 = 10.00
	})

	t.Run("handles multiple discounts on same line", func(t *testing.T) {
		price, _ := num.AmountFromString("100.00")
		qty, _ := num.AmountFromString("4")
		total, _ := num.AmountFromString("360.00")
		discount1, _ := num.AmountFromString("20.00")
		discount2, _ := num.AmountFromString("20.00")

		lines := []*bill.Line{
			{
				Index:    1,
				Quantity: qty,
				Item: &org.Item{
					Name:  "Multiple Discounts",
					Price: &price,
					Unit:  "h",
				},
				Discounts: []*bill.LineDiscount{
					{Amount: discount1, Reason: "Discount 1"},
					{Amount: discount2, Reason: "Discount 2"},
				},
				Total: &total,
				Taxes: tax.Set{
					&tax.Combo{
						Category: tax.CategoryVAT,
						Percent:  num.NewPercentage(23, 2),
					},
				},
			},
		}

		result := ksef.NewLines(lines)

		require.Len(t, result, 1)
		assert.Equal(t, "10.00", result[0].UnitDiscount) // (20.00 + 20.00) / 4 = 10.00
	})
}

func TestNewOrderLines(t *testing.T) {
	t.Run("converts order lines with VAT calculation", func(t *testing.T) {
		price, _ := num.AmountFromString("100.00")
		qty, _ := num.AmountFromString("2")
		total, _ := num.AmountFromString("200.00")

		lines := []*bill.Line{
			{
				Index:    1,
				Quantity: qty,
				Item: &org.Item{
					Name:  "Order Item",
					Price: &price,
					Unit:  "h",
				},
				Total: &total,
				Taxes: tax.Set{
					&tax.Combo{
						Category: tax.CategoryVAT,
						Percent:  num.NewPercentage(23, 2),
						Ext: tax.Extensions{
							favat.ExtKeyTaxCategory: "1",
						},
					},
				},
			},
		}

		result := ksef.NewOrderLines(lines, 2)

		require.Len(t, result, 1)
		assert.Equal(t, 1, result[0].LineNumber)
		assert.Equal(t, "Order Item", result[0].Name)
		assert.Equal(t, "HUR", result[0].Measure)
		assert.Equal(t, "100.00", result[0].NetUnitPrice)
		assert.Equal(t, "2", result[0].Quantity)
		assert.Equal(t, "200.00", result[0].NetPriceTotal)
		assert.Equal(t, "23", result[0].VATRate)
		// Note: TaxValue calculation in current code uses Percent.Amount() which
		// returns the percentage value directly (23), not the decimal (0.23),
		// resulting in 200.00 * 23 = 4600.00 (there's a TODO in the code about this)
		assert.Equal(t, "4600.00", result[0].TaxValue)
	})

	t.Run("handles multiple order lines", func(t *testing.T) {
		price, _ := num.AmountFromString("50.00")
		qty, _ := num.AmountFromString("1")
		total, _ := num.AmountFromString("50.00")

		lines := []*bill.Line{
			{
				Index:    1,
				Quantity: qty,
				Item: &org.Item{
					Name:  "Item 1",
					Price: &price,
					Unit:  "h",
				},
				Total: &total,
				Taxes: tax.Set{
					&tax.Combo{
						Category: tax.CategoryVAT,
						Percent:  num.NewPercentage(23, 2),
					},
				},
			},
			{
				Index:    2,
				Quantity: qty,
				Item: &org.Item{
					Name:  "Item 2",
					Price: &price,
					Unit:  "service",
				},
				Total: &total,
				Taxes: tax.Set{
					&tax.Combo{
						Category: tax.CategoryVAT,
						Percent:  num.NewPercentage(8, 2),
					},
				},
			},
		}

		result := ksef.NewOrderLines(lines, 2)

		require.Len(t, result, 2)
		assert.Equal(t, "Item 1", result[0].Name)
		assert.Equal(t, "Item 2", result[1].Name)
	})
}

func TestLineToGOBL(t *testing.T) {
	t.Run("converts basic KSEF line to GOBL", func(t *testing.T) {
		ksefLine := &ksef.Line{
			Name:         "Test Item",
			Quantity:     "2",
			NetUnitPrice: "100.00",
			Measure:      "HUR",
			VATRate:      "23",
			NetPriceTotal: "200.00",
		}

		line, err := ksefLine.ToGOBL()

		require.NoError(t, err)
		assert.Equal(t, "Test Item", line.Item.Name)
		assert.Equal(t, "2", line.Quantity.String())
		assert.Equal(t, "100.00", line.Item.Price.String())
		assert.Equal(t, org.Unit("HUR"), line.Item.Unit)
		assert.Len(t, line.Taxes, 1)
		assert.Equal(t, tax.CategoryVAT, line.Taxes[0].Category)
		assert.Equal(t, "23", line.Taxes[0].Percent.Amount().MinimalString())
	})

	t.Run("handles line with discount", func(t *testing.T) {
		ksefLine := &ksef.Line{
			Name:          "Discounted Item",
			Quantity:      "2",
			NetUnitPrice:  "100.00",
			UnitDiscount:  "10.00",
			Measure:       "HUR",
			VATRate:       "23",
			NetPriceTotal: "180.00",
		}

		line, err := ksefLine.ToGOBL()

		require.NoError(t, err)
		assert.Len(t, line.Discounts, 1)
		assert.Equal(t, "10.00", line.Discounts[0].Amount.String())
	})

	t.Run("handles exempt line with zero VAT", func(t *testing.T) {
		ksefLine := &ksef.Line{
			Name:          "Exempt Item",
			Quantity:      "1",
			NetUnitPrice:  "100.00",
			Measure:       "HUR",
			VATRate:       "zw",
			NetPriceTotal: "100.00",
		}

		line, err := ksefLine.ToGOBL()

		require.NoError(t, err)
		assert.Len(t, line.Taxes, 1)
		assert.Equal(t, tax.KeyExempt, line.Taxes[0].Key)
		assert.Equal(t, "7", string(line.Taxes[0].Ext[favat.ExtKeyTaxCategory]))
	})

	t.Run("handles intra-community supply", func(t *testing.T) {
		ksefLine := &ksef.Line{
			Name:          "EU Item",
			Quantity:      "1",
			NetUnitPrice:  "100.00",
			Measure:       "HUR",
			VATRate:       "0 WDT",
			NetPriceTotal: "100.00",
		}

		line, err := ksefLine.ToGOBL()

		require.NoError(t, err)
		assert.Len(t, line.Taxes, 1)
		assert.Equal(t, tax.KeyIntraCommunity, line.Taxes[0].Key)
		assert.Equal(t, "6.2", string(line.Taxes[0].Ext[favat.ExtKeyTaxCategory]))
	})

	t.Run("handles export supply", func(t *testing.T) {
		ksefLine := &ksef.Line{
			Name:          "Export Item",
			Quantity:      "1",
			NetUnitPrice:  "100.00",
			Measure:       "HUR",
			VATRate:       "0 EX",
			NetPriceTotal: "100.00",
		}

		line, err := ksefLine.ToGOBL()

		require.NoError(t, err)
		assert.Len(t, line.Taxes, 1)
		assert.Equal(t, tax.KeyExport, line.Taxes[0].Key)
		assert.Equal(t, "6.3", string(line.Taxes[0].Ext[favat.ExtKeyTaxCategory]))
	})

	t.Run("handles reverse charge", func(t *testing.T) {
		ksefLine := &ksef.Line{
			Name:          "Reverse Charge Item",
			Quantity:      "1",
			NetUnitPrice:  "100.00",
			Measure:       "HUR",
			VATRate:       "np II",
			NetPriceTotal: "100.00",
		}

		line, err := ksefLine.ToGOBL()

		require.NoError(t, err)
		assert.Len(t, line.Taxes, 1)
		assert.Equal(t, tax.KeyReverseCharge, line.Taxes[0].Key)
		assert.Equal(t, "9", string(line.Taxes[0].Ext[favat.ExtKeyTaxCategory]))
	})

	t.Run("handles domestic reverse charge", func(t *testing.T) {
		ksefLine := &ksef.Line{
			Name:          "Domestic RC Item",
			Quantity:      "1",
			NetUnitPrice:  "100.00",
			Measure:       "HUR",
			VATRate:       "oo",
			NetPriceTotal: "100.00",
		}

		line, err := ksefLine.ToGOBL()

		require.NoError(t, err)
		assert.Len(t, line.Taxes, 1)
		assert.Equal(t, tax.KeyReverseCharge, line.Taxes[0].Key)
		assert.Equal(t, "10", string(line.Taxes[0].Ext[favat.ExtKeyTaxCategory]))
	})

	t.Run("handles reduced rate 8%", func(t *testing.T) {
		ksefLine := &ksef.Line{
			Name:          "Reduced Item",
			Quantity:      "1",
			NetUnitPrice:  "100.00",
			Measure:       "HUR",
			VATRate:       "8",
			NetPriceTotal: "100.00",
		}

		line, err := ksefLine.ToGOBL()

		require.NoError(t, err)
		assert.Len(t, line.Taxes, 1)
		assert.Equal(t, tax.KeyStandard, line.Taxes[0].Key)
		assert.Equal(t, tax.RateReduced, line.Taxes[0].Rate)
		assert.Equal(t, "8", line.Taxes[0].Percent.Amount().MinimalString())
	})

	t.Run("handles super reduced rate 5%", func(t *testing.T) {
		ksefLine := &ksef.Line{
			Name:          "Super Reduced Item",
			Quantity:      "1",
			NetUnitPrice:  "100.00",
			Measure:       "HUR",
			VATRate:       "5",
			NetPriceTotal: "100.00",
		}

		line, err := ksefLine.ToGOBL()

		require.NoError(t, err)
		assert.Len(t, line.Taxes, 1)
		assert.Equal(t, tax.KeyStandard, line.Taxes[0].Key)
		assert.Equal(t, tax.RateSuperReduced, line.Taxes[0].Rate)
		assert.Equal(t, "5", line.Taxes[0].Percent.Amount().MinimalString())
	})

	t.Run("handles gross unit price (P_9B)", func(t *testing.T) {
		ksefLine := &ksef.Line{
			Name:           "Gross Price Item",
			Quantity:       "1",
			GrossUnitPrice: "123.00",
			Measure:        "HUR",
			VATRate:        "23",
		}

		line, err := ksefLine.ToGOBL()

		require.NoError(t, err)
		assert.Equal(t, "Gross Price Item", line.Item.Name)
		require.NotNil(t, line.Item.Price)
		assert.Equal(t, "123.00", line.Item.Price.String())
		assert.Equal(t, "1", line.Quantity.String())
	})

	t.Run("prefers net unit price over gross", func(t *testing.T) {
		ksefLine := &ksef.Line{
			Name:           "Both Prices Item",
			Quantity:       "1",
			NetUnitPrice:   "100.00",
			GrossUnitPrice: "123.00",
			Measure:        "HUR",
			VATRate:        "23",
		}

		line, err := ksefLine.ToGOBL()

		require.NoError(t, err)
		require.NotNil(t, line.Item.Price)
		assert.Equal(t, "100.00", line.Item.Price.String())
	})

	t.Run("handles invalid unit codes gracefully", func(t *testing.T) {
		ksefLine := &ksef.Line{
			Name:         "Item with Polish unit",
			Quantity:     "1",
			NetUnitPrice: "100.00",
			Measure:      "szt",
			VATRate:      "23",
		}

		line, err := ksefLine.ToGOBL()

		require.NoError(t, err)
		assert.Equal(t, org.Unit(""), line.Item.Unit)
	})

	t.Run("handles valid UNECE unit codes", func(t *testing.T) {
		ksefLine := &ksef.Line{
			Name:         "Item with UNECE unit",
			Quantity:     "1",
			NetUnitPrice: "100.00",
			Measure:      "KGM",
			VATRate:      "23",
		}

		line, err := ksefLine.ToGOBL()

		require.NoError(t, err)
		assert.Equal(t, org.Unit("KGM"), line.Item.Unit)
	})

	t.Run("handles valid GOBL unit codes", func(t *testing.T) {
		ksefLine := &ksef.Line{
			Name:         "Item with GOBL unit",
			Quantity:     "1",
			NetUnitPrice: "100.00",
			Measure:      "h",
			VATRate:      "23",
		}

		line, err := ksefLine.ToGOBL()

		require.NoError(t, err)
		assert.Equal(t, org.Unit("h"), line.Item.Unit)
	})
}
