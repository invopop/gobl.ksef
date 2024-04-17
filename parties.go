package ksef

import (
	"github.com/invopop/gobl/l10n"
	"github.com/invopop/gobl/org"
)

// Address defines the XML structure for KSeF addresses
type Address struct {
	CountryCode string `xml:"KodKraju"`
	AddressL1   string `xml:"AdresL1"`
	AddressL2   string `xml:"AdresL2,omitempty"`
}

// Seller defines the XML structure for KSeF seller
type Seller struct {
	NIP     string          `xml:"DaneIdentyfikacyjne>NIP"`
	Name    string          `xml:"DaneIdentyfikacyjne>Nazwa"`
	Address *Address        `xml:"Adres"`
	Contact *ContactDetails `xml:"DaneKontaktowe,omitempty"`
}

// ContactDetails defines the XML structure for KSeF contact
type ContactDetails struct {
	Phone string `xml:"Telefon,omitempty"`
	Email string `xml:"Email,omitempty"`
}

// Buyer defines the XML structure for KSeF buyer
type Buyer struct {
	NIP string `xml:"DaneIdentyfikacyjne>NIP,omitempty"`
	// or
	UECode      string `xml:"DaneIdentyfikacyjne>KodUE,omitempty"`
	UEVatNumber string `xml:"DaneIdentyfikacyjne>NrVatUE,omitempty"`
	// or
	CountryCode string `xml:"DaneIdentyfikacyjne>KodKraju,omitempty"`
	IDNumber    string `xml:"DaneIdentyfikacyjne>NrID,omitempty"`
	// or
	NoID int `xml:"DaneIdentyfikacyjne>BrakID,omitempty"`

	Name    string          `xml:"DaneIdentyfikacyjne>Nazwa,omitempty"`
	Address *Address        `xml:"Adres,omitempty"`
	Contact *ContactDetails `xml:"DaneKontaktowe,omitempty"`
}

// newAddress gets the address data from GOBL address
func newAddress(address *org.Address) *Address {
	adres := &Address{
		CountryCode: string(address.Country),
		AddressL1:   addressLine1(address),
		AddressL2:   addressLine2(address),
	}

	return adres
}

// nameToString get the seller name out of the organization
func nameToString(name org.Name) string {
	return name.Prefix + nameMaybe(name.Given) +
		nameMaybe(name.Middle) + nameMaybe(name.Surname) +
		nameMaybe(name.Surname2) + nameMaybe(name.Suffix)
}

// NewSeller converts a GOBL Party into a KSeF seller
func NewSeller(supplier *org.Party) *Seller {
	var name string
	if supplier.Name != "" {
		name = supplier.Name
	} else {
		name = nameToString(supplier.People[0].Name)
	}
	seller := &Seller{
		Address: newAddress(supplier.Addresses[0]),
		NIP:     string(supplier.TaxID.Code),
		Name:    name,
	}
	if len(supplier.Telephones) > 0 {
		seller.Contact = &ContactDetails{
			Phone: supplier.Telephones[0].Number,
		}
	}
	if len(supplier.Emails) > 0 {
		if seller.Contact == nil {
			seller.Contact = &ContactDetails{}
		}
		seller.Contact.Email = supplier.Emails[0].Address
	}

	return seller
}

// NewBuyer converts a GOBL Party into a KSeF buyer
func NewBuyer(customer *org.Party) *Buyer {

	buyer := &Buyer{
		Name: customer.Name,
	}

	if len(customer.TaxID.Code) == 0 {
		buyer.NoID = 1
	} else if customer.TaxID.Country == l10n.PL {
		buyer.NIP = string(customer.TaxID.Code)
	} else if isEUCountry(customer.TaxID.Country) {
		buyer.UEVatNumber = string(customer.TaxID.Code)
		buyer.UECode = string(customer.TaxID.Country)
	} else {
		buyer.IDNumber = string(customer.TaxID.Code)
		buyer.CountryCode = string(customer.TaxID.Country)
	}

	if len(customer.Addresses) > 0 {
		buyer.Address = newAddress(customer.Addresses[0])
	}

	if len(customer.Telephones) > 0 {
		buyer.Contact = &ContactDetails{
			Phone: customer.Telephones[0].Number,
		}
	}
	if len(customer.Emails) > 0 {
		if buyer.Contact == nil {
			buyer.Contact = &ContactDetails{}
		}
		buyer.Contact.Email = customer.Emails[0].Address
	}

	return buyer
}

func EUCountries() []l10n.CountryCode {
	return []l10n.CountryCode{
		l10n.AT, // Austria,
		l10n.BE, // Belgium,
		l10n.BG, // Bulgaria,
		l10n.HR, // Croatia,
		l10n.CY, // Republic of Cyprus,
		l10n.CZ, // Czech Republic,
		l10n.DK, // Denmark,
		l10n.EE, // Estonia,
		l10n.FI, // Finland,
		l10n.FR, // France,
		l10n.DE, // Germany,
		l10n.GR, // Greece,
		l10n.HU, // Hungary,
		l10n.IE, // Ireland,
		l10n.IT, // Italy,
		l10n.LV, // Latvia,
		l10n.LT, // Lithuania,
		l10n.LU, // Luxembourg,
		l10n.MT, // Malta,
		l10n.NL, // Netherlands,
		l10n.PT, // Portugal,
		l10n.RO, // Romania,
		l10n.SK, // Slovakia,
		l10n.SL, // Slovenia,
		l10n.ES, // Spain,
		l10n.SE, // Sweden
	}
}

func isEUCountry(country l10n.CountryCode) bool {
	return contains(EUCountries(), country)
}

func contains(slice []l10n.CountryCode, item l10n.CountryCode) bool {
	for _, value := range slice {
		if value == item {
			return true
		}
	}
	return false
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
