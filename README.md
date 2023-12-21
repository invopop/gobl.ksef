# GOBL to KSeF Conversion

Convert GOBL to the Polish FA_VAT format and send to KSeF.

Copyright [Invopop Ltd.](https://invopop.com) 2023. Released publicly under the [Apache License Version 2.0](LICENSE). For commercial licenses please contact the [dev team at invopop](mailto:dev@invopop.com). In order to accept contributions to this library we will require transferring copyrights to Invopop Ltd.

## Project Development Objectives

The following list the steps to follow through on in order to accomplish the goal of using GOBL to submit electronic invoices to the Polish authorities:

1. Add the PL (`pl`) tax regime to [GOBL](https://github.com/invopop/gobl). Figure out local taxes, tax ID validation rules, and any "extensions" that may be required to be defined in GOBL, and send in a PR. For examples of existing regimes, see the [regimes](https://github.com/invopop/gobl/tree/main/regimes) directory. Key Concerns:
   - Basic B2B invoices support.
   - Tax ID validation as per local rules.
   - Support for "simplified" invoices.
   - Requirements for credit-notes or "rectified" invoices and the correction options definition for the tax regime.
   - Any additional fields that need to be validated, like payment terms.
2. Convert GOBL into FA_VAT format in library. A couple of good examples: [gobl.cfdi for mexico](https://github.com/invopop/gobl.cfdi) and [gobl.facturae for Spain](https://github.com/invopop/gobl.facturae). Library would just be able to run tests in the first version.
3. Build a CLI (copy from gobl.cfdi and gobl.facture projects) to convert GOBL JSON documents into FA_VAT XML.
4. Build a second part of this project that allows documents to be sent directly to the KSeF. A partial example of this can be found in the [gobl.ticketbai project](https://github.com/invopop/gobl.ticketbai/tree/refactor/internal/gateways). It'd probably be useful to be able to upload via the CLI too.

## KSeF API

Useful links:

- [National e-Invoice System](https://www.podatki.gov.pl/ksef/) - for details on system in general.
- [KSeF Test Zone](https://www.podatki.gov.pl/ksef/strefa-testowa-ksef/)

KSeF provide three environments:

1.  [Test Environment](https://ksef-test.mf.gov.pl/) for application development with fictitious data.
2.  [Pre-production "demo"](https://ksef-demo.mf.gov.pl/) area with production data, but not officially declared.
3.  [Production](https://ksef.mf.gov.pl)

A translation of the Interface Specification 1.5 is available in the [docs](./docs) folder.

OpenAPI documentation is available for three specific interfaces:

1. Batches ([test openapi 'batch' spec](https://ksef-test.mf.gov.pl/openapi/gtw/svc/api/KSeF-batch.yaml)) - for sending multiple documents at the same time.
2. Common ([test openapi 'common' spec](https://ksef-test.mf.gov.pl/openapi/gtw/svc/api/KSeF-common.yaml)) - general operations that don't require authentication.
3. Interactive ([test openapi 'online' spec](https://ksef-test.mf.gov.pl/openapi/gtw/svc/api/KSeF-online.yaml))

## Authentication

Authentication with KSeF appears to be done using digital certificates issued by trusted service providers approved by [NCCert Poland](https://www.nccert.pl/).

There is an online process to register a company:

- [Test Company login](https://ksef-test.mf.gov.pl/web/login)
- [Generate a fake NIP (tax ID)](http://generatory.it/)

Once inside the test environment, you can create an Authorization token to use to make requests to the API.
