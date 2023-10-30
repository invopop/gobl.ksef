# GOBL to KSeF Conversion

Convert GOBL to the Polish FA_VAT format and send to KSeF.

Copyright [Invopop Ltd.](https://invopop.com) 2023. Released publicly under the [Apache License Version 2.0](LICENSE). For commercial licenses please contact the [dev team at invopop](mailto:dev@invopop.com). In order to accept contributions to this library we will require transferring copyrights to Invopop Ltd.

## Project Development Objectives

The following list the steps to follow through on in order to accomplish the goal of using GOBL to submit electronic invoices to the Polish authorities:

1. Figure out local taxes, Tax ID validation rules, and any "extensions" that may be required to be defined in GOBL, and send in a PR with the changes to the [GOBL Repository](https://github.com/invopop/gobl). For examples, see the [regimes](https://github.com/invopop/gobl/tree/main/regimes) directory. Key Concerns:
  - Basic B2B invoices support.
  - Tax ID validation as per local rules.
  - Support for "simplified" invoices.
  - Requirements for credit-notes or "rectified" invoices and the correction options definition for the tax regime.
  - Any additional fields that need to be validated, like payment terms.
2. Convert GOBL into FA_VAT format in library. A couple of good examples: [gobl.cfdi for mexico](https://github.com/invopop/gobl.cfdi) and [gobl.facturae for Spain](https://github.com/invopop/gobl.facturae). Library would just be able to run tests in the first version.
3. Build a CLI (copy from gobl.cfdi and gobl.facture projects) to convert GOBL JSON documents into FA_VAT XML.
4. Build a second part of this project that allows documents to be sent directly to the KSeF. A partial example of this can be found in the [gobl.ticketbai project](https://github.com/invopop/gobl.ticketbai/tree/refactor/internal/gateways). It'd probably be useful to be able to upload via the CLI too.