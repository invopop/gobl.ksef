package ksef

/**/
import (
	"slices"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/regimes/pl"
	"github.com/invopop/gobl/tax"
)

const (
	regionDomestic = "domestic"
	regionEU       = "EU"
	regionNonEU    = "non-EU"
)

// Inv defines the XML structure for KSeF invoice
type Inv struct {
	CurrencyCode                      string        `xml:"KodWaluty"`
	IssueDate                         string        `xml:"P_1"`
	IssuePlace                        string        `xml:"P_1M,omitempty"`
	SequentialNumber                  string        `xml:"P_2"`
	CompletionDate                    string        `xml:"P_6,omitempty"`
	StartDate                         string        `xml:"P_6_Od,omitempty"`
	EndDate                           string        `xml:"P_6_Do,omitempty"`
	StandardRateNetSale               string        `xml:"P_13_1,omitempty"`
	StandardRateTax                   string        `xml:"P_14_1,omitempty"`
	StandardRateTaxConvertedToPln     string        `xml:"P_14_1W,omitempty"`
	ReducedRateNetSale                string        `xml:"P_13_2,omitempty"`
	ReducedRateTax                    string        `xml:"P_14_2,omitempty"`
	ReducedRateTaxConvertedToPln      string        `xml:"P_14_2W,omitempty"`
	SuperReducedRateNetSale           string        `xml:"P_13_3,omitempty"`
	SuperReducedRateTax               string        `xml:"P_14_3,omitempty"`
	SuperReducedRateTaxConvertedToPln string        `xml:"P_14_3W,omitempty"`
	TaxiRateNetSale                   string        `xml:"P_13_4,omitempty"`
	TaxiRateTax                       string        `xml:"P_14_4,omitempty"`
	TaxiRateTaxConvertedToPln         string        `xml:"P_14_4W,omitempty"`
	SpecialProcedureNetSale           string        `xml:"P_13_5,omitempty"`
	SpecialProcedureTax               string        `xml:"P_14_5,omitempty"`
	DomesticZeroTaxNetSale            string        `xml:"P_13_6_1,omitempty"`
	EUZeroTaxNetSale                  string        `xml:"P_13_6_2,omitempty"`
	ExportNetSale                     string        `xml:"P_13_6_3,omitempty"`
	TaxExemptNetSale                  string        `xml:"P_13_7,omitempty"`
	TaxNAInternationalNetSale         string        `xml:"P_13_8,omitempty"`
	TaxNAEUNetSale                    string        `xml:"P_13_9,omitempty"`
	EUServiceNetSale                  string        `xml:"P_13_10,omitempty"`
	MarginNetSale                     string        `xml:"P_13_11,omitempty"`
	TotalAmountReceivable             string        `xml:"P_15"`
	Annotations                       *Annotations  `xml:"Adnotacje"`
	InvoiceType                       string        `xml:"RodzajFaktury"`
	CorrectionReason                  string        `xml:"PrzyczynaKorekty,omitempty"`
	CorrectionType                    string        `xml:"TypKorekty,omitempty"`
	CorrectedInv                      *CorrectedInv `xml:"DaneFaKorygowanej,omitempty"`
	Lines                             []*Line       `xml:"FaWiersz"`
	Payment                           *Payment      `xml:"Platnosc"`
}

// Annotations defines the XML structure for KSeF annotations
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

// newAnnotations sets annotations data
func newAnnotations(inv *bill.Invoice) *Annotations {
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

	if inv.Tax != nil && slices.Contains(inv.Tax.Tags, tax.TagReverseCharge) {
		Annotations.ReverseCharge = 1
	}

	return Annotations
}

// NewInv gets invoice data from GOBL invoice
func NewInv(inv *bill.Invoice) *Inv {
	cu := inv.Currency.Def().Subunits
	Inv := &Inv{
		Annotations:           newAnnotations(inv),
		CurrencyCode:          string(inv.Currency),
		IssueDate:             inv.IssueDate.String(),
		SequentialNumber:      invoiceNumber(inv.Series, inv.Code),
		TotalAmountReceivable: inv.Totals.Payable.Rescale(cu).String(),
		Lines:                 NewLines(inv.Lines),
		Payment:               NewPayment(inv.Payment, inv.Totals),
	}

	if len(inv.Preceding) > 0 {
		for _, prc := range inv.Preceding {
			Inv.CorrectedInv = NewCorrectedInv(prc)
			Inv.CorrectionReason = prc.Reason
			if prc.Ext.Has(pl.ExtKeyKSeFEffectiveDate) {
				Inv.CorrectionType = prc.Ext[pl.ExtKeyKSeFEffectiveDate].Code().String()
			}
		}
	}

	ss := inv.ScenarioSummary()
	Inv.InvoiceType = ss.Codes[pl.KeyFAVATInvoiceType].String()
	if inv.OperationDate != nil {
		Inv.CompletionDate = inv.OperationDate.String()
	}

	reg := region(inv)

	for _, cat := range inv.Totals.Taxes.Categories {
		if cat.Code != tax.CategoryVAT {
			continue
		}

		for _, rate := range cat.Rates {
			setTaxRate(Inv, rate, cu, reg)
		}
	}

	return Inv
}

func invoiceNumber(series string, code string) string {
	if series == "" {
		return code
	}
	return series + "-" + code
}

func setTaxRate(inv *Inv, rate *tax.RateTotal, cu uint32, region string) {
	if rate.Percent == nil {
		return
	}

	base := rate.Base.Rescale(cu).String()
	taxAmount := rate.Amount.Rescale(cu).String()

	switch rate.Key {
	case tax.RateStandard:
		inv.StandardRateNetSale = base
		inv.StandardRateTax = taxAmount
	case tax.RateReduced:
		inv.ReducedRateNetSale = base
		inv.ReducedRateTax = taxAmount
	case tax.RateSuperReduced:
		inv.SuperReducedRateNetSale = base
		inv.SuperReducedRateTax = taxAmount
	case tax.RateSpecial:
		if rate.Ext.Has(pl.ExtKeyKSeFVATSpecial) && rate.Ext[pl.ExtKeyKSeFVATSpecial].String() == "taxi" {
			inv.TaxiRateNetSale = base
			inv.TaxiRateTax = taxAmount
		}
	case tax.RateZero:
		switch region {
		case regionDomestic:
			inv.DomesticZeroTaxNetSale = base
		case regionEU:
			inv.EUZeroTaxNetSale = base
		case regionNonEU:
			inv.ExportNetSale = base
		}
	case tax.RateExempt:
		inv.TaxExemptNetSale = base
	case pl.TaxRateNotPursuant:
		switch region {
		case regionEU:
			inv.TaxNAEUNetSale = base
		case regionNonEU:
			inv.TaxNAInternationalNetSale = base
		}
	}
}

func region(inv *bill.Invoice) string {
	if inv.Supplier == nil || inv.Customer == nil || inv.Supplier.TaxID == nil || inv.Customer.TaxID == nil {
		return regionDomestic
	}
	if isEUCountry(inv.Supplier.TaxID.Country) || isEUCountry(inv.Customer.TaxID.Country) {
		return regionEU
	}
	if inv.Supplier.TaxID.Country != l10n.PL || inv.Customer.TaxID.Country != l10n.PL {
		return regionNonEU
	}
	return regionDomestic
}
