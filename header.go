// Package ksef implements the conversion from GOBL to FA_VAT XML
package ksef

import (
	"cloud.google.com/go/civil"
	"github.com/invopop/gobl/bill"
	"github.com/invopop/gobl/cal"
)

// KSEF schema constants
const (
	kodSystemowy      = "FA (2)"
	kodFormularza     = "FA"
	wersjaSchemy      = "1-0E"
	wariantFormularza = 2
	systemInfo        = "GOBL.KSEF"
)

type Naglowek struct {
	KodFormularza     *KodFormularza `xml:"KodFormularza"`
	WariantFormularza int            `xml:"WariantFormularza"`
	DataWytworzeniaFa string         `xml:"DataWytworzeniaFa"`
	SystemInfo        string         `xml:"SystemInfo"`
}

type KodFormularza struct {
	KodSystemowy  string `xml:"kodSystemowy,attr"`
	WersjaSchemy  string `xml:"wersjaSchemy,attr"`
	KodFormularza string `xml:",chardata"`
}

func NewNaglowek(inv *bill.Invoice) *Naglowek {
	date := formatIssueDate(inv.IssueDate)

	naglowek := &Naglowek{
		KodFormularza: &KodFormularza{
			KodSystemowy:  kodSystemowy,
			WersjaSchemy:  wersjaSchemy,
			KodFormularza: kodFormularza,
		},
		WariantFormularza: wariantFormularza,
		DataWytworzeniaFa: date,
		SystemInfo:        systemInfo,
	}

	return naglowek
}

func formatIssueDate(date cal.Date) string {
	dateTime := civil.DateTime{Date: date.Date, Time: civil.Time{}}
	return dateTime.String() + "Z"
}
