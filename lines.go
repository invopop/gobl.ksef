package ksef

import "github.com/invopop/gobl/bill"

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
		UnitDiscount:  line.Total.Subtract(line.Sum).Divide(line.Quantity).String(), // not sure if there should be some rescale or matchPrecision
		NetPriceTotal: line.Total.String(),
		TaxRate:       line.Taxes[0].Percent.Rescale(2).StringWithoutSymbol(),
	}

	return Line
}

// NewLines generates lines for the KSeF invoice
func NewLines(lines []*bill.Line) []*Line {
	var Lines []*Line

	for _, line := range lines {
		Lines = append(Lines, newLine(line))
	}

	return Lines
}
