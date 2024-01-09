package ksef

import (
	"fmt"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cbc"
	"github.com/invopop/gobl/regimes/pl"
	"github.com/invopop/gobl/tax"
)

type AdvancePayment struct {
	PaymentAmount string `xml:"KwotaZaplatyCzesciowej,omitempty"`
	PaymentDate   string `xml:"DataZaplatyCzesciowej,omitempty"`
}

type PartialPayment struct {
	Date        string `xml:"Termin,omitempty"`
	Description string `xml:"TerminOpis,omitempty"`
}

type BankAccount struct {
	AccountNumber         string `xml:"NrRB"`
	SWIFT                 string `xml:"SWIFT,omitempty"`
	BankSelfAccountMarker int    `xml:"RachunekWlasnyBanku,omitempty"` // enum - 1,2,3, not sure what exactly they mean
	BankName              string `xml:"NazwaBanku,omitempty"`
	AccountDescription    string `xml:"OpisRachunku,omitempty"`
}

type Discount struct { // TODO
	Conditions string `xml:"WarunkiSkonta,omitempty"`
	Amount     string `xml:"WysokoscSkonta,omitempty"`
}

type Payment struct {
	PaidMarker             string            `xml:"Zaplacono,omitempty"`
	PaymentDate            string            `xml:"DataZaplaty,omitempty"`
	PartiallyPaidMarker    string            `xml:"ZnacznikZaplatyCzesciowej,omitempty"`
	AdvancePayments        []*AdvancePayment `xml:"ZaplataCzesciowa,omitempty"`
	PartialPayments        []*PartialPayment `xml:"TerminPlatnosci,omitempty"`
	PaymentMean            string            `xml:"FormaPlatnosci,omitempty"`
	OtherPaymentMeanMarker string            `xml:"PlatnoscInna,omitempty"`
	OtherPaymentMean       string            `xml:"OpisPlatnosci,omitempty"`
	BankAccounts           []*BankAccount    `xml:"RachunekBankowy,omitempty"`
	FactorBankAccounts     []*BankAccount    `xml:"RachunekBankowyFaktora,omitempty"` // not sure if supported by gobl
	Discount               *Discount         `xml:"Skonto,omitempty"`                 // it's some special discount for early payments
}

func NewPayment(inv *bill.Invoice) *Payment {

	var payment = &Payment{
		PartialPayments: []*PartialPayment{},
		AdvancePayments: []*AdvancePayment{},
	}

	PaymentMeansCode, err := findPaymentMeansCode(inv.Payment.Instructions.Key)

	if err != nil {
		payment.OtherPaymentMeanMarker = "1"
		payment.OtherPaymentMean = inv.Payment.Instructions.Key.String()
	} else {
		payment.PaymentMean = PaymentMeansCode
	}

	if terms := inv.Payment.Terms; terms != nil {
		for _, dueDate := range inv.Payment.Terms.DueDates {
			payment.PartialPayments = append(payment.PartialPayments, &PartialPayment{
				Date:        dueDate.Date.String(),
				Description: dueDate.Amount.String(),
			})
		}
	}

	if advances := inv.Payment.Advances; advances != nil {
		if len(advances) > 1 || len(inv.Payment.Terms.DueDates) > 0 {
			payment.PartiallyPaidMarker = "1"
			for _, advance := range inv.Payment.Advances {
				payment.AdvancePayments = append(payment.AdvancePayments, &AdvancePayment{
					PaymentAmount: advance.Amount.String(),
					PaymentDate:   advance.Date.String(),
				})
			}
		} else {
			if len(advances) == 1 {
				payment.PaidMarker = "1"
				payment.PaymentDate = advances[0].Date.String()
			}
		}
	}

	payment.BankAccounts = []*BankAccount{}
	payment.FactorBankAccounts = []*BankAccount{}

	for _, account := range inv.Payment.Instructions.CreditTransfer {
		payment.BankAccounts = append(payment.BankAccounts, &BankAccount{
			AccountNumber: account.Number,
			SWIFT:         account.BIC,
			BankName:      account.Name,
		})
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

func findPaymentKeyDefinition(key cbc.Key) *tax.KeyDefinition {
	for _, keyDef := range regime.PaymentMeansKeys {
		if key == keyDef.Key {
			return keyDef
		}
	}
	return nil
}
