package ksef

import (
	"fmt"
	"time"

	"github.com/invopop/gobl/addons/pl/favat"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/org"
	"github.com/invopop/gobl/tax"
)

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
	TransactionConditions              *TransactionConditions       `xml:"WarunkiTransakcji,omitempty"`
	Payment                            *Payment                     `xml:"Platnosc,omitempty"`
	Order                              *Order                       `xml:"Zamowienie,omitempty"` // for ZAL and KOR_ZAL types
}

type InvoicePeriod struct {
	StartDate string `xml:"P_6_Od,omitempty"`
	EndDate   string `xml:"P_6_Do,omitempty"`
}

// Annotations defines the XML structure for KSeF annotations
type Annotations struct {
	CashAccounting                      string             `xml:"P_16"`
	SelfBilling                         string             `xml:"P_17"`
	ReverseCharge                       string             `xml:"P_18"`
	SplitPaymentMechanism               string             `xml:"P_18A"`
	TaxExemption                        *TaxExemption      `xml:"Zwolnienie,omitempty"`
	NewTransportMeans                   *NewTransportMeans `xml:"NoweSrodkiTransportu,omitempty"`
	SimplifiedProcedureBySecondTaxpayer string             `xml:"P_23"`
	MarginScheme                        *MarginScheme      `xml:"PMarzy,omitempty"`
}

// TaxExemption defines the XML structure for tax exemption details
type TaxExemption struct {
	Marker           string `xml:"P_19,omitempty"`
	PolishLawBasis   string `xml:"P_19A,omitempty"`
	EUDirectiveBasis string `xml:"P_19B,omitempty"`
	OtherLegalBasis  string `xml:"P_19C,omitempty"`
	NoExemption      string `xml:"P_19N,omitempty"`
}

// NewTransportMeans defines the XML structure for new means of transport
type NewTransportMeans struct {
	Marker                 int                      `xml:"P_22,omitempty"`
	Art42Obligation        string                   `xml:"P_42_5,omitempty"`
	NewTransportMeansItems []*NewTransportMeansItem `xml:"NowySrodekTransportu,omitempty"`
	NoNewTransportMeans    string                   `xml:"P_22N,omitempty"`
}

// NewTransportMeansItem defines details for a single new transport means item
type NewTransportMeansItem struct {
	FirstUseDate       string `xml:"P_22A"`
	LineNumber         int    `xml:"P_NrWierszaNST"`
	Brand              string `xml:"P_22BMK,omitempty"`
	Model              string `xml:"P_22BMD,omitempty"`
	Color              string `xml:"P_22BK,omitempty"`
	RegistrationNumber string `xml:"P_22BNR,omitempty"`
	ProductionYear     string `xml:"P_22BRP,omitempty"`
	// For land vehicles
	Mileage       string `xml:"P_22B,omitempty"`
	VIN           string `xml:"P_22B1,omitempty"`
	BodyNumber    string `xml:"P_22B2,omitempty"`
	ChassisNumber string `xml:"P_22B3,omitempty"`
	FrameNumber   string `xml:"P_22B4,omitempty"`
	VehicleType   string `xml:"P_22BT,omitempty"`
	// For watercraft
	OperatingHoursWater string `xml:"P_22C,omitempty"`
	HullNumber          string `xml:"P_22C1,omitempty"`
	// For aircraft
	OperatingHoursAir string `xml:"P_22D,omitempty"`
	FactoryNumber     string `xml:"P_22D1,omitempty"`
}

// MarginScheme defines the XML structure for margin scheme
type MarginScheme struct {
	Marker                        string `xml:"P_PMarzy,omitempty"`
	TravelAgencyMargin            string `xml:"P_PMarzy_2,omitempty"`
	UsedGoodsMargin               string `xml:"P_PMarzy_3_1,omitempty"`
	ArtWorksMargin                string `xml:"P_PMarzy_3_2,omitempty"`
	CollectiblesAndAntiquesMargin string `xml:"P_PMarzy_3_3,omitempty"`
	NoMarginScheme                string `xml:"P_PMarzyN,omitempty"`
}

// AdditionalDescriptionLine defines the XML structure for KSeF additional description line (`DodatkowyOpis`)
type AdditionalDescriptionLine struct {
	LineNumber string `xml:"NrWiersza,omitempty"`
	Key        string `xml:"Klucz"`
	Value      string `xml:"Wartosc"`
}

// Order defines the XML structure for KSeF "Zamowienie" (order) field, required for ZAL and KOR_ZAL types
type Order struct {
	OrderAmount string       `xml:"WartoscZamowienia"`
	LineItems   []*OrderLine `xml:"ZamowienieWiersz,omitempty"`
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
	Contracts         []*Contract  `xml:"Umowy>Umowa,omitempty"`
	Orders            []*OrderRef  `xml:"Zamowienia>Zamowienie,omitempty"`
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

func invoiceNumber(series cbc.Code, code cbc.Code) string {
	if series == "" {
		return code.String()
	}
	return fmt.Sprintf("%s-%s", series, code)
}

func newInvoicePeriod(ordering *bill.Ordering) *InvoicePeriod {
	if ordering == nil || ordering.Period == nil {
		return nil
	}

	return &InvoicePeriod{
		StartDate: ordering.Period.Start.String(),
		EndDate:   ordering.Period.End.String(),
	}
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
	if inv.Annotations != nil {
		goblInv.Tax = &bill.Tax{
			Ext: make(tax.Extensions),
		}

		// Set invoice type extension
		if inv.InvoiceType != "" {
			goblInv.Tax.Ext[favat.ExtKeyInvoiceType] = cbc.Code(inv.InvoiceType)
		}

		// Cash accounting
		if inv.Annotations.CashAccounting == "1" {
			goblInv.Tax.Ext[favat.ExtKeyCashAccounting] = "1"
		}

		// Self billing
		if inv.Annotations.SelfBilling == "1" {
			goblInv.Tax.Ext[favat.ExtKeySelfBilling] = "1"
			goblInv.Tags.List = append(goblInv.Tags.List, tax.TagSelfBilled)
		}

		// Reverse charge
		if inv.Annotations.ReverseCharge == "1" {
			goblInv.Tax.Ext[favat.ExtKeyReverseCharge] = "1"
			goblInv.Tags.List = append(goblInv.Tags.List, tax.TagReverseCharge)
		}

		// Split payment
		if inv.Annotations.SplitPaymentMechanism == "1" {
			goblInv.Tax.Ext[favat.ExtKeySplitPayment] = "1"
		}

		// Tax exemption
		if inv.Annotations.TaxExemption != nil && inv.Annotations.TaxExemption.Marker == "1" {
			// Determine exemption code
			var exemptionCode string
			var exemptionText string

			if inv.Annotations.TaxExemption.PolishLawBasis != "" {
				exemptionCode = "A"
				exemptionText = inv.Annotations.TaxExemption.PolishLawBasis
			} else if inv.Annotations.TaxExemption.EUDirectiveBasis != "" {
				exemptionCode = "B"
				exemptionText = inv.Annotations.TaxExemption.EUDirectiveBasis
			} else if inv.Annotations.TaxExemption.OtherLegalBasis != "" {
				exemptionCode = "C"
				exemptionText = inv.Annotations.TaxExemption.OtherLegalBasis
			}

			if exemptionCode != "" {
				goblInv.Tax.Ext[favat.ExtKeyExemption] = cbc.Code(exemptionCode)

				// Add note with exemption text
				if goblInv.Notes == nil {
					goblInv.Notes = []*org.Note{}
				}
				goblInv.Notes = append(goblInv.Notes, &org.Note{
					Key:  org.NoteKeyLegal,
					Code: cbc.Code(exemptionCode),
					Src:  favat.ExtKeyExemption,
					Text: exemptionText,
				})
			}
		}

		// Margin scheme
		if inv.Annotations.MarginScheme != nil && inv.Annotations.MarginScheme.Marker == "1" {
			var marginCode string
			if inv.Annotations.MarginScheme.TravelAgencyMargin == "1" {
				marginCode = "2"
			} else if inv.Annotations.MarginScheme.UsedGoodsMargin == "1" {
				marginCode = "3.1"
			} else if inv.Annotations.MarginScheme.ArtWorksMargin == "1" {
				marginCode = "3.2"
			} else if inv.Annotations.MarginScheme.CollectiblesAndAntiquesMargin == "1" {
				marginCode = "3.3"
			}

			if marginCode != "" {
				goblInv.Tax.Ext[favat.ExtKeyMarginScheme] = cbc.Code(marginCode)
			}
		}
	}

	// Parse additional description as notes
	if len(inv.AdditionalDescription) > 0 {
		if goblInv.Notes == nil {
			goblInv.Notes = []*org.Note{}
		}
		for _, desc := range inv.AdditionalDescription {
			goblInv.Notes = append(goblInv.Notes, &org.Note{
				Code: cbc.Code(desc.Key),
				Text: desc.Value,
			})
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
	hasGrossPricing := false
	for _, ksefLine := range inv.Lines {
		line, err := ksefLine.ToGOBL()
		if err != nil {
			return fmt.Errorf("parsing line %d: %w", ksefLine.LineNumber, err)
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

// newAnnotations sets annotations data
func newAnnotations(invoice *bill.Invoice) *Annotations {
	// default values for the most common case,
	// For fields P_16 to P_18 and field P_23 2 means "no", 1 means "yes".
	// for others 1 means "yes", no value means "no"
	Annotations := &Annotations{
		CashAccounting:        "2",
		SelfBilling:           "2",
		ReverseCharge:         "2",
		SplitPaymentMechanism: "2",
		TaxExemption: &TaxExemption{
			NoExemption: "1",
		},
		NewTransportMeans: &NewTransportMeans{
			NoNewTransportMeans: "1",
		},
		SimplifiedProcedureBySecondTaxpayer: "2",
		MarginScheme: &MarginScheme{
			NoMarginScheme: "1",
		},
	}

	if invoice.Tax == nil {
		return Annotations
	}

	if invoice.Tax.Ext.Get(favat.ExtKeyCashAccounting) == "1" {
		Annotations.CashAccounting = "1"
	}

	if invoice.Tax.Ext.Get(favat.ExtKeySelfBilling) == "1" {
		Annotations.SelfBilling = "1"
	}

	if invoice.Tax.Ext.Get(favat.ExtKeyReverseCharge) == "1" {
		Annotations.ReverseCharge = "1"
	}

	if invoice.Tax.Ext.Get(favat.ExtKeySplitPayment) == "1" {
		Annotations.SplitPaymentMechanism = "1"
	}

	if invoice.Tax.Ext.Get(favat.ExtKeyExemption) != "" {
		// Find the note in notes with key legal
		Annotations.TaxExemption = &TaxExemption{
			Marker: "1",
		}
		for _, note := range invoice.Notes {
			if note.Key == org.NoteKeyLegal && note.Src == favat.ExtKeyExemption {
				switch invoice.Tax.Ext.Get(favat.ExtKeyExemption) {
				case "A": // polish law basis
					Annotations.TaxExemption.PolishLawBasis = note.Text
				case "B": // EU directive basis
					Annotations.TaxExemption.EUDirectiveBasis = note.Text
				case "C": // other legal basis
					Annotations.TaxExemption.OtherLegalBasis = note.Text
				}
				break
			}
		}
	}

	if invoice.Tax.Ext.Get(favat.ExtKeyMarginScheme) != "" {
		Annotations.MarginScheme = &MarginScheme{
			Marker: "1",
		}
		switch invoice.Tax.Ext.Get(favat.ExtKeyMarginScheme) {
		case "2": // travel agency margin scheme
			Annotations.MarginScheme.TravelAgencyMargin = "1"
		case "3.1": // used goods margin scheme
			Annotations.MarginScheme.UsedGoodsMargin = "1"
		case "3.2": // art works margin scheme
			Annotations.MarginScheme.ArtWorksMargin = "1"
		case "3.3": // collectibles and antiques margin scheme
			Annotations.MarginScheme.CollectiblesAndAntiquesMargin = "1"
		}
	}

	return Annotations
}
