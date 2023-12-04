package ksef

import (
	"github.com/invopop/gobl/bill"
)

type Footer struct {
}

func NewFooter(inv *bill.Invoice) *Footer {

	footer := &Footer{}

	return footer
}
