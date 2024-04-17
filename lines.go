package ksef

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/regimes/pl"
	"github.com/invopop/gobl/tax"
)

// Line defines the XML structure for KSeF item line
type Line struct {
	LineNumber              int    `xml:"NrWierszaFa"`
	Name                    string `xml:"P_7,omitempty"`
	Measure                 string `xml:"P_8A,omitempty"`
	Quantity                string `xml:"P_8B,omitempty"`
	NetUnitPrice            string `xml:"P_9A,omitempty"`
	UnitDiscount            string `xml:"P_10,omitempty"`
	NetPriceTotal           string `xml:"P_11,omitempty"`
	TaxRate                 string `xml:"P_12,omitempty"`
	ExciseDuty              string `xml:"KwotaAkcyzy,omitempty"`
	SpecialGoodsCode        string `xml:"GTU,omitempty"` // values GTU_1 to GTU_13
	OSSTaxRate              string `xml:"P_12_XII,omitempty"`
	Attachment15GoodsMarker string `xml:"P_12_Zal_15,omitempty"`
	Procedure               string `xml:"Procedura,omitempty"`
	BeforeCorrectionMarker  string `xml:"StanPrzed,omitempty"`
}

func newLine(line *bill.Line) *Line {
	Line := &Line{
		LineNumber:    line.Index,
		Name:          line.Item.Name,
		Measure:       string(line.Item.Unit.UNECE()),
		NetUnitPrice:  line.Item.Price.String(),
		Quantity:      line.Quantity.String(),
		UnitDiscount:  unitDiscount(line),
		NetPriceTotal: line.Total.String(),
		TaxRate:       newTaxRate(line.Taxes.Get(tax.CategoryVAT)),
	}

	return Line
}

// newTaxRate returns tax rate as string value with one of the values:
// "23", "22", "8", "7", "5", "4", "3", "0", "np", "zw"
func newTaxRate(t *tax.Combo) string {
	switch t.Rate {
	case tax.RateZero:
		return "0"
	case tax.RateExempt:
		return "zw"
	case pl.TaxRateNotPursuant:
		return "np"
	default:
		return t.Percent.Rescale(2).StringWithoutSymbol()
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

// NewLines generates lines for the KSeF invoice
func NewLines(lines []*bill.Line) []*Line {
	var Lines []*Line

	for _, line := range lines {
		Lines = append(Lines, newLine(line))
	}

	return Lines
}
