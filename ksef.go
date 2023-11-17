// Package ksef implements the conversion from GOBL to FA_VAT XML
package ksef

import (
	"encoding/xml"
	"fmt"

	"github.com/invopop/gobl"
	"github.com/invopop/gobl/bill"
)

// Document is a pseudo-model for containing the XML document being created
type Document struct {
	Naglowek *Naglowek `xml:"Naglowek"`
	Podmiot1 *Podmiot1 `xml:"Podmiot1"`
	Podmiot2 *Podmiot2 `xml:"Podmiot2"`
	Fa       *Fa       `xml:"Fa"`
	Stopka   *Stopka   `xml:"Stopka,omitempty"`
}

// NewDocument converts a GOBL envelope into a FA_VAT document
func NewDocument(env *gobl.Envelope) (*Document, error) {
	inv, ok := env.Extract().(*bill.Invoice)
	if !ok {
		return nil, fmt.Errorf("invalid type %T", env.Document)
	}

	document := &Document{

		Naglowek: NewNaglowek(inv),
	}

	return document, nil
}

// Bytes returns the XML representation of the document in bytes
func (d *Document) Bytes() ([]byte, error) {
	bytes, err := xml.MarshalIndent(d, "", "  ")
	if err != nil {
		return nil, err
	}

	return append([]byte(xml.Header), bytes...), nil
}
