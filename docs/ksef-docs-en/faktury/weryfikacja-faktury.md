# Invoice Verification
15.01.2026

An invoice sent to the KSeF system undergoes a series of technical and semantic checks. Verification includes the following criteria:

## XSD Schema Compliance
The invoice must be prepared in XML format, encoded in UTF-8 without the BOM marker (first 3 bytes 0xEF 0xBB 0xBF), compliant with the declared schema specified when opening the session.

## Invoice Uniqueness
- KSeF detects duplicate invoices globally, based on data stored in the system. The duplicate identification criterion is a combination of:
  1. Seller's NIP (`Podmiot1:NIP`)
  2. Invoice type (`RodzajFaktury`)
  3. Invoice number (`P_2`)
- In case of a duplicate, error code 440 ("Duplicate invoice") is returned.
- Invoice uniqueness is maintained in KSeF for a period of 10 full years, counted from the end of the calendar year in which the invoice was issued.
- The uniqueness criterion always refers to the seller (Podmiot1:NIP). In cases where different units issue invoices on behalf of the same entity (e.g., branches, organizational units of local government entities, other authorized entities), they must agree on numbering rules to avoid duplicates.

## Date Validation
The invoice issue date (`P_1`) cannot be later than the date of document acceptance into the KSeF system.

## NIP Number Validation
  - NIP checksum verification for: `Podmiot1`, `Podmiot2`, `Podmiot3`, and `PodmiotUpowazniony` (if present).
  - Applies only to production environment.

## NIP Number Validation in Internal Identifier
  - NIP checksum verification in internal identifier (`InternalId`) for `Podmiot3` - if this identifier is present.
  - Applies only to production environment.

## File Size
- Maximum invoice size without attachments: **1 MB \*** (1,000,000 bytes).
- Maximum invoice size with attachments: **3 MB \*** (3,000,000 bytes).

## Quantity Limits
- The maximum number of invoices in a single session (both interactive and batch) is 10,000 *.
- In a batch upload, a maximum of 50 ZIP files can be sent; the size of each file before encryption cannot exceed 100 MB (100,000,000 bytes), and the total size of the ZIP package - 5 GB (5,000,000,000 bytes).

## Proper Encryption
- The invoice should be encrypted using the AES-256-CBC algorithm (256-bit symmetric key, 128-bit IV, with PKCS#7 padding).
- The symmetric key is encrypted using the RSAES-OAEP algorithm (SHA-256/MGF1).

## Invoice Metadata Compliance in Interactive Session
- Calculation and verification of the invoice hash along with file size.
- Calculation and verification of the encrypted invoice hash along with file size.

## Attachment Restrictions
- Sending invoices with attachments is only allowed in batch mode.
**Exception:** When sending a [technical correction of an offline invoice](../offline/korekta-techniczna.md), the use of an interactive session is permitted.
- The ability to send invoices with attachments requires prior registration of this option in the `e-Tax Office` service.

## Authorization Requirements
Sending an invoice to KSeF requires having the appropriate authorization to issue it in the context of the given entity.

\* **Note:** If the available [limits](../limity/limity.md) are insufficient for your organization's business scenarios, please contact KSeF support for an individual analysis and selection of an appropriate solution.
