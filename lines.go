package ksef

import "github.com/invopop/gobl/bill"

type FaWiersz struct {
	NrWierszaFa   int    `xml:"NrWierszaFa"`
	Name          string `xml:"P_7"`
	Measure       string `xml:"P_8A"`
	Quantity      string `xml:"P_8B"`
	NetUnitPrice  string `xml:"P_9A"`
	NetPriceTotal string `xml:"P_11"`
	TaxRate       string `xml:"P_12"`
}

func NewFaWiersz(line *bill.Line) *FaWiersz {
	FaWiersz := &FaWiersz{
		NrWierszaFa:   line.Index,
		Name:          line.Item.Name,
		Measure:       string(line.Item.Unit.UNECE()),
		NetUnitPrice:  line.Item.Price.String(),
		Quantity:      line.Quantity.String(),
		NetPriceTotal: line.Sum.String(),
		TaxRate:       line.Taxes[0].Percent.String(),
	}

	return FaWiersz
}
