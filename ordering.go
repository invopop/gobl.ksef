package ksef

import (
	"fmt"
	"strconv"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

// Order defines the XML structure for KSeF "Zamowienie" (order) field, required for ZAL and KOR_ZAL types
type Order struct {
	OrderAmount string       `xml:"WartoscZamowienia"`
	LineItems   []*OrderLine `xml:"ZamowienieWiersz,omitempty"`
}

// OrderLine defines the XML structure for KSeF item line (element type ZamowienieWiersz, for ZAL and KOR_ZAL type invoices)
type OrderLine struct {
	LineNumber              int    `xml:"NrWierszaZam"`
	UniqueID                string `xml:"UU_IDZ,omitempty"`
	Name                    string `xml:"P_7Z,omitempty"`
	InternalCode            string `xml:"IndeksZ,omitempty"`
	GTIN                    string `xml:"GTINZ,omitempty"`
	PKWiU                   string `xml:"PKWiUZ,omitempty"`
	CN                      string `xml:"CNZ,omitempty"`
	PKOB                    string `xml:"PKOBZ,omitempty"`
	Measure                 string `xml:"P_8AZ,omitempty"`
	Quantity                string `xml:"P_8BZ,omitempty"`
	NetUnitPrice            string `xml:"P_9AZ,omitempty"`
	NetPriceTotal           string `xml:"P_11NettoZ,omitempty"`
	TaxValue                string `xml:"P_11VatZ,omitempty"`
	VATRate                 string `xml:"P_12Z,omitempty"`
	OSSTaxRate              string `xml:"P_12Z_XII,omitempty"` // one stop shop
	Attachment15GoodsMarker int    `xml:"P_12Z_Zal_15,omitempty"`
	SpecialGoodsCode        string `xml:"GTUZ,omitempty"` // values GTU_01 to GTU_13
	Procedure               string `xml:"ProceduraZ,omitempty"`
	ExciseDuty              string `xml:"KwotaAkcyzyZ,omitempty"`
	BeforeCorrectionMarker  int    `xml:"StanPrzedZ,omitempty"`
}

// TransactionConditions defines the XML structure for transaction conditions
type TransactionConditions struct {
	Contracts         []*Contract  `xml:"Umowy,omitempty"`
	Orders            []*OrderRef  `xml:"Zamowienia,omitempty"`
	BatchNumbers      []string     `xml:"NrPartiiTowaru,omitempty"`
	DeliveryTerms     string       `xml:"WarunkiDostawy,omitempty"`
	ContractRate      string       `xml:"KursUmowny,omitempty"`
	ContractCurrency  string       `xml:"WalutaUmowna,omitempty"`
	Transport         []*Transport `xml:"Transport,omitempty"`
	IntermediaryParty int          `xml:"PodmiotPosredniczacy,omitempty"`
}

// Contract defines the XML structure for contract reference
type Contract struct {
	Date   string `xml:"DataUmowy"`
	Number string `xml:"NrUmowy"`
}

// OrderRef defines the XML structure for order reference
type OrderRef struct {
	Date   string `xml:"DataZamowienia"`
	Number string `xml:"NrZamowienia"`
}

// Transport defines the XML structure for transport information
type Transport struct {
	TransportType        string     `xml:"RodzajTransportu,omitempty"`
	OtherTransportType   int        `xml:"TransportInny,omitempty"`
	OtherTransportDesc   string     `xml:"OpisInnegoTransportu,omitempty"`
	Carrier              *Carrier   `xml:"Przewoznik,omitempty"`
	TransportOrderNumber string     `xml:"NrZleceniaTransportu,omitempty"`
	CargoType            string     `xml:"OpisLadunku,omitempty"`
	OtherCargoType       int        `xml:"LadunekInny,omitempty"`
	OtherCargoDesc       string     `xml:"OpisInnegoLadunku,omitempty"`
	PackagingUnit        string     `xml:"JednostkaOpakowania,omitempty"`
	TransportStartTime   string     `xml:"DataGodzRozpTransportu,omitempty"`
	TransportEndTime     string     `xml:"DataGodzZakTransportu,omitempty"`
	ShipFrom             *Address   `xml:"WysylkaZ,omitempty"`
	ShipVia              []*Address `xml:"WysylkaPrzez,omitempty"`
	ShipTo               *Address   `xml:"WysylkaDo,omitempty"`
}

// Carrier defines the XML structure for carrier information
type Carrier struct {
	IdentificationData *Buyer   `xml:"DaneIdentyfikacyjne"`
	Address            *Address `xml:"AdresPrzewoznika"`
}

// newInvoicePeriod converts GOBL ordering period to KSeF InvoicePeriod.
func newInvoicePeriod(ordering *bill.Ordering) *InvoicePeriod {
	if ordering == nil || ordering.Period == nil {
		return nil
	}

	return &InvoicePeriod{
		StartDate: ordering.Period.Start.String(),
		EndDate:   ordering.Period.End.String(),
	}
}

// newOrder converts ordering purchases to KSeF Order and TransactionConditions.
// This is the inverse of parseOrderingLines().
func newOrder(ordering *bill.Ordering) (*Order, *TransactionConditions) {
	if ordering == nil || len(ordering.Purchases) == 0 {
		return nil, nil
	}

	ref := ordering.Purchases[0]

	// Build TransactionConditions with order reference
	tc := &TransactionConditions{
		Orders: []*OrderRef{
			{
				Number: ref.Code.String(),
			},
		},
	}
	if ref.IssueDate != nil {
		tc.Orders[0].Date = ref.IssueDate.String()
	}

	// Build Order
	order := &Order{}
	if ref.Payable != nil {
		order.OrderAmount = ref.Payable.String()
	}

	// Build order lines from tax rates
	if ref.Tax != nil && len(ref.Tax.Categories) > 0 {
		cat := ref.Tax.Categories[0]
		for i, rate := range cat.Rates {
			ol := &OrderLine{
				LineNumber: i + 1,
				Name:       ref.Description,
				Quantity:   "1",
			}
			if !rate.Base.IsZero() {
				ol.NetPriceTotal = rate.Base.String()
				ol.NetUnitPrice = rate.Base.String()
			}
			if !rate.Amount.IsZero() {
				ol.TaxValue = rate.Amount.String()
			}
			if rate.Percent != nil {
				ol.VATRate = rate.Percent.Amount().MinimalString()
			}
			order.LineItems = append(order.LineItems, ol)
		}
	}

	return order, tc
}

// parseOrderingLines maps Zamowienie (order) data into Ordering.Purchases.
// Any invoice type can have Zamowienie and WarunkiTransakcji data.
func (inv *Inv) parseOrderingLines(goblInv *bill.Invoice) error {
	if inv.Order == nil {
		return nil
	}

	ref := &org.DocumentRef{}

	// Set order date and number from WarunkiTransakcji.Zamowienia
	if inv.TransactionConditions != nil && len(inv.TransactionConditions.Orders) > 0 {
		order := inv.TransactionConditions.Orders[0]
		if order.Date != "" {
			date, err := parseDate(order.Date)
			if err != nil {
				return fmt.Errorf("parsing order date: %w", err)
			}
			ref.IssueDate = &date
		}
		if order.Number != "" {
			ref.Code = cbc.Code(order.Number)
		}
	}

	// DocumentRef requires a non-empty code
	if ref.Code == "" {
		ref.Code = "unknown"
	}

	// Set payable from order total
	if inv.Order.OrderAmount != "" {
		payable, err := parseAmount(inv.Order.OrderAmount)
		if err != nil {
			return fmt.Errorf("parsing order amount: %w", err)
		}
		ref.Payable = &payable
	}

	// Build tax total and description from order lines (one rate per line)
	if len(inv.Order.LineItems) > 0 {
		var rates []*tax.RateTotal
		var taxSum num.Amount
		var descriptions []string

		for _, ol := range inv.Order.LineItems {
			if ol.Name != "" {
				descriptions = append(descriptions, ol.Name)
			}

			rt := &tax.RateTotal{}

			if ol.NetPriceTotal != "" {
				netAmt, err := parseAmount(ol.NetPriceTotal)
				if err != nil {
					return fmt.Errorf("parsing order line net: %w", err)
				}
				rt.Base = netAmt
			}

			if ol.TaxValue != "" {
				taxAmt, err := parseAmount(ol.TaxValue)
				if err != nil {
					return fmt.Errorf("parsing order line tax: %w", err)
				}
				rt.Amount = taxAmt
				taxSum = taxSum.MatchPrecision(taxAmt)
				taxSum = taxSum.Add(taxAmt)
			}

			if ol.VATRate != "" {
				rateF, err := strconv.ParseFloat(ol.VATRate, 64)
				if err == nil && rateF > 0 {
					pct := num.MakePercentage(int64(rateF*10), 3)
					rt.Percent = &pct
				}
			}

			rates = append(rates, rt)
		}

		if len(descriptions) > 0 {
			ref.Description = descriptions[0]
			for _, d := range descriptions[1:] {
				ref.Description += ", " + d
			}
		}

		if len(rates) > 0 {
			ref.Tax = &tax.Total{
				Categories: []*tax.CategoryTotal{
					{
						Code:   tax.CategoryVAT,
						Rates:  rates,
						Amount: taxSum,
					},
				},
				Sum: taxSum,
			}
		}
	}

	// Add to ordering purchases
	if goblInv.Ordering == nil {
		goblInv.Ordering = &bill.Ordering{}
	}
	goblInv.Ordering.Purchases = append(goblInv.Ordering.Purchases, ref)

	return nil
}
