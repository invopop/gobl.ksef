// Package ksef implements the conversion from GOBL to FA_VAT XML
package ksef

import (
	"encoding/xml"
	"fmt"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/bill"
)

const (
	XSINamespace = "http://www.w3.org/2001/XMLSchema-instance"
	XSDNamespace = "http://www.w3.org/2001/XMLSchema"
	XMLNamespace = "http://crd.gov.pl/wzor/2023/06/29/12648/"
)

// Faktura is a pseudo-model for containing the XML document being created
type Faktura struct {
	XSINamespace string    `xml:"xmlns:xsi,attr"`
	XSDNamespace string    `xml:"xmlns:xsd,attr"`
	XMLNamespace string    `xml:"xmlns,attr"`
	Naglowek     *Naglowek `xml:"Naglowek"`
	Podmiot1     *Podmiot1 `xml:"Podmiot1"`
	Podmiot2     *Podmiot2 `xml:"Podmiot2"`
	Fa           *Fa       `xml:"Fa"`
	Stopka       *Stopka   `xml:"Stopka,omitempty"`
}

// NewDocument converts a GOBL envelope into a FA_VAT document
func NewDocument(env *gobl.Envelope) (*Faktura, error) {
	inv, ok := env.Extract().(*bill.Invoice)
	if !ok {
		return nil, fmt.Errorf("invalid type %T", env.Document)
	}

	faktura := &Faktura{
		XSINamespace: XSINamespace,
		XSDNamespace: XSDNamespace,
		XMLNamespace: XMLNamespace,

		Naglowek: NewNaglowek(inv),
		Podmiot1: NewPodmiot1(inv.Supplier),
		Podmiot2: NewPodmiot2(inv.Customer),
		Fa:       NewFa(inv),
		Stopka:   NewStopka(inv),
	}

	return faktura, nil
}

// Bytes returns the XML representation of the document in bytes
func (d *Faktura) Bytes() ([]byte, error) {
	bytes, err := xml.MarshalIndent(d, "", "  ")
	if err != nil {
		return nil, err
	}

	return append([]byte(xml.Header), bytes...), nil
}
