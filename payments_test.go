package ksef_test

import (
	"testing"
	"time"

	ksef "github.com/invopop/gobl.ksef"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
	"github.com/invopop/gobl/num"
	"github.com/invopop/gobl/pay"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPayment(t *testing.T) {
	t.Run("should return nil when no payment data passed", func(t *testing.T) {
		pay := ksef.NewPayment(nil, nil)
		assert.Nil(t, pay)
	})

	t.Run("should return payment if there are payment instructions", func(t *testing.T) {
		payment := &bill.Payment{
			Instructions: &pay.Instructions{},
		}
		totals := &bill.Totals{}

		pay := ksef.NewPayment(payment, totals)
		result := &ksef.Payment{
			PaidMarker:             "",
			PaymentDate:            "",
			PartiallyPaidMarker:    "",
			AdvancePayments:        []*ksef.AdvancePayment{},
			DueDates:               []*ksef.DueDate{},
			PaymentMean:            "",
			OtherPaymentMeanMarker: "1",
			OtherPaymentMean:       "",
			BankAccounts:           []*ksef.BankAccount{},
			FactorBankAccounts:     []*ksef.BankAccount{},
			Discount:               (*ksef.Discount)(nil),
		}

		assert.Equal(t, result, pay)
	})

	t.Run("should return set payment method from payment instructions", func(t *testing.T) {
		payment := &bill.Payment{
			Instructions: &pay.Instructions{
				Key: "credit-transfer",
			},
		}
		totals := &bill.Totals{}
		pay := ksef.NewPayment(payment, totals)
		result := &ksef.Payment{
			PaidMarker:             "",
			PaymentDate:            "",
			PartiallyPaidMarker:    "",
			AdvancePayments:        []*ksef.AdvancePayment{},
			DueDates:               []*ksef.DueDate{},
			PaymentMean:            "6",
			OtherPaymentMeanMarker: "",
			OtherPaymentMean:       "",
			BankAccounts:           []*ksef.BankAccount{},
			FactorBankAccounts:     []*ksef.BankAccount{},
			Discount:               (*ksef.Discount)(nil),
		}

		assert.Equal(t, result, pay)
	})

	t.Run("should set payment terms", func(t *testing.T) {
		x := time.Date(2023, time.July, 28, 0, 0, 0, 0, time.UTC)
		d := cal.DateOf(x)
		num, err := num.AmountFromString("245.890")
		require.NoError(t, err)

		payment := &bill.Payment{
			Terms: &pay.Terms{
				DueDates: []*pay.DueDate{{Date: &d, Amount: num}},
			},
		}
		totals := &bill.Totals{}
		pay := ksef.NewPayment(payment, totals)
		result := &ksef.Payment{
			PaidMarker:             "",
			PaymentDate:            "",
			PartiallyPaidMarker:    "",
			AdvancePayments:        []*ksef.AdvancePayment{},
			DueDates:               []*ksef.DueDate{{Date: d.String(), Description: num.String()}},
			PaymentMean:            "",
			OtherPaymentMeanMarker: "",
			OtherPaymentMean:       "",
			BankAccounts:           []*ksef.BankAccount(nil),
			FactorBankAccounts:     []*ksef.BankAccount(nil),
			Discount:               (*ksef.Discount)(nil),
		}

		assert.Equal(t, result, pay)
	})

	t.Run("advances should set paid marker and date", func(t *testing.T) {
		x := time.Date(2023, time.July, 28, 0, 0, 0, 0, time.UTC)
		d := cal.DateOf(x)
		firstNum, err := num.AmountFromString("245.890")
		require.NoError(t, err)
		zero, err := num.AmountFromString("0")
		require.NoError(t, err)

		payment := &bill.Payment{
			Advances: []*pay.Advance{{Date: &d, Amount: firstNum}},
		}
		totals := &bill.Totals{
			Due:      &zero,
			Advances: &firstNum,
		}
		pay := ksef.NewPayment(payment, totals)
		result := &ksef.Payment{
			PaidMarker:             "1",
			PaymentDate:            d.String(),
			PartiallyPaidMarker:    "",
			AdvancePayments:        []*ksef.AdvancePayment{},
			DueDates:               []*ksef.DueDate{},
			PaymentMean:            "",
			OtherPaymentMeanMarker: "",
			OtherPaymentMean:       "",
			BankAccounts:           []*ksef.BankAccount(nil),
			FactorBankAccounts:     []*ksef.BankAccount(nil),
			Discount:               (*ksef.Discount)(nil),
		}

		assert.Equal(t, result, pay)
	})

	t.Run("multiple advances sets partially paid marker and advance fields", func(t *testing.T) {
		x := time.Date(2023, time.July, 28, 0, 0, 0, 0, time.UTC)
		d := cal.DateOf(x)
		firstNum, err := num.AmountFromString("245.890")
		require.NoError(t, err)
		secondNum, err := num.AmountFromString("45.990")
		require.NoError(t, err)

		payment := &bill.Payment{
			Advances: []*pay.Advance{{Date: &d, Amount: firstNum}, {Date: &d, Amount: secondNum}},
		}
		totals := &bill.Totals{
			Due:      &secondNum,
			Advances: &firstNum,
		}
		pay := ksef.NewPayment(payment, totals)
		result := &ksef.Payment{
			PaidMarker:             "",
			PaymentDate:            "",
			PartiallyPaidMarker:    "1",
			AdvancePayments:        []*ksef.AdvancePayment{{PaymentAmount: firstNum.String(), PaymentDate: d.String()}, {PaymentAmount: secondNum.String(), PaymentDate: d.String()}},
			DueDates:               []*ksef.DueDate{},
			PaymentMean:            "",
			OtherPaymentMeanMarker: "",
			OtherPaymentMean:       "",
			BankAccounts:           []*ksef.BankAccount(nil),
			FactorBankAccounts:     []*ksef.BankAccount(nil),
			Discount:               (*ksef.Discount)(nil),
		}

		assert.Equal(t, result, pay)
	})
}
