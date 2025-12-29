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
2. Convert GOBL into FA_VAT format in library. A couple of good examples: [gobl.cfdi for Mexico](https://github.com/invopop/gobl.cfdi) and [gobl.facturae for Spain](https://github.com/invopop/gobl.facturae). Library would just be able to run tests in the first version.
3. Build a CLI (copy from gobl.cfdi and gobl.facture projects) to convert GOBL JSON documents into FA_VAT XML.
4. Build a second part of this project that allows documents to be sent directly to the KSeF. A partial example of this can be found in the [gobl.ticketbai project](https://github.com/invopop/gobl.ticketbai/tree/refactor/internal/gateways). It'd probably be useful to be able to upload via the CLI too.

## FA_VAT documentation

FA_VAT is the Polish electronic invoice format. The format uses XML.

- [XML schema](https://github.com/CIRFMF/ksef-docs/blob/main/faktury/schemy/FA/schemat_FA(3)_v1-0E.xsd) for V3 (description of fields is in Polish)
- [Types definition](https://raw.githubusercontent.com/CIRFMF/ksef-docs/refs/heads/main/faktury/schemy/FA/bazowe/ElementarneTypyDanych_v10-0E.xsd) (description of fields is in Polish) - we have to open it as raw, as [the original link](https://github.com/CIRFMF/ksef-docs/blob/main/faktury/schemy/FA/bazowe/StrukturyDanych_v10-0E.xsd) does not add newlines
- [Complex types definition](https://raw.githubusercontent.com/CIRFMF/ksef-docs/refs/heads/main/faktury/schemy/FA/bazowe/StrukturyDanych_v10-0E.xsd) (description of fields is in Polish) - we have to open it as raw, as [the original link](https://github.com/CIRFMF/ksef-docs/blob/main/faktury/schemy/FA/bazowe/StrukturyDanych_v10-0E.xsd) does not add newlines

## KSeF API

KSeF is the Polish system for submitting electronic invoices to the Polish authorities.

Useful links:

- [National e-Invoice System](https://www.podatki.gov.pl/ksef/) - for details on system in general.
- [KSeF Test Zone](https://www.podatki.gov.pl/ksef/strefa-testowa-ksef/)
- [API documentation](https://ksef-test.mf.gov.pl/docs/v2/index.html) for the test environment (in Polish)

KSeF provide three environments:

1.  [Test Environment](https://ksef-test.mf.gov.pl/) for application development with fictitious data.
2.  [Pre-production "demo"](https://ksef-demo.mf.gov.pl/) area with production data, but not officially declared.
3.  [Production](https://ksef.mf.gov.pl)

A translation of the Interface Specification 1.5 is available in the [docs](./docs) folder.

OpenAPI documentation is available for three specific interfaces:

1. Batches ([test openapi 'batch' spec](https://ksef-test.mf.gov.pl/openapi/gtw/svc/api/KSeF-batch.yaml)) - for sending multiple documents at the same time.
2. Common ([test openapi 'common' spec](https://ksef-test.mf.gov.pl/openapi/gtw/svc/api/KSeF-common.yaml)) - general operations that don't require authentication.
3. Interactive ([test openapi 'online' spec](https://ksef-test.mf.gov.pl/openapi/gtw/svc/api/KSeF-online.yaml)) - sending a single document in each request.

## Authentication

- [Authentication in KSEF](https://github.com/CIRFMF/ksef-docs/blob/main/uwierzytelnianie.md) (in Polish)
- [XAdEs digital signature](https://github.com/CIRFMF/ksef-docs/blob/main/auth/podpis-xades.md) (in Polish)
- [How to use the official .NET client to generate a test XAdEs certificate](https://github.com/CIRFMF/ksef-docs/blob/main/auth/testowe-certyfikaty-i-podpisy-xades.md) (in Polish)

### How to obtain a certificate

Authentication with KSeF appears to be done using digital certificates issued by trusted service providers approved by [NCCert Poland](https://www.nccert.pl/).

There is an online process to register a company:

- [Test Company login](https://ksef-test.mf.gov.pl/web/login)
- [Generate a fake NIP (tax ID)](http://generatory.it/)

Once inside the test environment, you can create an Authorization token to use to make requests to the API.

### Authentication to the KSeF API

This is translated from [Authentication in KSeF](https://github.com/CIRFMF/ksef-docs/blob/main/uwierzytelnianie.md) document:

To login, XAdES digital signature or a KSeF token is needed.

Base URL for the test environment: https://api-test.ksef.mf.gov.pl/v2

1. Submit `POST /auth/challenge` with no body, no headers, response has fields `challenge` (opaque string), `timestamp`, `timestampMs`.
2. Depending on the login method (XAdES / KSeF token) submit `POST /auth/xades-signature` (body is in XML and should contain challenge, context, signature) or `POST /auth/ksef-token` (body contains challenge, context, KSeF token + timestamp encrypted with public key). In both cases we receive JSON response with `referenceNumber` and `authenticationToken`.
3. Keep polling `GET /auth/[referenceNumber]` with header `Bearer [authenticationToken]` - field `status` will indicate: 100 - in progress, 200 - successful, 4xx/5xx - error.
4. When the endpoint above returns status 200, send `POST /auth/token/redeem` with header `Bearer [authenticationToken]` and no body. Response contains `accessToken` and `refreshToken` + their expiration times. Redeeming can be done once - more attempts will result in 40x errors.
`accessToken` can be used for most actions in the API.
5. To get a new `accessToken`, send `POST /auth/token/refresh` with header `Bearer [refreshToken]` and no body. Response contains a new `accessToken` + expiration time.
6. List of current login sessions is at `GET /auth/sessions`.
7. To logout, send `DELETE /auth/sessions/current` or `DELETE /auth/sessions/[referenceNumber]`.

A single subject can have multiple login sessions. One login session is associated with a single context.

Note that the API documentation uses the name `referenceNumber` in other endpoints for asynchronous operations (submit and poll for status), not only for identifying login sessions.

### What is subject and context?

- Subject = who is logging in
- Context = business entity the operations are about

E.g. context = company X, subject = accountant of company X, employee of an accounting company having contract with company X, etc. This way, a single accountant or accounting company can work with multiple companies.

Using an API endpoint, a subject having appropriate permissions to a given context can provide another subject with permissions to the same context. It's possible to revoke permissions for a subject. It's also possible to mark the permissions as transferable - this is useful when company X gives permissions to accounting company Y, and company Y gives permissions to one of its employees.

### What is the KSeF token?

KSEF token (separate type from authentication token, access token and refresh token - don’t confuse them!) is intended for API integration, and has specific permissions and description string (name entered by user). KSEF tokens can be revoked. Permissions on a KSEF token cannot be changed - to add new permissions, it's necessary to create a new token.

Users logged in with XAdES can create, list and delete KSEF tokens using the API.

### How to obtain the public key

To obtain the public key certificate, use `GET /security/public-key-certificates`.

[Documentation here](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Certyfikaty-klucza-publicznego)

Public key is needed to:
1. Login with KSeF token using the `POST /auth/ksef-token` endpoint.
2. Encrypt a symmetric AES key when uploading invoices (in both online and batch formats) and exporting (batch) incoming invoices. File upload and export endpoints don't require `accessToken`, but they accept or return chunks of the data respectively encrypted with the provided symmetric AES key. Online upload endpoint is a regular HTTP endpoint using `accessToken` for authentication, but also requires providing a key, and the invoice must be encrypted with that key.
