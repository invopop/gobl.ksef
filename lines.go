package ksef

import (
	"strings"

	"github.com/invopop/gobl/addons/pl/favat"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
	"github.com/invopop/gobl/uuid"
)

// Line defines the XML structure for KSeF item line (element type FaWiersz, for VAT and KOR type invoices)
type Line struct {
	LineNumber              int    `xml:"NrWierszaFa"`
	UniqueID                string `xml:"UU_ID,omitempty"`
	CompletionDate          string `xml:"P_6A,omitempty"`
	Name                    string `xml:"P_7,omitempty"`
	InternalCode            string `xml:"Indeks,omitempty"`
	GTIN                    string `xml:"GTIN,omitempty"`
	PKWiU                   string `xml:"PKWiU,omitempty"`
	CN                      string `xml:"CN,omitempty"`
	PKOB                    string `xml:"PKOB,omitempty"`
	Measure                 string `xml:"P_8A,omitempty"`
	Quantity                string `xml:"P_8B,omitempty"`
	NetUnitPrice            string `xml:"P_9A,omitempty"`
	GrossUnitPrice          string `xml:"P_9B,omitempty"`
	UnitDiscount            string `xml:"P_10,omitempty"`
	NetPriceTotal           string `xml:"P_11,omitempty"`
	GrossPriceTotal         string `xml:"P_11A,omitempty"`
	VATAmount               string `xml:"P_11Vat,omitempty"`
	VATRate                 string `xml:"P_12,omitempty"`
	OSSTaxRate              string `xml:"P_12_XII,omitempty"` // one stop shop
	Attachment15GoodsMarker int    `xml:"P_12_Zal_15,omitempty"`
	ExciseDuty              string `xml:"KwotaAkcyzy,omitempty"`
	SpecialGoodsCode        string `xml:"GTU,omitempty"` // values GTU_01 to GTU_13
	Procedure               string `xml:"Procedura,omitempty"`
	CurrencyRate            string `xml:"KursWaluty,omitempty"`
	BeforeCorrectionMarker  int    `xml:"StanPrzed,omitempty"`
}

// NewLines generates lines for the KSeF invoice
func NewLines(lines []*bill.Line) []*Line {
	var Lines []*Line

	for _, line := range lines {
		Lines = append(Lines, newLine(line))
	}

	return Lines
}

func newLine(line *bill.Line) *Line {
	l := &Line{
		LineNumber:    line.Index,
		UniqueID:      string(line.UUID),
		Name:          line.Item.Name,
		Measure:       string(line.Item.Unit.UNECE()),
		NetUnitPrice:  line.Item.Price.String(),
		Quantity:      line.Quantity.String(),
		UnitDiscount:  unitDiscount(line),
		NetPriceTotal: line.Total.String(),
	}
	if tc := line.Taxes.Get(tax.CategoryVAT); tc != nil {
		if tc.Ext.Get(favat.ExtKeyTaxCategory) == "5" {
			if tc.Percent != nil {
				l.OSSTaxRate = tc.Percent.Amount().MinimalString()
			}
		} else {
			l.VATRate = vatRate(tc)
		}
	}

	return l
}

// vatRate returns the VAT rate string and OSS tax rate string for a tax combo
// based on the tax category extension
func vatRate(tc *tax.Combo) string {

	// For non-zero percentages, use the percentage value
	if tc.Percent != nil && !tc.Percent.IsZero() {
		return tc.Percent.Amount().MinimalString()
	}

	// For zero/nil percentage, determine from tax category extension
	switch tc.Ext.Get(favat.ExtKeyTaxCategory) {
	case "6.1": // zero-rated goods and services in the country
		return "0 KR"
	case "6.2": // intra-community supply
		return "0 WDT"
	case "6.3": // export supply
		return "0 EX"
	case "7": // tax exempt supply
		return "zw"
	case "8": // outside scope supply
		return "np I"
	case "9": // reverse charge supply
		return "np II"
	case "10": // domestic reverse charge supply
		return "oo"
	default:
		return ""
	}
}

func unitDiscount(line *bill.Line) string {
	if len(line.Discounts) == 0 {
		return ""
	}

	amount := num.MakeAmount(0, 2)

	for _, discount := range line.Discounts {
		amount = amount.Add(discount.Amount)
	}

	discount := amount.Divide(line.Quantity)

	return discount.String()
}

// ToGOBL converts a KSEF Line to a GOBL Line.
func (l *Line) ToGOBL() (*bill.Line, error) {
	line := &bill.Line{
		Item: &org.Item{
			Name: l.Name,
		},
	}

	if l.UniqueID != "" {
		if id, err := uuid.Parse(l.UniqueID); err == nil {
			line.UUID = id
		}
	}

	// Parse quantity
	if l.Quantity != "" {
		qty, err := parseAmount(l.Quantity)
		if err != nil {
			return nil, err
		}
		line.Quantity = qty
	}

	// Parse unit price
	if l.NetUnitPrice != "" {
		price, err := parseAmount(l.NetUnitPrice)
		if err != nil {
			return nil, err
		}
		line.Item.Price = &price
	} else if l.GrossUnitPrice != "" {
		price, err := parseAmount(l.GrossUnitPrice)
		if err != nil {
			return nil, err
		}
		line.Item.Price = &price
	}

	// Parse unit of measure
	if l.Measure != "" {
		line.Item.Unit = parseUnit(l.Measure)
	}

	// Parse discount
	if l.UnitDiscount != "" {
		discount, err := parseAmount(l.UnitDiscount)
		if err != nil {
			return nil, err
		}
		if !discount.IsZero() {
			if !line.Quantity.IsZero() {
				discount = discount.Multiply(line.Quantity)
			}
			line.Discounts = []*bill.LineDiscount{
				{
					Amount: discount,
				},
			}
		}
	}

	// Parse VAT rate and create tax combo
	var rateStr string
	if l.OSSTaxRate != "" {
		rateStr = l.OSSTaxRate
	} else if l.VATRate != "" {
		rateStr = l.VATRate
	}

	if rateStr != "" {
		taxInfo := parseVATRate(rateStr)
		taxCombo := &tax.Combo{
			Category: tax.CategoryVAT,
			Key:      taxInfo.Key,
			Rate:     taxInfo.Rate,
			Percent:  taxInfo.Percent,
			Ext: tax.Extensions{
				favat.ExtKeyTaxCategory: taxInfo.TaxCategory,
			},
		}
		line.Taxes = tax.Set{taxCombo}
	}

	return line, nil
}

// parseAmount parses a string amount to num.Amount
func parseAmount(s string) (num.Amount, error) {
	amt, err := num.AmountFromString(s)
	if err != nil {
		return num.Amount{}, err
	}
	return amt, nil
}

// parseUnit converts KSEF unit code (UNECE) to GOBL unit.
// Returns empty if the code is not a valid GOBL or UN/ECE unit.
func parseUnit(code string) org.Unit {
	u := org.Unit(code)
	if err := u.Validate(); err != nil {
		return ""
	}
	return u
}

// TaxRateInfo contains the parsed tax rate information
type TaxRateInfo struct {
	Key         cbc.Key
	Rate        cbc.Key
	Percent     *num.Percentage
	TaxCategory cbc.Code
}

// parseVATRate converts KSEF VAT rate string to GOBL tax information.
// KSEF uses various formats:
// - "23", "8", "5" for standard rates
// - "0 KR" for zero-rated (6.1)
// - "0 WDT" for intra-community (6.2)
// - "0 EX" for export (6.3)
// - "zw" for exempt (7)
// - "np I" for outside scope (8)
// - "np II" for reverse charge (9)
// - "oo" for domestic reverse charge (10)
func parseVATRate(rateStr string) *TaxRateInfo {
	rateStr = strings.TrimSpace(rateStr)

	info := &TaxRateInfo{}

	switch rateStr {
	case "23":
		info.Key = tax.KeyStandard
		info.Rate = tax.RateGeneral
		pct := num.MakePercentage(230, 3)
		info.Percent = &pct
		info.TaxCategory = "1"
	case "22":
		info.Key = tax.KeyStandard
		pct := num.MakePercentage(220, 3)
		info.Percent = &pct
		info.TaxCategory = "1"
	case "8":
		info.Key = tax.KeyStandard
		info.Rate = tax.RateReduced
		pct := num.MakePercentage(80, 3)
		info.Percent = &pct
		info.TaxCategory = "2"
	case "7":
		info.Key = tax.KeyStandard
		pct := num.MakePercentage(70, 3)
		info.Percent = &pct
		info.TaxCategory = "2"
	case "5":
		info.Key = tax.KeyStandard
		info.Rate = tax.RateSuperReduced
		pct := num.MakePercentage(50, 3)
		info.Percent = &pct
		info.TaxCategory = "3"
	case "4":
		info.Key = tax.KeyStandard
		pct := num.MakePercentage(40, 3)
		info.Percent = &pct
		info.TaxCategory = "4"
	case "3":
		info.Key = tax.KeyStandard
		pct := num.MakePercentage(30, 3)
		info.Percent = &pct
		info.TaxCategory = "3"
	case "0 KR":
		info.Key = tax.KeyZero
		pct := num.MakePercentage(0, 3)
		info.Percent = &pct
		info.TaxCategory = "6.1"
	case "0 WDT":
		info.Key = tax.KeyIntraCommunity
		pct := num.MakePercentage(0, 3)
		info.Percent = &pct
		info.TaxCategory = "6.2"
	case "0 EX":
		info.Key = tax.KeyExport
		pct := num.MakePercentage(0, 3)
		info.Percent = &pct
		info.TaxCategory = "6.3"
	case "zw":
		info.Key = tax.KeyExempt
		info.TaxCategory = "7"
	case "np I":
		info.Key = tax.KeyOutsideScope
		info.TaxCategory = "8"
	case "np II":
		info.Key = tax.KeyReverseCharge
		info.TaxCategory = "9"
	case "oo":
		info.Key = tax.KeyReverseCharge
		info.TaxCategory = "10"
	}

	return info
}
