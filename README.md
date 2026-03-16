# GOBL ↔ KSeF Conversion

Bidirectional conversion between GOBL and the Polish FA_VAT XML format (KSeF).

## Main Conversion Entrypoints

**GOBL → KSeF:**
- `ksef.BuildFavat(env *gobl.Envelope) (*Invoice, error)` - Converts a GOBL envelope to a KSeF FA_VAT invoice model
- `(*Invoice).Bytes() ([]byte, error)` - Returns the XML representation as bytes

**KSeF → GOBL:**
- `ksef.ParseKSeF(xmlData []byte) (*gobl.Envelope, error)` - Converts KSeF FA_VAT XML to a GOBL envelope

Copyright [Invopop Ltd.](https://invopop.com) 2023. Released publicly under the [Apache License Version 2.0](LICENSE). For commercial licenses please contact the [dev team at invopop](mailto:dev@invopop.com). In order to accept contributions to this library we will require transferring copyrights to Invopop Ltd.

## Supported Features

The converter handles the following invoice types and features:

**Invoice types (GOBL → KSeF):**
- `VAT` - Standard invoices
- `ZAL` - Advance/prepayment invoices (tag: `partial`)
- `ROZ` - Settlement invoices (tag: `settlement`)
- `UPR` - Simplified invoices (tag: `simplified`)
- `KOR` - Correction invoices (credit notes)
- `KOR_ZAL` - Correction of advance invoices
- `KOR_ROZ` - Correction of settlement invoices

**Parties:**
- Seller (Podmiot1) with Polish NIP, address, contact details, and EU VAT prefix
- Buyer (Podmiot2) with Polish NIP, EU VAT number, or non-EU tax ID
- Third parties (Podmiot3) for JST (local government units) and Group VAT scenarios

**Tax rates:**
- Standard (23%), reduced (8%), super-reduced (5%), and other percentage rates
- Zero-rated (0 KR), intra-community (0 WDT), export (0 EX)
- Tax exempt (zw), outside scope (np I), reverse charge (np II), domestic reverse charge (oo)
- OSS (One Stop Shop) rates
- Margin scheme rates

**Annotations:**
- Cash accounting, self-billing, reverse charge, split payment mechanism
- Tax exemption with legal basis (Polish law, EU directive, or other)
- Margin scheme (travel agency, used goods, art works, collectibles/antiques)

**Other features:**
- Line item discounts
- Invoice periods (P_6_Od / P_6_Do)
- Correction/credit note references with KSeF numbers
- Payment details: means of payment, bank accounts, due dates, advance payments
- Additional description lines (DodatkowyOpis)
- Ordering data with order references and order lines (Zamowienie / WarunkiTransakcji)
- Rounding adjustments to reconcile KSeF and GOBL totals
- Gross pricing (P_9B) support in KSeF → GOBL direction
- Credit notes with before/after correction lines (StanPrzed)
- Prepayment invoices without line items (bypass mode with totals from tax summaries)

## CLI

The `gobl.ksef` CLI provides two commands:

- **`convert`** - Convert a GOBL JSON envelope into a FA_VAT XML document:
  ```bash
  gobl.ksef convert input.json output.xml
  ```

- **`send`** - Convert and send a GOBL JSON envelope to the KSeF API:
  ```bash
  gobl.ksef send input.json [nip] [token] [keyPath]
  ```

## Testing

The test suite includes tests for both conversion directions and round-trip validation.

### Running Tests

**Run all tests:**
```bash
go test ./test -v
```

**Update golden files:**
```bash
go test ./test --update -v
```

**With XSD schema validation (requires libxml2):**
```bash
# Using the helper script (sets LD_LIBRARY_PATH automatically)
./test/test.sh -v
./test/test.sh --update -v

# Or manually
LD_LIBRARY_PATH=/home/linuxbrew/.linuxbrew/opt/libxml2/lib:$LD_LIBRARY_PATH go test -tags xsdvalidate ./test -v
```

### Test Data

**GOBL → KSeF conversion:**
- **Input**: GOBL JSON files in `test/data/gobl.ksef/*.json`
- **Output**: KSeF XML files in `test/data/gobl.ksef/out/*.xml`

**KSeF → GOBL conversion:**
- **Input**: KSeF XML files in `test/data/ksef.gobl/*.xml`
- **Output**: GOBL JSON files in `test/data/ksef.gobl/out/*.json`

**Schema validation:**
- **Schema**: FA3 XSD and dependencies in `test/data/schema/`

## Unsupported fields

See [unsupported-fields.md](unsupported-fields.md) for the list of unsupported fields.

## FA_VAT documentation

FA_VAT is the Polish electronic invoice format. The format uses XML.

- [XML schema](https://github.com/CIRFMF/ksef-docs/blob/main/faktury/schemy/FA/schemat_FA(3)_v1-0E.xsd) for V3 (description of fields is in Polish)
- [Types definition](https://raw.githubusercontent.com/CIRFMF/ksef-docs/refs/heads/main/faktury/schemy/FA/bazowe/ElementarneTypyDanych_v10-0E.xsd) (description of fields is in Polish) - we have to open it as raw, as [the original link](https://github.com/CIRFMF/ksef-docs/blob/main/faktury/schemy/FA/bazowe/StrukturyDanych_v10-0E.xsd) does not add newlines
- [Complex types definition](https://raw.githubusercontent.com/CIRFMF/ksef-docs/refs/heads/main/faktury/schemy/FA/bazowe/StrukturyDanych_v10-0E.xsd) (description of fields is in Polish) - we have to open it as raw, as [the original link](https://github.com/CIRFMF/ksef-docs/blob/main/faktury/schemy/FA/bazowe/StrukturyDanych_v10-0E.xsd) does not add newlines

## Parsing (KSeF → GOBL)

The parsing functionality converts KSeF FA_VAT XML documents back into GOBL format. The implementation includes:

- **Party conversion**: Converts seller (Podmiot1), buyer (Podmiot2), and third parties (Podmiot3) to GOBL parties
- **Invoice data**: Parses invoice metadata including type, codes, dates, currency, and annotations
- **Line items**: Converts FA_VAT line items (FaWiersz) to GOBL invoice lines, including discounts and tax combos
- **Credit notes**: Handles line inversion for correction invoices, including before/after correction lines (StanPrzed)
- **Gross pricing**: Supports invoices using gross unit prices (P_9B) by setting `PricesInclude = VAT`
- **Ordering data**: Maps Zamowienie (order) and WarunkiTransakcji (transaction conditions) to GOBL ordering purchases
- **Payment details**: Extracts payment means, bank accounts, due dates, and advance payments
- **Prepayment invoices**: Handles advance invoices without line items (ZAL/KOR_ZAL) using bypass mode with totals from tax summary fields
- **Settlement invoices**: Derives advance payments for ROZ/KOR_ROZ invoices (see below)
- **Rounding adjustments**: Handles rounding differences between KSeF and GOBL calculation methods
- **Round-trip validation**: All GOBL → KSeF conversions are validated through round-trip tests (GOBL → KSeF → GOBL)

## Settlement Invoices (ROZ)

Settlement invoices (`ROZ`) finalize orders that had advance payments (`ZAL` invoices). Per Art. 106f sec. 3 of the Polish VAT Act, they must show the full order value in line items (`FaWiersz`) while `P_15` contains only the remaining amount after advance deductions. The advance invoice references appear in `FakturaZaliczkowa` elements.

This creates a structural mismatch: GOBL calculates Payable from line items (full amount), but P_15 is the remaining balance. For example, an order worth 68,363.40 PLN with a 13,672.68 PLN advance would have lines totalling 68,363.40 but P_15 = 54,690.72.

When converting KSeF → GOBL, the converter detects settlement invoices with `FakturaZaliczkowa` references and no explicit `ZaplataCzesciowa` (partial payment) entries, and derives the advance amount as `Payable - P_15`. This produces a natural GOBL invoice where lines represent the full order, advances represent prepaid amounts, and Due equals the remaining balance. When explicit `ZaplataCzesciowa` entries are present (as in some ERP systems), those are used directly and no derivation occurs.

## KSeF API

KSeF is the Polish system for submitting electronic invoices to the Polish authorities. The system uses API version 2.0, which introduced JWT-based authentication, mandatory invoice encryption, and a unified session model.

Useful links:

- [National e-Invoice System](https://ksef.mf.gov.pl/) - for details on system in general (English translation available - language picker is in the top right corner)
- [KSeF 2.0 Integrator's Guide](./docs/ksef-docs-en/README.md) - translated documentation with detailed integration instructions

KSeF provides three environments:

| Environment | Description | API Documentation |
| ----------- | ----------- | ----------------- |
| **Test** (Release Candidate) | For testing integration, contains RC versions | [api-test.ksef.mf.gov.pl](https://api-test.ksef.mf.gov.pl/docs/v2) |
| **Demo** (Pre-production) | Matches production configuration, for final validation | [api-demo.ksef.mf.gov.pl](https://api-demo.ksef.mf.gov.pl/docs/v2) |
| **Production** | Full legal validity, production data | [api.ksef.mf.gov.pl](https://api.ksef.mf.gov.pl/docs/v2) |

The OpenAPI specification is available at each environment's `/docs/v2/openapi.json` endpoint.

## Authentication

See [authentication.md](./authentication.md).

