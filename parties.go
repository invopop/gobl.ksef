package ksef

import (
	"github.com/invopop/gobl/org"
)

type Adres struct {
	KodKraju string `xml:"KodKraju"`
	AdresL1  string `xml:"AdresL1"`
	AdresL2  string `xml:"AdresL2"`
}

type Podmiot1 struct {
	Adres *Adres `xml:"Adres"`
	NIP   string `xml:"DaneIdentyfikacyjne>NIP"`
	Nazwa string `xml:"DaneIdentyfikacyjne>Nazwa"`
}

type Podmiot2 struct {
	Adres *Adres `xml:"Adres,omitempty"`

	NIP string `xml:"DaneIdentyfikacyjne>NIP,omitempty"`
	// or
	KodUE   string `xml:"DaneIdentyfikacyjne>KodUE,omitempty"`
	NrVatUE string `xml:"DaneIdentyfikacyjne>NrVatUE,omitempty"`
	// or
	KodKraju string `xml:"DaneIdentyfikacyjne>KodKraju,omitempty"`
	NrId     string `xml:"DaneIdentyfikacyjne>NrId,omitempty"`
	// or
	BrakID int `xml:"DaneIdentyfikacyjne>BrakID,omitempty"`

	Nazwa string `xml:"DaneIdentyfikacyjne>Nazwa,omitempty"`
}

func NewAdres(address *org.Address) *Adres {
	adres := &Adres{
		KodKraju: string(address.Country),
		AdresL1:  addressLine1(address),
		AdresL2:  addressLine2(address),
	}

	return adres
}

func NewPodmiot1(supplier *org.Party) *Podmiot1 {

	podmiot1 := &Podmiot1{
		Adres: NewAdres(supplier.Addresses[0]),
		Nazwa: supplier.Name,
		NIP:   string(supplier.TaxID.Code),
	}

	return podmiot1
}

func NewPodmiot2(customer *org.Party) *Podmiot2 {

	podmiot2 := &Podmiot2{
		Nazwa: customer.Name,
		NIP:   string(customer.TaxID.Code),
		// TODO other DaneIdentyfikacyjne types
	}
	if len(customer.Addresses) > 0 {
		podmiot2.Adres = NewAdres(customer.Addresses[0])
	}

	return podmiot2
}

func addressLine1(address *org.Address) string {
	if address.PostOfficeBox != "" {
		return address.PostOfficeBox
	}

	return address.Street +
		", " + address.Number +
		addressMaybe(address.Block) +
		addressMaybe(address.Floor) +
		addressMaybe(address.Door)
}

func addressLine2(address *org.Address) string {
	return address.Code + ", " + address.Locality
}

func addressMaybe(element string) string {
	if element != "" {
		return ", " + element
	}
	return ""
}
