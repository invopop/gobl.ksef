package ksef_test

import (
	"testing"

	ksef "github.com/invopop/gobl.ksef"
	"github.com/invopop/gobl/addons/pl/favat"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewFavatInv(t *testing.T) {
	baseInvoice := func() *bill.Invoice {
		return &bill.Invoice{
			Currency: currency.PLN,
			Supplier: &org.Party{
				TaxID: &tax.Identity{
					Country: l10n.PL.Tax(),
				},
			},
			Tax: &bill.Tax{
				Ext: tax.Extensions{
					favat.ExtKeyInvoiceType: "VAT",
				},
			},
			Totals: &bill.Totals{
				Taxes: &tax.Total{},
			},
		}
	}

	t.Run("sets preceding invoice", func(t *testing.T) {
		inv := baseInvoice()
		inv.Preceding = []*org.DocumentRef{
			{},
		}

		invoice := ksef.NewFavatInv(inv)

		assert.NotNil(t, invoice.CorrectedInv)
	})

	t.Run("sets correction reason", func(t *testing.T) {
		reason := "example reason"

		inv := baseInvoice()
		inv.Preceding = []*org.DocumentRef{
			{
				Reason: reason,
			},
		}

		invoice := ksef.NewFavatInv(inv)

		assert.Equal(t, reason, invoice.CorrectionReason)
	})

	t.Run("sets correction type", func(t *testing.T) {
		inv := baseInvoice()
		inv.Preceding = []*org.DocumentRef{
			{
				Ext: tax.Extensions{
					favat.ExtKeyEffectiveDate: "1",
				},
			},
		}

		invoice := ksef.NewFavatInv(inv)

		assert.Equal(t, "1", invoice.CorrectionType)
	})

	t.Run("sets the self-billing annotation to false in non-self-billed invoices", func(t *testing.T) {
		inv := baseInvoice()

		invoice := ksef.NewFavatInv(inv)

		assert.Equal(t, "2", invoice.Annotations.SelfBilling)
	})

	t.Run("sets the self-billing annotation to true in self-billed invoices", func(t *testing.T) {
		inv := baseInvoice()
		inv.Tax.Ext[favat.ExtKeySelfBilling] = "1"

		invoice := ksef.NewFavatInv(inv)

		assert.Equal(t, "1", invoice.Annotations.SelfBilling)
	})

	t.Run("sets reverse charge annotation", func(t *testing.T) {
		inv := baseInvoice()
		inv.Tax.Ext[favat.ExtKeyReverseCharge] = "1"

		invoice := ksef.NewFavatInv(inv)

		assert.Equal(t, "1", invoice.Annotations.ReverseCharge)
	})

	t.Run("sets cash accounting annotation", func(t *testing.T) {
		inv := baseInvoice()
		inv.Tax.Ext[favat.ExtKeyCashAccounting] = "1"

		invoice := ksef.NewFavatInv(inv)

		assert.Equal(t, "1", invoice.Annotations.CashAccounting)
	})

	t.Run("sets split payment annotation", func(t *testing.T) {
		inv := baseInvoice()
		inv.Tax.Ext[favat.ExtKeySplitPayment] = "1"

		invoice := ksef.NewFavatInv(inv)

		assert.Equal(t, "1", invoice.Annotations.SplitPaymentMechanism)
	})

	t.Run("sets tax exemption annotation with marker", func(t *testing.T) {
		inv := baseInvoice()
		inv.Tax.Ext[favat.ExtKeyExemption] = "A"

		invoice := ksef.NewFavatInv(inv)

		assert.Equal(t, "1", invoice.Annotations.TaxExemption.Marker)
	})

	t.Run("sets margin scheme travel agency", func(t *testing.T) {
		inv := baseInvoice()
		inv.Tax.Ext[favat.ExtKeyMarginScheme] = "2"

		invoice := ksef.NewFavatInv(inv)

		assert.Equal(t, "1", invoice.Annotations.MarginScheme.Marker)
		assert.Equal(t, "1", invoice.Annotations.MarginScheme.TravelAgencyMargin)
	})

	t.Run("sets margin scheme used goods", func(t *testing.T) {
		inv := baseInvoice()
		inv.Tax.Ext[favat.ExtKeyMarginScheme] = "3.1"

		invoice := ksef.NewFavatInv(inv)

		assert.Equal(t, "1", invoice.Annotations.MarginScheme.Marker)
		assert.Equal(t, "1", invoice.Annotations.MarginScheme.UsedGoodsMargin)
	})

	t.Run("sets margin scheme art works", func(t *testing.T) {
		inv := baseInvoice()
		inv.Tax.Ext[favat.ExtKeyMarginScheme] = "3.2"

		invoice := ksef.NewFavatInv(inv)

		assert.Equal(t, "1", invoice.Annotations.MarginScheme.Marker)
		assert.Equal(t, "1", invoice.Annotations.MarginScheme.ArtWorksMargin)
	})

	t.Run("sets margin scheme collectibles and antiques", func(t *testing.T) {
		inv := baseInvoice()
		inv.Tax.Ext[favat.ExtKeyMarginScheme] = "3.3"

		invoice := ksef.NewFavatInv(inv)

		assert.Equal(t, "1", invoice.Annotations.MarginScheme.Marker)
		assert.Equal(t, "1", invoice.Annotations.MarginScheme.CollectiblesAndAntiquesMargin)
	})

	t.Run("sets additional description from notes", func(t *testing.T) {
		inv := baseInvoice()
		inv.Notes = []*org.Note{
			{
				Key:  "general",
				Text: "Test note text",
			},
		}

		invoice := ksef.NewFavatInv(inv)

		assert.Len(t, invoice.AdditionalDescription, 1)
		assert.Equal(t, "general", invoice.AdditionalDescription[0].Key)
		assert.Equal(t, "Test note text", invoice.AdditionalDescription[0].Value)
	})

	t.Run("sets invoice number with series", func(t *testing.T) {
		inv := baseInvoice()
		inv.Series = "INV"
		inv.Code = "001"

		invoice := ksef.NewFavatInv(inv)

		assert.Equal(t, "INV-001", invoice.SequentialNumber)
	})

	t.Run("sets invoice number without series", func(t *testing.T) {
		inv := baseInvoice()
		inv.Code = "001"

		invoice := ksef.NewFavatInv(inv)

		assert.Equal(t, "001", invoice.SequentialNumber)
	})

	t.Run("sets total amount due when due is specified", func(t *testing.T) {
		inv := baseInvoice()
		due, err := num.AmountFromString("10.00")
		require.NoError(t, err)
		inv.Totals.Due = &due

		invoice := ksef.NewFavatInv(inv)

		assert.Equal(t, "10.00", invoice.TotalAmountDue)
	})

	t.Run("sets total amount due from payable when due is nil", func(t *testing.T) {
		inv := baseInvoice()
		payable, err := num.AmountFromString("25.00")
		require.NoError(t, err)
		inv.Totals.Payable = payable

		invoice := ksef.NewFavatInv(inv)

		assert.Equal(t, "25.00", invoice.TotalAmountDue)
	})
}

func TestParseLinesSetsPricesInclude(t *testing.T) {
	t.Run("sets PricesInclude when lines use gross pricing", func(t *testing.T) {
		doc := &ksef.Invoice{
			Seller: &ksef.Seller{
				NIP:  "1234567890",
				Name: "Test Supplier",
				Address: &ksef.Address{
					CountryCode: "PL",
					AddressL1:   "ul. Testowa 1 00-001 Warszawa",
				},
			},
			Buyer: &ksef.Buyer{
				NIP:  "9876543210",
				Name: "Test Buyer",
				Address: &ksef.Address{
					CountryCode: "PL",
					AddressL1:   "ul. Testowa 2 00-002 Warszawa",
				},
				JST: "2",
				GV:  "2",
			},
			Inv: &ksef.Inv{
				CurrencyCode:     "PLN",
				IssueDate:        "2024-06-15",
				SequentialNumber:  "FV-001",
				TotalAmountDue:   "123.00",
				InvoiceType:      "VAT",
				StandardRateNetSale: "100.00",
				StandardRateTax:    "23.00",
				Annotations: &ksef.Annotations{
					CashAccounting:                      "2",
					SelfBilling:                         "2",
					ReverseCharge:                       "2",
					SplitPaymentMechanism:               "2",
					SimplifiedProcedureBySecondTaxpayer: "2",
					TaxExemption: &ksef.TaxExemption{
						NoExemption: "1",
					},
					NewTransportMeans: &ksef.NewTransportMeans{
						NoNewTransportMeans: "1",
					},
					MarginScheme: &ksef.MarginScheme{
						NoMarginScheme: "1",
					},
				},
				Lines: []*ksef.Line{
					{
						LineNumber:     1,
						Name:           "Gross Price Item",
						Quantity:       "1",
						GrossUnitPrice: "123.00",
						Measure:        "HUR",
						VATRate:        "23",
					},
				},
			},
		}

		inv, err := doc.ToGOBL()
		require.NoError(t, err)
		require.NotNil(t, inv.Tax)
		assert.Equal(t, tax.CategoryVAT, inv.Tax.PricesInclude)
	})

	t.Run("does not set PricesInclude for net pricing", func(t *testing.T) {
		doc := &ksef.Invoice{
			Seller: &ksef.Seller{
				NIP:  "1234567890",
				Name: "Test Supplier",
				Address: &ksef.Address{
					CountryCode: "PL",
					AddressL1:   "ul. Testowa 1 00-001 Warszawa",
				},
			},
			Buyer: &ksef.Buyer{
				NIP:  "9876543210",
				Name: "Test Buyer",
				Address: &ksef.Address{
					CountryCode: "PL",
					AddressL1:   "ul. Testowa 2 00-002 Warszawa",
				},
				JST: "2",
				GV:  "2",
			},
			Inv: &ksef.Inv{
				CurrencyCode:     "PLN",
				IssueDate:        "2024-06-15",
				SequentialNumber:  "FV-001",
				TotalAmountDue:   "123.00",
				InvoiceType:      "VAT",
				StandardRateNetSale: "100.00",
				StandardRateTax:    "23.00",
				Annotations: &ksef.Annotations{
					CashAccounting:                      "2",
					SelfBilling:                         "2",
					ReverseCharge:                       "2",
					SplitPaymentMechanism:               "2",
					SimplifiedProcedureBySecondTaxpayer: "2",
					TaxExemption: &ksef.TaxExemption{
						NoExemption: "1",
					},
					NewTransportMeans: &ksef.NewTransportMeans{
						NoNewTransportMeans: "1",
					},
					MarginScheme: &ksef.MarginScheme{
						NoMarginScheme: "1",
					},
				},
				Lines: []*ksef.Line{
					{
						LineNumber:    1,
						Name:          "Net Price Item",
						Quantity:      "1",
						NetUnitPrice:  "100.00",
						Measure:       "HUR",
						VATRate:       "23",
						NetPriceTotal: "100.00",
					},
				},
			},
		}

		inv, err := doc.ToGOBL()
		require.NoError(t, err)
		require.NotNil(t, inv.Tax)
		assert.Empty(t, inv.Tax.PricesInclude)
	})
}

// testCreditNoteDoc creates a KSeF Invoice document for credit note tests.
func testCreditNoteDoc(lines []*ksef.Line, totalDue string) *ksef.Invoice {
	return &ksef.Invoice{
		Seller: &ksef.Seller{
			NIP:  "1234567890",
			Name: "Test Supplier",
			Address: &ksef.Address{
				CountryCode: "PL",
				AddressL1:   "ul. Testowa 1 00-001 Warszawa",
			},
		},
		Buyer: &ksef.Buyer{
			NIP:  "9876543210",
			Name: "Test Buyer",
			Address: &ksef.Address{
				CountryCode: "PL",
				AddressL1:   "ul. Testowa 2 00-002 Warszawa",
			},
			JST: "2",
			GV:  "2",
		},
		Inv: &ksef.Inv{
			CurrencyCode:       "PLN",
			IssueDate:          "2026-01-20",
			SequentialNumber:    "KOR-TEST",
			TotalAmountDue:     totalDue,
			InvoiceType:        "KOR",
			StandardRateNetSale: "-100.00",
			StandardRateTax:    "-23.00",
			CorrectionReason:   "Test correction",
			CorrectedInv: []*ksef.CorrectedInv{
				{IssueDate: "2026-01-15", SequentialNumber: "INV-001"},
			},
			Annotations: &ksef.Annotations{
				CashAccounting:                      "2",
				SelfBilling:                         "2",
				ReverseCharge:                       "2",
				SplitPaymentMechanism:               "2",
				SimplifiedProcedureBySecondTaxpayer: "2",
				TaxExemption: &ksef.TaxExemption{
					NoExemption: "1",
				},
				NewTransportMeans: &ksef.NewTransportMeans{
					NoNewTransportMeans: "1",
				},
				MarginScheme: &ksef.MarginScheme{
					NoMarginScheme: "1",
				},
			},
			Lines: lines,
		},
	}
}

func TestCreditNoteLineInversion(t *testing.T) {
	t.Run("inverts negative qty to positive for differences method", func(t *testing.T) {
		doc := testCreditNoteDoc([]*ksef.Line{
			{
				LineNumber:    1,
				Name:          "Software Services",
				Quantity:      "-10",
				NetUnitPrice:  "10.00",
				Measure:       "HUR",
				VATRate:       "23",
				NetPriceTotal: "-100.00",
			},
		}, "-123.00")

		inv, err := doc.ToGOBL()
		require.NoError(t, err)
		require.Len(t, inv.Lines, 1)

		// Quantity should be inverted from -10 to 10 (GOBL positive convention)
		assert.Equal(t, "10", inv.Lines[0].Quantity.String())
		// Payable should be positive
		assert.Equal(t, "123.00", inv.Totals.Payable.String())
	})

	t.Run("keeps StanPrzed lines positive", func(t *testing.T) {
		doc := testCreditNoteDoc([]*ksef.Line{
			// "After" line (no StanPrzed) — will be inverted
			{
				LineNumber:    1,
				Name:          "Item A",
				Quantity:      "0",
				NetUnitPrice:  "50.00",
				Measure:       "HUR",
				VATRate:       "23",
				NetPriceTotal: "0.00",
			},
			// "Before" line (StanPrzed=1) — stays positive
			{
				LineNumber:             1,
				Name:                   "Item A",
				Quantity:               "2",
				NetUnitPrice:           "50.00",
				Measure:                "HUR",
				VATRate:                "23",
				NetPriceTotal:          "100.00",
				BeforeCorrectionMarker: 1,
			},
		}, "-123.00")

		inv, err := doc.ToGOBL()
		require.NoError(t, err)
		require.Len(t, inv.Lines, 2)

		// "After" line: qty=0 inverted stays 0, no note
		assert.Equal(t, "0", inv.Lines[0].Quantity.String())
		assert.Empty(t, inv.Lines[0].Notes)
		// "Before" (StanPrzed) line: qty=2 stays positive, has note
		assert.Equal(t, "2", inv.Lines[1].Quantity.String())
		require.Len(t, inv.Lines[1].Notes, 1)
		assert.Contains(t, inv.Lines[1].Notes[0].Text, "Before correction")
	})

	t.Run("inverts discount amounts for non-StanPrzed lines", func(t *testing.T) {
		// In KSeF credit notes (differences method), qty is negative, but
		// unit discount is positive. ToGOBL multiplies discount * qty → negative.
		// parseLines then inverts both qty and discount → both positive.
		doc := testCreditNoteDoc([]*ksef.Line{
			{
				LineNumber:    1,
				Name:          "Discounted Item",
				Quantity:      "-2",
				NetUnitPrice:  "100.00",
				UnitDiscount:  "10.00",
				Measure:       "HUR",
				VATRate:       "23",
				NetPriceTotal: "-180.00",
			},
		}, "-221.40")

		inv, err := doc.ToGOBL()
		require.NoError(t, err)
		require.Len(t, inv.Lines, 1)

		// Quantity inverted from -2 to 2
		assert.Equal(t, "2", inv.Lines[0].Quantity.String())
		// Discount = 10.00 * (-2) = -20.00, then inverted to 20.00
		require.Len(t, inv.Lines[0].Discounts, 1)
		assert.Equal(t, "20.00", inv.Lines[0].Discounts[0].Amount.String())
	})

	t.Run("does not invert lines for standard invoices", func(t *testing.T) {
		doc := &ksef.Invoice{
			Seller: &ksef.Seller{
				NIP:  "1234567890",
				Name: "Test Supplier",
				Address: &ksef.Address{
					CountryCode: "PL",
					AddressL1:   "ul. Testowa 1 00-001 Warszawa",
				},
			},
			Buyer: &ksef.Buyer{
				NIP:  "9876543210",
				Name: "Test Buyer",
				Address: &ksef.Address{
					CountryCode: "PL",
					AddressL1:   "ul. Testowa 2 00-002 Warszawa",
				},
				JST: "2",
				GV:  "2",
			},
			Inv: &ksef.Inv{
				CurrencyCode:       "PLN",
				IssueDate:          "2026-01-20",
				SequentialNumber:    "FV-001",
				TotalAmountDue:     "123.00",
				InvoiceType:        "VAT",
				StandardRateNetSale: "100.00",
				StandardRateTax:    "23.00",
				Annotations: &ksef.Annotations{
					CashAccounting:                      "2",
					SelfBilling:                         "2",
					ReverseCharge:                       "2",
					SplitPaymentMechanism:               "2",
					SimplifiedProcedureBySecondTaxpayer: "2",
					TaxExemption: &ksef.TaxExemption{
						NoExemption: "1",
					},
					NewTransportMeans: &ksef.NewTransportMeans{
						NoNewTransportMeans: "1",
					},
					MarginScheme: &ksef.MarginScheme{
						NoMarginScheme: "1",
					},
				},
				Lines: []*ksef.Line{
					{
						LineNumber:    1,
						Name:          "Normal Item",
						Quantity:      "10",
						NetUnitPrice:  "10.00",
						Measure:       "HUR",
						VATRate:       "23",
						NetPriceTotal: "100.00",
					},
				},
			},
		}

		inv, err := doc.ToGOBL()
		require.NoError(t, err)
		require.Len(t, inv.Lines, 1)

		// Quantity should stay as-is for standard invoices
		assert.Equal(t, "10", inv.Lines[0].Quantity.String())
	})
}
