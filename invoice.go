package ksef

import (
	"fmt"
	"regexp"
	"time"

	"github.com/invopop/gobl/addons/pl/favat"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

// validCodeRe matches valid cbc.Code values (alphanumeric, hyphens, dots, underscores).
var validCodeRe = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)

// Inv defines the XML structure for KSeF invoice
type Inv struct {
	CurrencyCode                       string                       `xml:"KodWaluty"`
	IssueDate                          string                       `xml:"P_1"`
	IssuePlace                         string                       `xml:"P_1M,omitempty"`
	SequentialNumber                   string                       `xml:"P_2"`
	WarehouseDocuments                 []string                     `xml:"WZ,omitempty"`
	CompletionDate                     string                       `xml:"P_6,omitempty"`
	Period                             *InvoicePeriod               `xml:"OkresFa,omitempty"`
	ExchangeRate                       string                       `xml:"KursWalutyZ,omitempty"`
	StandardRateNetSale                string                       `xml:"P_13_1,omitempty"`
	StandardRateTax                    string                       `xml:"P_14_1,omitempty"`
	StandardRateTaxConvertedToPln      string                       `xml:"P_14_1W,omitempty"`
	ReducedRateNetSale                 string                       `xml:"P_13_2,omitempty"`
	ReducedRateTax                     string                       `xml:"P_14_2,omitempty"`
	ReducedRateTaxConvertedToPln       string                       `xml:"P_14_2W,omitempty"`
	SuperReducedRateNetSale            string                       `xml:"P_13_3,omitempty"`
	SuperReducedRateTax                string                       `xml:"P_14_3,omitempty"`
	SuperReducedRateTaxConvertedToPln  string                       `xml:"P_14_3W,omitempty"`
	TaxiRateNetSale                    string                       `xml:"P_13_4,omitempty"`
	TaxiRateTax                        string                       `xml:"P_14_4,omitempty"`
	TaxiRateTaxConvertedToPln          string                       `xml:"P_14_4W,omitempty"`
	OSSNetSale                         string                       `xml:"P_13_5,omitempty"`
	OSSTax                             string                       `xml:"P_14_5,omitempty"`
	ZeroTaxExceptIntraCommunityNetSale string                       `xml:"P_13_6_1,omitempty"`
	IntraCommunityNetSale              string                       `xml:"P_13_6_2,omitempty"`
	ExportNetSale                      string                       `xml:"P_13_6_3,omitempty"`
	TaxExemptNetSale                   string                       `xml:"P_13_7,omitempty"`
	OutsideScopeNetSale                string                       `xml:"P_13_8,omitempty"`
	ReverseChargeNetSale               string                       `xml:"P_13_9,omitempty"`
	DomesticReverseChargeNetSale       string                       `xml:"P_13_10,omitempty"`
	MarginNetSale                      string                       `xml:"P_13_11,omitempty"`
	TotalAmountDue                     string                       `xml:"P_15"`
	AmountBeforeCorrection             string                       `xml:"P_15ZK,omitempty"`
	Annotations                        *Annotations                 `xml:"Adnotacje"`
	InvoiceType                        string                       `xml:"RodzajFaktury"`
	CorrectionReason                   string                       `xml:"PrzyczynaKorekty,omitempty"`
	CorrectionType                     string                       `xml:"TypKorekty,omitempty"`
	CorrectedInv                       []*CorrectedInv              `xml:"DaneFaKorygowanej,omitempty"`
	AdvanceInvoices                    []*AdvanceInvoiceRef         `xml:"FakturaZaliczkowa,omitempty"`
	PartialAdvancePayments             []*PartialAdvancePayment     `xml:"ZaliczkaCzesciowa,omitempty"`
	FP                                 int                          `xml:"FP,omitempty"`
	TP                                 int                          `xml:"TP,omitempty"`
	ExciseTaxRefund                    int                          `xml:"ZwrotAkcyzy,omitempty"`
	AdditionalDescription              []*AdditionalDescriptionLine `xml:"DodatkowyOpis,omitempty"`
	Lines                              []*Line                      `xml:"FaWiersz,omitempty"` // empty for ZAL and KOR_ZAL, use Order instead
	Settlement                         *Settlement                  `xml:"Rozliczenie,omitempty"`
	Payment                            *Payment                     `xml:"Platnosc,omitempty"`
	TransactionConditions              *TransactionConditions       `xml:"WarunkiTransakcji,omitempty"`
	Order                              *Order                       `xml:"Zamowienie,omitempty"` // for ZAL and KOR_ZAL types
}

type InvoicePeriod struct {
	StartDate string `xml:"P_6_Od,omitempty"`
	EndDate   string `xml:"P_6_Do,omitempty"`
}

// AdditionalDescriptionLine defines the XML structure for KSeF additional description line (`DodatkowyOpis`)
type AdditionalDescriptionLine struct {
	LineNumber string `xml:"NrWiersza,omitempty"`
	Key        string `xml:"Klucz"`
	Value      string `xml:"Wartosc"`
}

// AdvanceInvoiceRef defines the XML structure for advance invoice reference
type AdvanceInvoiceRef struct {
	KSeFMarker           int    `xml:"NrKSeFZN,omitempty"`
	AdvanceInvoiceNo     string `xml:"NrFaZaliczkowej,omitempty"`
	KSeFAdvanceInvoiceNo string `xml:"NrKSeFFaZaliczkowej,omitempty"`
}

// PartialAdvancePayment defines the XML structure for partial advance payment (ZaliczkaCzesciowa)
type PartialAdvancePayment struct {
	PaymentDate          string `xml:"P_6Z"`
	PaymentAmount        string `xml:"P_15Z"`
	CurrencyExchangeRate string `xml:"KursWalutyZW,omitempty"`
}

// Settlement defines the XML structure for additional charges and deductions
type Settlement struct {
	Charges         []*ChargeOrDeduction `xml:"Obciazenia>Obciazenie,omitempty"`
	TotalCharges    string               `xml:"Obciazenia>SumaObciazen,omitempty"`
	Deductions      []*ChargeOrDeduction `xml:"Odliczenia>Odliczenie,omitempty"`
	TotalDeductions string               `xml:"Odliczenia>SumaOdliczen,omitempty"`
	AmountToPay     string               `xml:"DoZaplaty,omitempty"`
	AmountToSettle  string               `xml:"DoRozliczenia,omitempty"`
}

// ChargeOrDeduction defines the XML structure for a single charge or deduction
type ChargeOrDeduction struct {
	Amount string `xml:"Kwota"`
	Reason string `xml:"Powod"`
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

// NewFavatInv gets invoice data from GOBL invoice
func NewFavatInv(invoice *bill.Invoice) *Inv {

	inv := &Inv{
		CurrencyCode:     invoice.Currency.String(),
		IssueDate:        invoice.IssueDate.String(),
		Period:           newInvoicePeriod(invoice.Ordering),
		SequentialNumber: invoiceNumber(invoice.Series, invoice.Code),
		Annotations:      newAnnotations(invoice),
		Lines:            NewLines(invoice.Lines),
		Payment:          NewPayment(invoice.Payment, invoice.Totals),
	}

	if invoice.Totals.Due != nil {
		inv.TotalAmountDue = invoice.Totals.Due.String()
	} else {
		inv.TotalAmountDue = invoice.Totals.Payable.String()
	}

	if invoice.Tax != nil && invoice.Tax.Ext != nil {
		inv.InvoiceType = invoice.Tax.Ext.Get(favat.ExtKeyInvoiceType).String()
	}

	inv.setTaxRates(invoice.Totals.Taxes)

	if len(invoice.Notes) > 0 {
		for _, note := range invoice.Notes {
			inv.AdditionalDescription = append(inv.AdditionalDescription, &AdditionalDescriptionLine{
				Key:   note.Key.String(),
				Value: note.Text,
			})
		}
	}

	inv.Order = newOrder(invoice.Ordering)
	inv.TransactionConditions = newTransactionConditions(invoice.Ordering)

	if len(invoice.Preceding) > 0 {
		if invoice.Preceding[0].Reason != "" {
			inv.CorrectionReason = invoice.Preceding[0].Reason
		}
		if invoice.Preceding[0].Ext.Has(favat.ExtKeyEffectiveDate) {
			inv.CorrectionType = invoice.Preceding[0].Ext.Get(favat.ExtKeyEffectiveDate).String()
		}
		for _, prc := range invoice.Preceding {
			inv.CorrectedInv = append(inv.CorrectedInv, NewCorrectedInv(prc))
		}
	}

	return inv
}

// newTransactionConditions builds TransactionConditions from the invoice's
// ordering data. Currently reads order references from purchases; in the
// future this will also handle contracts and other document types.
func newTransactionConditions(ordering *bill.Ordering) *TransactionConditions {
	if ordering == nil {
		return nil
	}

	var tc TransactionConditions

	// Add order references from purchases
	for _, ref := range ordering.Purchases {
		or := &OrderRef{
			Number: ref.Code.String(),
		}
		if ref.IssueDate != nil {
			or.Date = ref.IssueDate.String()
		}
		tc.Orders = append(tc.Orders, or)
	}

	if len(tc.Orders) == 0 {
		return nil
	}

	return &tc
}

func invoiceNumber(series cbc.Code, code cbc.Code) string {
	if series == "" {
		return code.String()
	}
	return fmt.Sprintf("%s-%s", series, code)
}

func (inv *Inv) setTaxRates(taxes *tax.Total) {
	for _, cat := range taxes.Categories {
		if cat.Code != tax.CategoryVAT {
			continue
		}

		for _, rate := range cat.Rates {
			switch rate.Ext.Get(favat.ExtKeyTaxCategory) {
			case "1": // standard rate
				inv.StandardRateNetSale = rate.Base.String()
				inv.StandardRateTax = rate.Amount.String()
			case "2": // reduced rate
				inv.ReducedRateNetSale = rate.Base.String()
				inv.ReducedRateTax = rate.Amount.String()
			case "3": // super reduced rate
				inv.SuperReducedRateNetSale = rate.Base.String()
				inv.SuperReducedRateTax = rate.Amount.String()
			case "4": // taxi rate
				inv.TaxiRateNetSale = rate.Base.String()
				inv.TaxiRateTax = rate.Amount.String()
			case "5": // OSS rate
				inv.OSSNetSale = rate.Base.String()
				inv.OSSTax = rate.Amount.String()
			case "6.1": // zero tax except intra-community supply
				inv.ZeroTaxExceptIntraCommunityNetSale = rate.Base.String()
			case "6.2": // intra-community supply
				inv.IntraCommunityNetSale = rate.Base.String()
			case "6.3": // export supply
				inv.ExportNetSale = rate.Base.String()
			case "7": // tax exempt supply
				inv.TaxExemptNetSale = rate.Base.String()
			case "8": // outside scope supply
				inv.OutsideScopeNetSale = rate.Base.String()
			case "9": // reverse charge supply
				inv.ReverseChargeNetSale = rate.Base.String()
			case "10": // domestic reverse charge supply
				inv.DomesticReverseChargeNetSale = rate.Base.String()
			case "11": // margin supply
				inv.MarginNetSale = rate.Base.String()
			}
		}
	}
}

// parseInvoiceData converts KSEF invoice data to GOBL invoice fields
func (inv *Inv) parseInvoiceData(goblInv *bill.Invoice) error {
	// Parse invoice type and tags
	invType, tags := parseInvoiceType(inv.InvoiceType)
	goblInv.Type = invType
	if len(tags) > 0 {
		goblInv.Tags = tax.Tags{List: tags}
	}

	// Parse issue date
	if inv.IssueDate != "" {
		date, err := parseDate(inv.IssueDate)
		if err != nil {
			return fmt.Errorf("parsing issue date: %w", err)
		}
		goblInv.IssueDate = date
	}

	goblInv.Code = cbc.Code(inv.SequentialNumber)

	// Parse ordering period
	if inv.Period != nil {
		goblInv.Ordering = &bill.Ordering{
			Period: &cal.Period{},
		}
		if inv.Period.StartDate != "" {
			start, err := parseDate(inv.Period.StartDate)
			if err != nil {
				return fmt.Errorf("parsing period start date: %w", err)
			}
			goblInv.Ordering.Period.Start = start
		}
		if inv.Period.EndDate != "" {
			end, err := parseDate(inv.Period.EndDate)
			if err != nil {
				return fmt.Errorf("parsing period end date: %w", err)
			}
			goblInv.Ordering.Period.End = end
		}
	}

	// Parse annotations to tax extensions
	inv.parseAnnotations(goblInv)

	// Always set rounding to currency as per EN16931 standard
	if goblInv.Tax == nil {
		goblInv.Tax = &bill.Tax{}
	}
	goblInv.Tax.Rounding = tax.RoundingRuleCurrency

	// Parse additional description as notes
	if len(inv.AdditionalDescription) > 0 {
		if goblInv.Notes == nil {
			goblInv.Notes = []*org.Note{}
		}
		for _, desc := range inv.AdditionalDescription {
			note := &org.Note{}
			if validCodeRe.MatchString(desc.Key) {
				note.Code = cbc.Code(desc.Key)
				note.Text = desc.Value
			} else {
				// Key contains characters not valid for cbc.Code;
				// combine key and value into the text field.
				note.Text = desc.Key + ": " + desc.Value
			}
			goblInv.Notes = append(goblInv.Notes, note)
		}
	}

	// Parse corrected invoices (preceding documents for credit notes)
	if len(inv.CorrectedInv) > 0 {
		goblInv.Preceding = []*org.DocumentRef{}
		for _, corr := range inv.CorrectedInv {
			preceding := &org.DocumentRef{}

			if corr.SequentialNumber != "" {
				preceding.Code = cbc.Code(corr.SequentialNumber)
			}
			if corr.IssueDate != "" {
				date, err := parseDate(corr.IssueDate)
				if err != nil {
					return fmt.Errorf("parsing corrected invoice date: %w", err)
				}
				preceding.IssueDate = &date
			}
			if inv.CorrectionReason != "" {
				preceding.Reason = inv.CorrectionReason
			}
			if inv.CorrectionType != "" {
				preceding.Ext = tax.Extensions{
					favat.ExtKeyEffectiveDate: cbc.Code(inv.CorrectionType),
				}
			}
			if corr.KsefNumberPresent == 1 && corr.KsefNumber != "" {
				preceding.Stamps = []*head.Stamp{
					{
						Provider: favat.StampKSEFNumber,
						Value:    corr.KsefNumber,
					},
				}
			}

			goblInv.Preceding = append(goblInv.Preceding, preceding)
		}
	}

	return nil
}

// parseDate parses a date string in YYYY-MM-DD format
func parseDate(dateStr string) (cal.Date, error) {
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return cal.Date{}, fmt.Errorf("invalid date format: %w", err)
	}
	return cal.DateOf(t), nil
}

// parseInvoiceType converts KSEF invoice type to GOBL type and tags.
func parseInvoiceType(invType string) (cbc.Key, []cbc.Key) {
	tags := []cbc.Key{}

	switch invType {
	case "VAT":
		return bill.InvoiceTypeStandard, tags
	case "ZAL":
		tags = append(tags, tax.TagPartial)
		return bill.InvoiceTypeStandard, tags
	case "ROZ":
		tags = append(tags, favat.TagSettlement)
		return bill.InvoiceTypeStandard, tags
	case "UPR":
		tags = append(tags, tax.TagSimplified)
		return bill.InvoiceTypeStandard, tags
	case "KOR":
		return bill.InvoiceTypeCreditNote, tags
	case "KOR_ZAL":
		tags = append(tags, tax.TagPartial)
		return bill.InvoiceTypeCreditNote, tags
	case "KOR_ROZ":
		tags = append(tags, favat.TagSettlement)
		return bill.InvoiceTypeCreditNote, tags
	default:
		return bill.InvoiceTypeStandard, tags
	}
}

// parseLines converts KSEF lines to GOBL lines.
func (inv *Inv) parseLines(goblInv *bill.Invoice) error {
	// Always use bypass mode: totals are set directly from KSeF fields
	// and Calculate() is never called.
	goblInv.Tags.List = append(goblInv.Tags.List, tax.TagBypass)

	if len(inv.Lines) == 0 {
		return nil
	}

	goblInv.Lines = make([]*bill.Line, 0, len(inv.Lines))

	// NOTE: Mixed P_9A/P_9B invoices are not supported.
	//
	// KSeF allows P_9A (net unit price, art. 106e(2)-(3)) and P_9B (gross unit
	// price, art. 106e(7)-(8)) per line, with no schema-level mutual exclusivity.
	// In practice the Polish VAT act frames this as an invoice-wide choice, but a
	// document with some lines using P_9A and others using P_9B is technically
	// valid XML.
	//
	// The current approach sets PricesInclude = VAT at the invoice level when any
	// line uses gross pricing (P_9B without P_9A). This works for pure-gross
	// invoices but would produce wrong totals on a mixed invoice: GOBL would strip
	// VAT from the already-net prices on P_9A lines.
	//
	// If a mixed invoice is ever encountered, the fix is to drop PricesInclude and
	// instead convert gross→net explicitly inside Line.ToGOBL() using the line's
	// own VAT rate (P_12): net = gross / (1 + rate). This works uniformly for
	// pure-net, pure-gross, and mixed invoices without any invoice-level flag.
	isCreditNote := goblInv.Type == bill.InvoiceTypeCreditNote
	hasGrossPricing := false
	for i, ksefLine := range inv.Lines {
		line, err := ksefLine.ToGOBL()
		if err != nil {
			return fmt.Errorf("parsing line %d: %w", ksefLine.LineNumber, err)
		}
		line.Index = i + 1 // normally set by Calculate(), but bypass skips it

		// For credit notes: invert non-StanPrzed lines to positive GOBL convention.
		// StanPrzed lines represent the before-correction state and stay as-is
		// (they are the amounts being credited back to the customer).
		if isCreditNote {
			if ksefLine.BeforeCorrectionMarker == 1 {
				line.Notes = append(line.Notes, &org.Note{
					Text: "Before correction (stan przed korektą)",
				})
			} else {
				line.Quantity = line.Quantity.Invert()
				for _, d := range line.Discounts {
					d.Amount = d.Amount.Invert()
				}
				if line.Total != nil {
					inverted := line.Total.Invert()
					line.Total = &inverted
				}
			}
		}

		goblInv.Lines = append(goblInv.Lines, line)
		if ksefLine.NetUnitPrice == "" && ksefLine.GrossUnitPrice != "" {
			hasGrossPricing = true
		}
	}

	if hasGrossPricing {
		if goblInv.Tax == nil {
			goblInv.Tax = &bill.Tax{}
		}
		goblInv.Tax.PricesInclude = tax.CategoryVAT
	}

	return nil
}

// parseTotals builds bill.Totals directly from the invoice-level
// tax summary fields (P_13_X / P_14_X / P_15). The invoice must have the
// bypass tag set so that Calculate() does not overwrite these totals.
func (inv *Inv) parseTotals(goblInv *bill.Invoice) error {
	// Each entry maps a P_13_X / P_14_X pair to a tax category.
	// Well-known percentages are set when the rate is unambiguous;
	// otherwise only the amounts are included.
	type taxEntry struct {
		net      string
		tax      string
		category cbc.Code
		key      cbc.Key
		percent  *num.Percentage // nil when rate is unknown or exempt
	}

	pct23 := num.MakePercentage(230, 3)
	pct8 := num.MakePercentage(80, 3)
	pct5 := num.MakePercentage(50, 3)
	pct0 := num.MakePercentage(0, 3)

	entries := []taxEntry{
		{inv.StandardRateNetSale, inv.StandardRateTax, "1", tax.KeyStandard, &pct23},
		{inv.ReducedRateNetSale, inv.ReducedRateTax, "2", tax.KeyStandard, &pct8},
		{inv.SuperReducedRateNetSale, inv.SuperReducedRateTax, "3", tax.KeyStandard, &pct5},
		{inv.TaxiRateNetSale, inv.TaxiRateTax, "4", tax.KeyStandard, nil},
		{inv.OSSNetSale, inv.OSSTax, "5", tax.KeyStandard, nil},
		{inv.ZeroTaxExceptIntraCommunityNetSale, "", "6.1", tax.KeyZero, &pct0},
		{inv.IntraCommunityNetSale, "", "6.2", tax.KeyIntraCommunity, &pct0},
		{inv.ExportNetSale, "", "6.3", tax.KeyExport, &pct0},
		{inv.TaxExemptNetSale, "", "7", tax.KeyExempt, nil},
		{inv.OutsideScopeNetSale, "", "8", tax.KeyOutsideScope, nil},
		{inv.ReverseChargeNetSale, "", "9", tax.KeyReverseCharge, nil},
		{inv.DomesticReverseChargeNetSale, "", "10", tax.KeyReverseCharge, nil},
	}

	// For credit notes (e.g. KOR_ZAL), KSeF P_13/P_14/P_15 values are negative.
	// GOBL expects positive totals, so we invert them. Invoice.Invert() cannot
	// be used with the bypass tag, so we invert each amount individually.
	isCreditNote := goblInv.Type == bill.InvoiceTypeCreditNote

	var rates []*tax.RateTotal
	var netSum, taxSum num.Amount

	for _, e := range entries {
		if e.net == "" {
			continue
		}

		netAmt, err := parseAmount(e.net)
		if err != nil {
			return fmt.Errorf("parsing prepayment net for category %s: %w", e.category, err)
		}
		if isCreditNote {
			netAmt = netAmt.Invert()
		}

		rt := &tax.RateTotal{
			Key:     e.key,
			Base:    netAmt,
			Percent: e.percent,
			Ext: tax.Extensions{
				favat.ExtKeyTaxCategory: e.category,
			},
		}

		if e.tax != "" {
			taxAmt, err := parseAmount(e.tax)
			if err != nil {
				return fmt.Errorf("parsing prepayment tax for category %s: %w", e.category, err)
			}
			if isCreditNote {
				taxAmt = taxAmt.Invert()
			}
			rt.Amount = taxAmt
			taxSum = taxSum.MatchPrecision(taxAmt)
			taxSum = taxSum.Add(taxAmt)
		}

		netSum = netSum.MatchPrecision(netAmt)
		netSum = netSum.Add(netAmt)
		rates = append(rates, rt)
	}

	totals := &bill.Totals{
		Sum:   netSum,
		Total: netSum,
		Tax:   taxSum,
	}

	if len(rates) > 0 {
		totals.Taxes = &tax.Total{
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

	totals.TotalWithTax = netSum.Add(taxSum)

	if inv.TotalAmountDue != "" {
		payable, err := parseAmount(inv.TotalAmountDue)
		if err != nil {
			return fmt.Errorf("parsing total amount due: %w", err)
		}
		if isCreditNote {
			payable = payable.Invert()
		}
		totals.Payable = payable
	} else {
		totals.Payable = totals.TotalWithTax
	}

	// Sum advances from payment details
	if goblInv.Payment != nil && len(goblInv.Payment.Advances) > 0 {
		var advanceSum num.Amount
		for _, adv := range goblInv.Payment.Advances {
			advanceSum = advanceSum.MatchPrecision(adv.Amount)
			advanceSum = advanceSum.Add(adv.Amount)
		}
		totals.Advances = &advanceSum
		due := totals.Payable.Subtract(advanceSum)
		totals.Due = &due
	}

	goblInv.Totals = totals
	return nil
}
