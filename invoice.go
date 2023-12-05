package ksef

/**/
import (
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/regimes/common"
	"github.com/invopop/gobl/regimes/pl"
)

type Inv struct {
	CurrencyCode                       string       `xml:"KodWaluty"`
	IssueDate                          string       `xml:"P_1"`
	IssuePlace                         string       `xml:"P_1M,omitempty"`
	SequentialNumber                   string       `xml:"P_2"`
	CompletionDate                     string       `xml:"P_6,omitempty"`
	StartDate                          string       `xml:"P_6_Od,omitempty"`
	EndDate                            string       `xml:"P_6_Do,omitempty"`
	BasicRateNetSale                   string       `xml:"P_13_1,omitempty"`
	BasicRateTax                       string       `xml:"P_14_1,omitempty"`
	BasicRateTaxConvertedToPln         string       `xml:"P_14_1W,omitempty"`
	FirstReducedRateNetSale            string       `xml:"P_13_2,omitempty"`
	FirstReducedRateTax                string       `xml:"P_14_2,omitempty"`
	FirstReducedRateTaxConvertedToPln  string       `xml:"P_14_2W,omitempty"`
	SecondReducedRateNetSale           string       `xml:"P_13_3,omitempty"`
	SecondReducedRateTax               string       `xml:"P_14_3,omitempty"`
	SecondReducedRateTaxConvertedToPln string       `xml:"P_14_3W,omitempty"`
	TaxiRateNetSale                    string       `xml:"P_13_4,omitempty"`
	TaxiRateTax                        string       `xml:"P_14_4,omitempty"`
	TaxiRateTaxConvertedToPln          string       `xml:"P_14_4W,omitempty"`
	SpecialProcedureNetSale            string       `xml:"P_13_5,omitempty"`
	SpecialProcedureTax                string       `xml:"P_14_5,omitempty"`
	ZeroTaxExceptIntraCommunityNetSale string       `xml:"P_13_6_1,omitempty"`
	IntraCommunityNetSale              string       `xml:"P_13_6_2,omitempty"`
	ExportNetSale                      string       `xml:"P_13_6_3,omitempty"`
	TaxExemptNetSale                   string       `xml:"P_13_7,omitempty"`
	P_13_8                             string       `xml:"P_13_8,omitempty"`
	P_13_9                             string       `xml:"P_13_9,omitempty"`
	P_13_10                            string       `xml:"P_13_10,omitempty"`
	P_13_11                            string       `xml:"P_13_11,omitempty"`
	TotalAmountRecivable               string       `xml:"P_15"`
	ExchangeRate                       string       `xml:"KursWalutyZ"`
	Annotations                        *Annotations `xml:"Adnotacje"`
	InvoiceType                        string       `xml:"RodzajFaktury"`
	FP                                 string       `xml:"FP"`
	TP                                 string       `xml:"TP"`
	Lines                              []*Line      `xml:"FaWiersz"`
	Payment                            *Payment     `xml:"Platnosc"`
}

type Annotations struct {
	CashAccounting                      int `xml:"P_16"`
	SelfBilling                         int `xml:"P_17"`
	ReverseCharge                       int `xml:"P_18"`
	SplitPaymentMechanism               int `xml:"P_18A"`
	NoTaxExemptGoods                    int `xml:"Zwolnienie>P_19N"`
	NoNewTransportIntraCommunitySupply  int `xml:"NoweSrodkiTransportu>P_22N"`
	SimplifiedProcedureBySecondTaxpayer int `xml:"P_23"`
	NoMarginProcedures                  int `xml:"PMarzy>P_PMarzyN"`
}

func NewAnnotations(inv *bill.Invoice) *Annotations {
	// default values for the most common case,
	// For fields P_16 to P_18 and field P_23 2 means "no", 1 means "yes".
	// for others 1 means "yes", no value means "no"
	Annotations := &Annotations{
		CashAccounting:                      2,
		SelfBilling:                         2,
		ReverseCharge:                       2,
		SplitPaymentMechanism:               2,
		NoTaxExemptGoods:                    1,
		NoNewTransportIntraCommunitySupply:  1,
		SimplifiedProcedureBySecondTaxpayer: 2,
		NoMarginProcedures:                  1,
	}
	return Annotations
}

func NewInv(inv *bill.Invoice) *Inv {
	cu := inv.Currency.Def().Units
	Inv := &Inv{
		Annotations:          NewAnnotations(inv),
		CurrencyCode:         string(inv.Currency),
		IssueDate:            inv.IssueDate.String(),
		SequentialNumber:     inv.Series + inv.Code,
		TotalAmountRecivable: inv.Totals.Payable.Rescale(cu).String(),
		Lines:                NewLines(inv.Lines),
		Payment:              NewPayment(inv),
	}

	ss := inv.ScenarioSummary()
	Inv.InvoiceType = ss.Codes[pl.KeyFacturaEInvoiceClass].String()
	if inv.OperationDate != nil {
		Inv.CompletionDate = inv.OperationDate.String()
	}
	for _, cat := range inv.Totals.Taxes.Categories {
		if cat.Code != common.TaxCategoryVAT {
			continue
		}

		for _, rate := range cat.Rates {
			if rate.Key == common.TaxRateStandard {
				Inv.BasicRateNetSale = rate.Base.Rescale(cu).String()
				Inv.BasicRateTax = rate.Amount.Rescale(cu).String()
			}
			if rate.Key == common.TaxRateReduced {
				Inv.FirstReducedRateNetSale = rate.Base.Rescale(cu).String()
				Inv.FirstReducedRateTax = rate.Amount.Rescale(cu).String()
			}
			if rate.Key == common.TaxRateSuperReduced {
				Inv.SecondReducedRateNetSale = rate.Base.Rescale(cu).String()
				Inv.SecondReducedRateTax = rate.Amount.Rescale(cu).String()
			}
		}
	}

	return Inv
}
