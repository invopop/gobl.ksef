package ksef

import (
	"fmt"

	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/num"
)

// RoundingError represents a rounding discrepancy that exceeded acceptable thresholds.
// When this error is returned, the inv parameter passed to AdjustRounding has had the
// rounding adjustment applied, allowing callers to use it despite the warning.
type RoundingError struct {
	Diff       num.Amount
	MaxAllowed num.Amount
}

func (e *RoundingError) Error() string {
	return fmt.Sprintf("rounding error in totals too high: %s (max allowed: %s)",
		e.Diff.String(), e.MaxAllowed.String())
}

// AdjustRounding checks and, if needed, adjusts the rounding in the GOBL invoice to match the
// KSEF total amount. KSEF calculates totals by rounding each line and then summing,
// which can lead to a mismatch with the total amount in GOBL.
func AdjustRounding(inv *bill.Invoice, ksefTotalDue string) error {
	// First calculate the GOBL totals
	if err := inv.Calculate(); err != nil {
		return err
	}

	if inv.Totals == nil {
		return fmt.Errorf("invoice totals are nil after calculation")
	}

	// Parse the KSEF total amount
	expectedTotal, err := parseAmount(ksefTotalDue)
	if err != nil {
		return fmt.Errorf("parsing KSEF total amount: %w", err)
	}

	// Calculate the difference between the expected and the calculated totals
	var calculatedTotal num.Amount
	if inv.Totals.Due != nil && !inv.Totals.Due.IsZero() {
		calculatedTotal = *inv.Totals.Due
	} else {
		calculatedTotal = inv.Totals.Payable
	}

	diff := expectedTotal.Subtract(calculatedTotal)
	if diff.IsZero() {
		// No difference. No adjustment needed
		return nil
	}

	// Check if the difference can be attributed to rounding
	maxErr := MaxRoundingError(inv)
	if diff.Abs().Compare(maxErr) == 1 {
		// Too much difference. Apply the adjustment anyway and return a warning
		inv.Totals.Rounding = &diff
		return &RoundingError{
			Diff:       diff,
			MaxAllowed: maxErr,
		}
	}

	// Apply the rounding adjustment
	inv.Totals.Rounding = &diff

	return nil
}

// MaxRoundingError returns the maximum error that can be attributed to rounding in an invoice.
// It calculates 1 of the smallest subunit of the currency per line.
func MaxRoundingError(inv *bill.Invoice) num.Amount {
	// 1 of the smallest subunit of the currency per line
	subunits := inv.Currency.Def().Subunits
	return num.MakeAmount(int64(len(inv.Lines)), subunits)
}
