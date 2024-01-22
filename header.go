// Package ksef implements the conversion from GOBL to FA_VAT XML
package ksef

import (
	"cloud.google.com/go/civil"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
)

// KSEF schema constants
const (
	systemCode    = "FA (2)"
	formCode      = "FA"
	schemaVersion = "1-0E"
	formVariant   = 2
	systemInfo    = "GOBL.KSEF"
)

// Header defines the XML structure for KSeF header
type Header struct {
	FormCode     *FormCode `xml:"KodFormularza"`
	FormVariant  int       `xml:"WariantFormularza"`
	CreationDate string    `xml:"DataWytworzeniaFa"`
	SystemInfo   string    `xml:"SystemInfo"`
}

// FormCode defines the XML structure for KSeF schema versioning
type FormCode struct {
	SystemCode    string `xml:"kodSystemowy,attr"`
	SchemaVersion string `xml:"wersjaSchemy,attr"`
	FormCode      string `xml:",chardata"`
}

// NewHeader gets header data from GOBL invoice
func NewHeader(inv *bill.Invoice) *Header {
	date := formatIssueDate(inv.IssueDate)

	header := &Header{
		FormCode: &FormCode{
			SystemCode:    systemCode,
			SchemaVersion: schemaVersion,
			FormCode:      formCode,
		},
		FormVariant:  formVariant,
		CreationDate: date,
		SystemInfo:   systemInfo,
	}

	return header
}

func formatIssueDate(date cal.Date) string {
	dateTime := civil.DateTime{Date: date.Date, Time: civil.Time{}}
	return dateTime.String() + "Z"
}
