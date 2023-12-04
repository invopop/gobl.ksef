package ksef

import (
	"github.com/invopop/gobl/org"
)

type Address struct {
	CountryCode string `xml:"KodKraju"`
	AddressL1   string `xml:"AdresL1"`
	AddressL2   string `xml:"AdresL2"`
}

type Seller struct {
	NIP     string   `xml:"DaneIdentyfikacyjne>NIP"`
	Name    string   `xml:"DaneIdentyfikacyjne>Nazwa"`
	Address *Address `xml:"Adres"`
}

type Buyer struct {
	NIP string `xml:"DaneIdentyfikacyjne>NIP,omitempty"`
	// or
	UECode      string `xml:"DaneIdentyfikacyjne>KodUE,omitempty"`
	UEVatNumber string `xml:"DaneIdentyfikacyjne>NrVatUE,omitempty"`
	// or
	CountryCode string `xml:"DaneIdentyfikacyjne>KodKraju,omitempty"`
	IdNumber    string `xml:"DaneIdentyfikacyjne>NrId,omitempty"`
	// or
	NoId int `xml:"DaneIdentyfikacyjne>BrakID,omitempty"`

	Name    string   `xml:"DaneIdentyfikacyjne>Nazwa,omitempty"`
	Address *Address `xml:"Adres,omitempty"`
}

func NewAddress(address *org.Address) *Address {
	adres := &Address{
		CountryCode: string(address.Country),
		AddressL1:   addressLine1(address),
		AddressL2:   addressLine2(address),
	}

	return adres
}

func NameToString(name org.Name) string {
	return name.Prefix + nameMaybe(name.Given) +
		nameMaybe(name.Middle) + nameMaybe(name.Surname) +
		nameMaybe(name.Surname2) + nameMaybe(name.Suffix)
}

func NewSeller(supplier *org.Party) *Seller {
	var name string
	if supplier.Name != "" {
		name = supplier.Name
	} else {
		name = NameToString(supplier.People[0].Name)
	}
	seller := &Seller{
		Address: NewAddress(supplier.Addresses[0]),
		Name:    name,
		NIP:     string(supplier.TaxID.Code),
	}

	return seller
}

func NewBuyer(customer *org.Party) *Buyer {

	buyer := &Buyer{
		Name: customer.Name,
		NIP:  string(customer.TaxID.Code),
		// TODO other DaneIdentyfikacyjne types
	}
	if len(customer.Addresses) > 0 {
		buyer.Address = NewAddress(customer.Addresses[0])
	}

	return buyer
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

func nameMaybe(element string) string {
	if element != "" {
		return " " + element
	}
	return ""
}
