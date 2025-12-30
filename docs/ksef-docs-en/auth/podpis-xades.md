## XAdES Signature
https://www.w3.org/TR/XAdES/

Allowed signature formats:
- enveloped
- enveloping

External (detached) signature format is not accepted.

Allowed transforms in XAdES signature:
- http://www.w3.org/TR/1999/REC-xpath-19991116 - not(ancestor-or-self::ds:Signature)
- http://www.w3.org/2002/06/xmldsig-filter2
- http://www.w3.org/2000/09/xmldsig#enveloped-signature
- http://www.w3.org/2000/09/xmldsig#base64
- http://www.w3.org/2006/12/xml-c14n11
- http://www.w3.org/2006/12/xml-c14n11#WithComments
- http://www.w3.org/2001/10/xml-exc-c14n#
- http://www.w3.org/2001/10/xml-exc-c14n#WithComments
- http://www.w3.org/TR/2001/REC-xml-c14n-20010315
- http://www.w3.org/TR/2001/REC-xml-c14n-20010315#WithComments

### Allowed Certificate Types

Allowed certificate types in XAdES signature:
* Qualified certificate of a natural person - containing the PESEL or NIP number of the person authorized to act on behalf of the company,
* Qualified certificate of an organization (so-called company seal) - containing the NIP number,
* Trusted Profile (ePUAP) - allows document signing; used by natural persons,
* KSeF internal certificate - issued by the KSeF system. This certificate is not a qualified certificate, but is honored in the authentication process.

**Qualified certificate** - a certificate issued by a qualified trust service provider, registered in the EU registry [EU Trusted List (EUTL)](https://eidas.ec.europa.eu/efda/trust-services/browse/eidas/tls), in accordance with the eIDAS regulation. KSeF accepts qualified certificates issued in Poland and in other European Union member states.

### Required Attributes of Qualified Certificates

#### Qualified Signature Certificates (issued to natural persons)

Required subject attributes:<br/>
| Identifier (OID) | Name           | Meaning                                  |
|------------------|----------------|------------------------------------------|
| 2.5.4.42         | givenName      | first name                               |
| 2.5.4.4          | surname        | last name                                |
| 2.5.4.5          | serialNumber   | serial number                            |
| 2.5.4.3          | commonName     | common name of the certificate owner     |
| 2.5.4.6          | countryName    | country name, ISO 3166 code              |

Recognized patterns for `serialNumber` attribute:<br>
**(PNOPL|PESEL).\*?(?<number>\\d{11})**<br>
**(TINPL|NIP).\*?(?<number>\\d{10})**<br>

#### Qualified Seal Certificates (issued to organizations)

Required subject attributes:<br/>
| Identifier (OID) | Name                    | Meaning                                                              |
|------------------|-------------------------|----------------------------------------------------------------------|
| 2.5.4.10         | organizationName        | full formal name of the entity for which the certificate is issued   |
| 2.5.4.97         | organizationIdentifier  | entity identifier                                                    |
| 2.5.4.3          | commonName              | common name of the organization                                      |
| 2.5.4.6          | countryName             | country name, ISO 3166 code                                          |

Prohibited subject attributes:
| Identifier (OID) | Name        | Meaning    |
|------------------|------------ |------------|
| 2.5.4.42         | givenName   | first name |
| 2.5.4.4          | surname     | last name  |

Recognized patterns for `organizationIdentifier` attribute:<br>
**(VATPL).\*?(?<number>\\d{10})**<br>

### Certificate Fingerprint

For qualified certificates that do not have proper identifiers stored in the subject attribute OID.2.5.4.5, authentication with such a certificate is possible after prior authorization on the SHA-256 hash (so-called fingerprint) of that certificate.
