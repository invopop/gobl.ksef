package ksef

import (
	"fmt"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/pl"
)

// AdvancePayment defines the XML structure for KSeF advance payments
type AdvancePayment struct {
	PaymentAmount string `xml:"KwotaZaplatyCzesciowej,omitempty"`
	PaymentDate   string `xml:"DataZaplatyCzesciowej,omitempty"`
}

// DueDate defines the XML structure for KSeF due date
type DueDate struct {
	Date string `xml:"Termin,omitempty"`
}

// BankAccount defines the XML structure for KSeF bank accounts
type BankAccount struct {
	AccountNumber         string `xml:"NrRB"`
	SWIFT                 string `xml:"SWIFT,omitempty"`
	BankSelfAccountMarker int    `xml:"RachunekWlasnyBanku,omitempty"` // enum - 1,2,3, not sure what exactly they mean
	BankName              string `xml:"NazwaBanku,omitempty"`
	AccountDescription    string `xml:"OpisRachunku,omitempty"`
}

// Discount defines the XML structure for KSeF discount
type Discount struct { // TODO
	Conditions string `xml:"WarunkiSkonta,omitempty"`
	Amount     string `xml:"WysokoscSkonta,omitempty"`
}

// Payment defines the XML structure for KSeF payment
type Payment struct {
	PaidMarker             string            `xml:"Zaplacono,omitempty"`
	PaymentDate            string            `xml:"DataZaplaty,omitempty"`
	PartiallyPaidMarker    string            `xml:"ZnacznikZaplatyCzesciowej,omitempty"`
	AdvancePayments        []*AdvancePayment `xml:"ZaplataCzesciowa,omitempty"`
	DueDates               []*DueDate        `xml:"TerminPlatnosci,omitempty"`
	PaymentMean            string            `xml:"FormaPlatnosci,omitempty"` // enum: 1 = cash, 2 = card etc. (see KSeF documentation)
	OtherPaymentMeanMarker string            `xml:"PlatnoscInna,omitempty"`
	OtherPaymentMean       string            `xml:"OpisPlatnosci,omitempty"`
	BankAccounts           []*BankAccount    `xml:"RachunekBankowy,omitempty"`
	FactorBankAccounts     []*BankAccount    `xml:"RachunekBankowyFaktora,omitempty"` // not sure if supported by gobl
	Discount               *Discount         `xml:"Skonto,omitempty"`                 // it's some special discount for early payments
}

// NewPayment gets payment data from GOBL invoice
func NewPayment(pay *bill.PaymentDetails, totals *bill.Totals) *Payment {
	if pay == nil {
		return nil
	}

	var payment = &Payment{
		DueDates:        []*DueDate{},
		AdvancePayments: []*AdvancePayment{},
	}

	if instructions := pay.Instructions; instructions != nil {
		PaymentMeansCode, err := findPaymentMeansCode(instructions.Key)

		if err != nil {
			payment.OtherPaymentMeanMarker = "1"
			payment.OtherPaymentMean = instructions.Key.String()
		} else {
			payment.PaymentMean = PaymentMeansCode
		}

		payment.BankAccounts = []*BankAccount{}
		payment.FactorBankAccounts = []*BankAccount{}

		for _, account := range instructions.CreditTransfer {
			payment.BankAccounts = append(payment.BankAccounts, &BankAccount{
				AccountNumber: account.Number,
				SWIFT:         account.BIC,
				BankName:      account.Name,
			})
		}
	}

	if terms := pay.Terms; terms != nil {
		for _, dueDate := range pay.Terms.DueDates {
			payment.DueDates = append(payment.DueDates, &DueDate{
				Date: dueDate.Date.String(),
			})
		}
	}

	// According to FA_VAT v3 schema:
	// If an invoice is paid in full in one payment, PaidMarker should be "1"
	// Otherwise, set PartiallyPaidMarker with the following values:
	// 1 = invoice paid partially
	// 2 = paid in full after partial payments, and the last payment is the final one
	// If the invoice is not paid at all, do not add PaidMarker or PartiallyPaidMarker.

	if advances := pay.Advances; advances != nil {
		if len(advances) == 1 && totals.Due.IsZero() {
			// Invoice already paid in full in one payment
			payment.PaidMarker = "1"
			payment.PaymentDate = advances[len(advances)-1].Date.String()
		} else {
			if totals.Due.IsZero() {
				// Invoice already paid in full in multiple payments
				payment.PartiallyPaidMarker = "2"
			}
			if !totals.Due.IsZero() && len(advances) > 0 {
				// Invoice paid partially
				payment.PartiallyPaidMarker = "1"
			}
			// Otherwise, not paid at all - no markers needed

			for _, advance := range advances {
				payment.AdvancePayments = append(payment.AdvancePayments, &AdvancePayment{
					PaymentAmount: advance.Amount.String(),
					PaymentDate:   advance.Date.String(),
				})
			}
		}
	}

	return payment
}

func findPaymentMeansCode(key cbc.Key) (string, error) {
	keyDef := findPaymentKeyDefinition(key)

	if keyDef == nil {
		return "", fmt.Errorf("FormaPlatnosci Code not found for payment method key '%s'", key)
	}

	code := keyDef.Map[pl.KeyFAVATPaymentType]
	if code == "" {
		return "", fmt.Errorf("FormaPlatnosci Code not found for payment method key '%s'", key)
	}

	return code.String(), nil
}

func findPaymentKeyDefinition(key cbc.Key) *cbc.Definition {
	// TODO in the newest gobl library version it's moved from regime to addon
	// The addon will be at github.com/invopop/gobl/addons/pl/favat

	for _, keyDef := range regime.PaymentMeansKeys {
		if key == keyDef.Key {
			return keyDef
		}
	}
	return nil
}
