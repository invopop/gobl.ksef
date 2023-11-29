package ksef

import "github.com/invopop/gobl/bill"

type Line struct {
	LineNumber    int    `xml:"NrWierszaFa"`
	Name          string `xml:"P_7"`
	Measure       string `xml:"P_8A"`
	Quantity      string `xml:"P_8B"`
	NetUnitPrice  string `xml:"P_9A"`
	NetPriceTotal string `xml:"P_11"`
	TaxRate       string `xml:"P_12"`
}

func NewLine(line *bill.Line) *Line {
	Line := &Line{
		LineNumber:    line.Index,
		Name:          line.Item.Name,
		Measure:       string(line.Item.Unit.UNECE()),
		NetUnitPrice:  line.Item.Price.String(),
		Quantity:      line.Quantity.String(),
		NetPriceTotal: line.Sum.String(),
		TaxRate:       line.Taxes[0].Percent.Rescale(2).StringWithoutSymbol(),
	}

	return Line
}

func NewLines(lines []*bill.Line) []*Line {
	var Lines []*Line

	for _, line := range lines {
		Lines = append(Lines, NewLine(line))
	}

	return Lines
}
