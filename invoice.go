package ksef

import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/regimes/common"
)

type Fa struct {
	FaWiersz                           []FaWiersz `xml:"FaWiersz"`
	Adnotacje                          *Adnotacje `xml:"Adnotacje"`
	KodWaluty                          string     `xml:"KodWaluty"`
	RodzajFaktury                      string     `xml:"RodzajFaktury"`
	IssueDate                          string     `xml:"P_1"`
	SequentialNumber                   string     `xml:"P_2"`
	CompletionDate                     string     `xml:"P_6,omitempty"`
	StartDate                          string     `xml:"P_6_Od,omitempty"`
	EndDate                            string     `xml:"P_6_Do,omitempty"`
	BasicRateNetSale                   string     `xml:"P_13_1,omitempty"`
	BasicRateTax                       string     `xml:"P_14_1,omitempty"`
	BasicRateTaxConvertedToPln         string     `xml:"P_14_1W,omitempty"`
	FirstReducedRateNetSale            string     `xml:"P_13_2,omitempty"`
	FirstReducedRateTax                string     `xml:"P_14_2,omitempty"`
	FirstReducedRateTaxConvertedToPln  string     `xml:"P_14_2W,omitempty"`
	SecondReducedRateNetSale           string     `xml:"P_13_3,omitempty"`
	SecondReducedRateTax               string     `xml:"P_14_3,omitempty"`
	SecondReducedRateTaxConvertedToPln string     `xml:"P_14_3W,omitempty"`
	TotalAmountRecivable               string     `xml:"P_15"`
}

type Adnotacje struct {
	P_16      int `xml:"P_16"`
	P_17      int `xml:"P_17"`
	P_18      int `xml:"P_18"`
	P_18A     int `xml:"P_18A"`
	P_19N     int `xml:"Zwolnienie>P_19N"`
	P_22N     int `xml:"NoweSrodkiTransportu>P_22N"`
	P_23      int `xml:"P_23"`
	P_PMarzyN int `xml:"PMarzy>P_PMarzyN"`
}

func NewAdnotacje(inv *bill.Invoice) *Adnotacje {
	Adnotacje := &Adnotacje{ //default values for the most common case
		P_16:      2,
		P_17:      2,
		P_18:      2,
		P_18A:     2,
		P_19N:     1,
		P_22N:     1,
		P_23:      2,
		P_PMarzyN: 1,
	}
	return Adnotacje
}

func NewFa(inv *bill.Invoice) *Fa {
	Fa := &Fa{
		Adnotacje:        NewAdnotacje(inv),
		KodWaluty:        string(inv.Currency),
		RodzajFaktury:    string(inv.Type),
		IssueDate:        inv.IssueDate.String(),
		SequentialNumber: inv.Series + inv.Code,
		CompletionDate:   inv.OperationDate.String(),
	}
	cu := inv.Currency.Def().Units
	for _, cat := range inv.Totals.Taxes.Categories {
		if cat.Code != common.TaxCategoryVAT {
			continue
		}

		for _, rate := range cat.Rates {
			if rate.Key == common.TaxRateStandard {
				Fa.BasicRateNetSale = rate.Base.Rescale(cu).String()
				Fa.BasicRateTax = rate.Amount.Rescale(cu).String()
			}
			if rate.Key == common.TaxRateReduced {
				Fa.FirstReducedRateNetSale = rate.Base.Rescale(cu).String()
				Fa.FirstReducedRateTax = rate.Amount.Rescale(cu).String()
			}
			if rate.Key == common.TaxRateSuperReduced {
				Fa.SecondReducedRateNetSale = rate.Base.Rescale(cu).String()
				Fa.SecondReducedRateTax = rate.Amount.Rescale(cu).String()
			}
		}
	}

	return Fa
}
