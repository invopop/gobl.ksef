package ksef

import (
	"github.com/invopop/gobl/head"
	"github.com/invopop/gobl/org"
)

// CorrectedInv defines the XML structure for KSeF correction invoice
type CorrectedInv struct {
	IssueDate           string `xml:"DataWystFaKorygowanej,omitempty"`
	SequentialNumber    string `xml:"NrFaKorygowanej,omitempty"`
	KsefNumberPresent   int    `xml:"NrKSeF,omitempty"`
	NoKsefNumberPresent int    `xml:"NrKSeFN,omitempty"`
	KsefNumber          string `xml:"NrKSeFFaKorygowanej,omitempty"`
}

// NewCorrectedInv gets credit note data from GOBL invoice
func NewCorrectedInv(prc *org.DocumentRef) *CorrectedInv {
	inv := &CorrectedInv{
		SequentialNumber: invoiceNumber(prc.Series, prc.Code),
	}

	if prc.IssueDate != nil {
		inv.IssueDate = prc.IssueDate.String()
	}

	if id := findStamp(prc.Stamps, "ksef-id"); id != -1 {
		inv.KsefNumberPresent = 1
		inv.KsefNumber = prc.Stamps[id].Value
	} else {
		inv.NoKsefNumberPresent = 1
	}

	return inv
}

func findStamp(a []*head.Stamp, x string) int {
	for i, n := range a {
		if x == string(n.Provider) {
			return i
		}
	}
	return -1
}
