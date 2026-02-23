package ksef_test

import (
	"testing"

	ksef "github.com/invopop/gobl.ksef"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/currency"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func orderingBaseInvoice() *bill.Invoice {
	return &bill.Invoice{
		Currency: currency.PLN,
		Supplier: &org.Party{
			TaxID: &tax.Identity{
				Country: l10n.PL.Tax(),
			},
		},
		Tax: &bill.Tax{
			Ext: tax.Extensions{
				"pl-favat-invoice-type": "ZAL",
			},
		},
		Totals: &bill.Totals{
			Taxes: &tax.Total{},
		},
	}
}

func mustAmount(s string) num.Amount {
	a, err := num.AmountFromString(s)
	if err != nil {
		panic(err)
	}
	return a
}

func TestNewOrder(t *testing.T) {
	t.Run("skips when no ordering", func(t *testing.T) {
		inv := orderingBaseInvoice()

		result := ksef.NewFavatInv(inv)

		assert.Nil(t, result.Order)
		assert.Nil(t, result.TransactionConditions)
	})

	t.Run("skips when ordering has no purchases", func(t *testing.T) {
		inv := orderingBaseInvoice()
		inv.Ordering = &bill.Ordering{}

		result := ksef.NewFavatInv(inv)

		assert.Nil(t, result.Order)
		assert.Nil(t, result.TransactionConditions)
	})

	t.Run("sets order ref number and date", func(t *testing.T) {
		inv := orderingBaseInvoice()
		issueDate := cal.MakeDate(2026, 1, 15)
		inv.Ordering = &bill.Ordering{
			Purchases: []*org.DocumentRef{
				{
					Code:      "PO-12345",
					IssueDate: &issueDate,
				},
			},
		}

		result := ksef.NewFavatInv(inv)

		require.NotNil(t, result.TransactionConditions)
		require.Len(t, result.TransactionConditions.Orders, 1)
		assert.Equal(t, "PO-12345", result.TransactionConditions.Orders[0].Number)
		assert.Equal(t, "2026-01-15", result.TransactionConditions.Orders[0].Date)
	})

	t.Run("sets order ref without date", func(t *testing.T) {
		inv := orderingBaseInvoice()
		inv.Ordering = &bill.Ordering{
			Purchases: []*org.DocumentRef{
				{
					Code: "PO-99",
				},
			},
		}

		result := ksef.NewFavatInv(inv)

		require.NotNil(t, result.TransactionConditions)
		assert.Equal(t, "PO-99", result.TransactionConditions.Orders[0].Number)
		assert.Equal(t, "", result.TransactionConditions.Orders[0].Date)
	})

	t.Run("sets order amount from payable", func(t *testing.T) {
		inv := orderingBaseInvoice()
		payable := mustAmount("6150.00")
		inv.Ordering = &bill.Ordering{
			Purchases: []*org.DocumentRef{
				{
					Code:    "PO-1",
					Payable: &payable,
				},
			},
		}

		result := ksef.NewFavatInv(inv)

		require.NotNil(t, result.Order)
		assert.Equal(t, "6150.00", result.Order.OrderAmount)
	})

	t.Run("builds order lines from tax rates", func(t *testing.T) {
		inv := orderingBaseInvoice()
		pct23 := num.MakePercentage(230, 3)
		payable := mustAmount("6150.00")
		inv.Ordering = &bill.Ordering{
			Purchases: []*org.DocumentRef{
				{
					Code:        "PO-1",
					Description: "Advance Payment",
					Payable:     &payable,
					Tax: &tax.Total{
						Categories: []*tax.CategoryTotal{
							{
								Code: tax.CategoryVAT,
								Rates: []*tax.RateTotal{
									{
										Base:    mustAmount("5000.00"),
										Percent: &pct23,
										Amount:  mustAmount("1150.00"),
									},
								},
								Amount: mustAmount("1150.00"),
							},
						},
						Sum: mustAmount("1150.00"),
					},
				},
			},
		}

		result := ksef.NewFavatInv(inv)

		require.NotNil(t, result.Order)
		assert.Equal(t, "6150.00", result.Order.OrderAmount)
		require.Len(t, result.Order.LineItems, 1)

		ol := result.Order.LineItems[0]
		assert.Equal(t, 1, ol.LineNumber)
		assert.Equal(t, "Advance Payment", ol.Name)
		assert.Equal(t, "1", ol.Quantity)
		assert.Equal(t, "5000.00", ol.NetPriceTotal)
		assert.Equal(t, "5000.00", ol.NetUnitPrice)
		assert.Equal(t, "1150.00", ol.TaxValue)
		assert.Equal(t, "23", ol.VATRate)
	})

	t.Run("builds multiple order lines from multiple rates", func(t *testing.T) {
		inv := orderingBaseInvoice()
		pct23 := num.MakePercentage(230, 3)
		pct8 := num.MakePercentage(80, 3)
		payable := mustAmount("10000.00")
		inv.Ordering = &bill.Ordering{
			Purchases: []*org.DocumentRef{
				{
					Code:        "PO-2",
					Description: "Mixed rates order",
					Payable:     &payable,
					Tax: &tax.Total{
						Categories: []*tax.CategoryTotal{
							{
								Code: tax.CategoryVAT,
								Rates: []*tax.RateTotal{
									{
										Base:    mustAmount("3000.00"),
										Percent: &pct23,
										Amount:  mustAmount("690.00"),
									},
									{
										Base:    mustAmount("2000.00"),
										Percent: &pct8,
										Amount:  mustAmount("160.00"),
									},
								},
							},
						},
					},
				},
			},
		}

		result := ksef.NewFavatInv(inv)

		require.NotNil(t, result.Order)
		require.Len(t, result.Order.LineItems, 2)

		assert.Equal(t, 1, result.Order.LineItems[0].LineNumber)
		assert.Equal(t, "3000.00", result.Order.LineItems[0].NetPriceTotal)
		assert.Equal(t, "690.00", result.Order.LineItems[0].TaxValue)
		assert.Equal(t, "23", result.Order.LineItems[0].VATRate)

		assert.Equal(t, 2, result.Order.LineItems[1].LineNumber)
		assert.Equal(t, "2000.00", result.Order.LineItems[1].NetPriceTotal)
		assert.Equal(t, "160.00", result.Order.LineItems[1].TaxValue)
		assert.Equal(t, "8", result.Order.LineItems[1].VATRate)
	})

	t.Run("handles purchase without tax data", func(t *testing.T) {
		inv := orderingBaseInvoice()
		payable := mustAmount("1000.00")
		inv.Ordering = &bill.Ordering{
			Purchases: []*org.DocumentRef{
				{
					Code:    "PO-3",
					Payable: &payable,
				},
			},
		}

		result := ksef.NewFavatInv(inv)

		require.NotNil(t, result.Order)
		assert.Equal(t, "1000.00", result.Order.OrderAmount)
		assert.Empty(t, result.Order.LineItems)
	})

	t.Run("handles purchase without payable", func(t *testing.T) {
		inv := orderingBaseInvoice()
		inv.Ordering = &bill.Ordering{
			Purchases: []*org.DocumentRef{
				{
					Code: "PO-4",
				},
			},
		}

		result := ksef.NewFavatInv(inv)

		require.NotNil(t, result.Order)
		assert.Equal(t, "", result.Order.OrderAmount)
	})

	t.Run("does not create empty Contracts in TransactionConditions", func(t *testing.T) {
		inv := orderingBaseInvoice()
		inv.Ordering = &bill.Ordering{
			Purchases: []*org.DocumentRef{
				{
					Code: "PO-5",
				},
			},
		}

		result := ksef.NewFavatInv(inv)

		require.NotNil(t, result.TransactionConditions)
		assert.Nil(t, result.TransactionConditions.Contracts)
	})
}

func TestNewOrderRoundTrip(t *testing.T) {
	t.Run("ordering data survives GOBL to KSeF to GOBL", func(t *testing.T) {
		// Build a KSeF doc with order data (simulating what newOrder produces)
		doc := testPrepaymentDoc()

		// Parse back to GOBL
		inv, err := doc.ToGOBL()
		require.NoError(t, err)
		require.NotNil(t, inv.Ordering)
		require.Len(t, inv.Ordering.Purchases, 1)

		ref := inv.Ordering.Purchases[0]
		assert.Equal(t, cbc.Code("PO-12345"), ref.Code)
		require.NotNil(t, ref.IssueDate)
		assert.Equal(t, "2026-01-10", ref.IssueDate.String())
		require.NotNil(t, ref.Payable)
		assert.Equal(t, "5000.00", ref.Payable.String())
		assert.Equal(t, "Widget A, Widget B", ref.Description)

		// Now convert that GOBL invoice back to KSeF
		result := ksef.NewFavatInv(inv)

		// Verify the KSeF structures match
		require.NotNil(t, result.TransactionConditions)
		require.Len(t, result.TransactionConditions.Orders, 1)
		assert.Equal(t, "PO-12345", result.TransactionConditions.Orders[0].Number)
		assert.Equal(t, "2026-01-10", result.TransactionConditions.Orders[0].Date)

		require.NotNil(t, result.Order)
		assert.Equal(t, "5000.00", result.Order.OrderAmount)
		require.Len(t, result.Order.LineItems, 2)
		assert.Equal(t, "3000.00", result.Order.LineItems[0].NetPriceTotal)
		assert.Equal(t, "690.00", result.Order.LineItems[0].TaxValue)
		assert.Equal(t, "23", result.Order.LineItems[0].VATRate)
		assert.Equal(t, "2000.00", result.Order.LineItems[1].NetPriceTotal)
		assert.Equal(t, "460.00", result.Order.LineItems[1].TaxValue)
	})
}
