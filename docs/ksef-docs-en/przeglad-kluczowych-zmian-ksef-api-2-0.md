# KSeF API 2.0 – Overview of Key Changes
09.06.2025

# Introduction

This document is intended for technical teams and developers with experience in KSeF API version 1.0. It provides an overview of the most important changes introduced in version 2.0, along with a discussion of new capabilities and practical improvements in integration.

The purpose of this document is to:

- Highlight the main differences compared to version 1.0
- Present the benefits of migrating to version 2.0
- Facilitate the preparation of integrations for the system requirements being implemented on February 1, 2026

---

## Documentation and Tools Supporting KSeF API 2.0 Integration

To facilitate the transition to the new API version and ensure proper implementation, a set of official materials and tools has been made available to support integrators:

**Technical Documentation (OpenAPI)**

KSeF API version 2.0 has been described in the OpenAPI standard, which enables both easy documentation browsing by developers and automatic generation of integration code.

* **Documentation** (interactive online version):
  A documentation interface in the form of a technical portal, containing descriptions of methods, data structures, parameters, and examples. Intended for use by developers and integration analysts.
  <br> Test environment (TE): \[[link](https://ksef-test.mf.gov.pl/docs/v2/index.html)\]
  <br> Production environment (PRD)

* **Specification** (OpenAPI JSON file):
  A raw OpenAPI specification file in JSON format, intended for use in integration automation tools (e.g., code generators, API contract validators).
  <br> Test environment (TE): \[[link](https://ksef-test.mf.gov.pl/docs/v2/openapi.json)\]
  <br> Production environment (PRD)

**Official KSeF 2.0 Client Integration Library (open source)**

A public library released under open source principles, developed in parallel with subsequent API versions and maintained in full compliance with the specification. It is the recommended integration tool, enabling tracking of changes and reducing the risk of incompatibility with current system releases.

* **C\#:** \[[link](https://github.com/CIRFMF/ksef-client-csharp)\]

* **Java:** \[[link](https://github.com/CIRFMF/ksef-client-java)\]

**Published Packages**

The KSeF 2.0 Client library will be available in official package repositories for the most popular programming languages. For the .NET platform, it will be published as a NuGet package, while for the Java environment – as a Maven Central artifact. Publication in these repositories will enable easy inclusion of the library in projects and automatic tracking of updates compatible with subsequent API versions.

**Step-by-Step Guide**

* **Integration guide / tutorial:**
  Practical step-by-step instructions with code snippets illustrating how to use the key system endpoints.
  <br/>\[[link](https://github.com/CIRFMF/ksef-docs)\]

# Key Changes in API 2.0

## New JWT-Based Authentication Model ##

In version 1.0, authentication was tightly coupled with opening an interactive session, which introduced many limitations and complicated integration.

In version 2.0:

* Authentication has been **separated as an independent process**, independent of session initialization.

* **Standard JWT tokens** have been introduced, which are used for authorization of all protected operations.

Benefits:

* compliance with market practices,
* ability to reuse the token to create multiple sessions,
* **support for token refresh and revocation**.

Authentication process details: \[[link](https://github.com/CIRFMF/ksef-docs/blob/main/uwierzytelnianie.md)\]

## Unified Initialization Process for Batch and Interactive Sessions

In API 2.0, the session opening process has been unified and made independent of the operating mode. After obtaining an authentication token, you can open both an interactive session: POST /sessions/online, and a batch session: POST /sessions/batch.

In both cases, a simple JSON is passed containing:

* form code (formCode),

* encrypted AES key for invoice data encryption (encryptionKey).

In the case of batch submission, a list of partial files along with metadata included in the package is also passed.

Details and usage examples:
* interactive submission \[[link](https://github.com/CIRFMF/ksef-docs/blob/main/sesja-interaktywna.md)\]
* batch submission \[[link](https://github.com/CIRFMF/ksef-docs/blob/main/sesja-wsadowa.md)\]

## Mandatory Encryption for All Invoices

In version 1.0, invoice encryption was mandatory only in batch mode. In interactive mode, the encryption option existed but was optional.

In version 2.0, encryption of all invoices – both in batch and interactive mode – **is required**.

Each invoice or invoice package must be encrypted locally by the client using an **AES key**, which:

* is generated individually for each session,

* is passed to the system during session opening (encryptionKey).


## Asymmetric Encryption

In version 2.0, `RSA-OAEP` with `SHA-256` and `MGF1-SHA256` has been introduced. KSeF token encryption is performed with a separate key from the symmetric key encryption used for invoices.

Current **public key certificates** are returned by the public endpoint: GET `/security/public-key-certificates`

## Consistency and New Naming Convention in API 2.0

One of the key changes in KSeF API 2.0 is the unification and simplification of naming conventions for resources, parameters, and JSON models. In version 1.0, the API contained a number of inconsistencies and excessive complexity resulting from system evolution.

In version 2.0:

* **Endpoints** have gained clear, RESTful naming (e.g., sessions/online, auth/token, permissions/entities/grants).

* **Operation names** have been simplified and reflect the actual action (e.g., grant, revoke, refresh).

* The structure of **headers, parameters, and data formats** has been organized to be consistent and aligned with good REST API design practices.

* **Data structures** are flat and clear – identifier and permission types have explicitly defined enum types (Nip, Pesel, Fingerprint), without the need to analyze subtypes.

Changes in version 2.0 also include updates to names and data structures. Although a complete map of these changes is not presented in this document, it is available in the OpenAPI v2 documentation and in the code examples in the official GitHub repository.

It should be emphasized that the changes are not drastic – they **do not affect** the overall logic of the KSeF system, but only organize and simplify naming and structures, making the API more transparent and intuitive to use.

Migration to version 2.0 should be treated as a change in the integration contract and requires adaptation of implementation on the side of external systems. It is recommended to use the official **KSeF 2.0 Client** integration library, developed and maintained by the team responsible for the API. This library implements all available endpoints and data models, which significantly facilitates the migration process and provides stable support for future system versions as well.

## New Module for Internal Certificate Management

As part of KSeF API version 2.0, mechanisms have been introduced to enable the issuance and management of internal **KSeF certificates** \[link to documentation\]. Certificates will enable authentication in KSeF and are necessary for issuing invoices in offline mode \[link to documentation\].

Entities that successfully complete the authentication process will be able to:

* submit an application for an internal KSeF certificate containing selected attributes from the signature certificate used during authentication,

* download the issued certificate in digital form,

* check the status of the submitted certificate application,

* download the list of metadata for issued certificates,

* check the available certificate limit.

## Improvement of the Batch Submission Process

In KSeF API version 2.0, a significant improvement has been introduced in the processing of batch sessions. The previous solution available in API 1.0 was inefficient – if even one invoice in a package contained an error, the entire submission was rejected. This approach effectively limited the use of batch mode by integrators and caused significant operational difficulties.

In the new solution, when submitting a package of documents:

* each invoice is processed independently,

* any errors affect only specific invoices, not the entire submission,

* the number of erroneous invoices is returned for the session status,

* a dedicated endpoint is available to retrieve the detailed status of incorrectly processed invoices along with error information.

This change significantly increases the reliability and efficiency of batch mode and is based on the same package submission model without the risk of losing the entire package due to individual errors.

## Invoice Duplicate Verification
The method of duplicate detection has been changed – business data from the invoice (Podmiot1:NIP, RodzajFaktury, P_2) is now checked, not the file hash. Details – [Duplicate Verification](faktury/weryfikacja-faktury.md).

## Changes in the Permissions Module

Changes in the permissions module are related to changes in some aspects of how they function in KSeF 2.0.

In response to reported comments, version 2.0 of the system introduced a mechanism for granting permissions indirectly, which replaces the previous principle of inheriting permissions for viewing and issuing invoices. The new interface allows for separation of viewing client (partner) invoices and the ability to issue invoices on their behalf from viewing invoices and issuing invoices for the entity itself (e.g., an accounting office).

The mechanism consists of granting a permission to an entity by the client for viewing or issuing invoices with a special option enabled that allows further transfer of this permission by the authorized entity. After receiving such a permission, the entity can grant it to, for example, its employees. After performing these actions, these employees will be able to service the indicated client within the scope defined by the granted permissions.

It is also possible for an entity to grant so-called general permissions, which allow an employee authorized in this way to service all of the entity's clients – of course, to the extent that these clients have authorized this entity and taking into account the scope of the employee's permissions (for viewing and/or issuing invoices).

Thanks to this mechanism, the granting and functioning of permissions for viewing and issuing invoices within the entity itself are not linked to permissions for servicing client invoices. This gives entities better possibilities for profiling employee permissions than the previously functioning inheritance mechanism in KSeF. This was because permissions were granted to employees for viewing and issuing invoices within the entity, and if the entity had appropriate permissions from clients, these permissions automatically passed to employees (were inherited). As a result – an employee could only service clients when they simultaneously had the right to view and/or issue the entity's invoices. And in many cases, this was an excessive permission, which could cause problems in organizations.

Additionally, a new permission type and new login capabilities have been introduced in the system, which enables self-invoicing support by EU entities. It is now possible to log in within a context defined by a Polish entity (identified by NIP) and an entity from an EU country identified by the VAT number of the EU country. In such a defined context, it is possible for authenticated representatives of the indicated EU entity to issue invoices in self-invoicing mode on behalf of the indicated Polish entity.

Within the API definition, all permissions have been organized into logical groups corresponding to individual functional areas of the system.

## API Call Limits (Rate Limiting) ##
In KSeF API version 2.0, a precise and predictable rate limiting mechanism has been introduced, which replaces the previous solutions known from version 1.0.
Each endpoint in the system is subject to a limit on the number of requests in given time intervals: per second, minute, hour.

The ranges and limit values are:
* publicly available: [API limits](limity/limity-api.md),
* differentiated depending on the environment (the test environment has less restrictive limits than production),
* adapted to the nature of the operation:
  * for protected endpoints – limits are applied per context and IP address,
  * for open endpoints – limits are applied per IP address.

The new limit model has been designed not to restrict typical testing of applications integrating with the system.
This solution provides greater transparency, predictability, and better system resilience, both in test and production environments.

## Auxiliary API for Test Data Generation (Test Environment)

In the KSeF API 2.0 test environment, a dedicated **auxiliary API for test data generation** will be made available, enabling quick creation of companies, organizational structures, and contexts necessary for conducting integration tests.

Thanks to this solution, it will be possible to:

* **simulate the establishment of a new business entity**,

* **simulate granting permissions via ZAW-FA**,

* **create units within the JST structure**,

* create tax **VAT groups** (GVAT) along with related entities,

* **define enforcement authorities** **and court bailiffs**.

Typically, the process of registering companies and granting permissions for real entities in the production environment requires formal actions (e.g., a visit to the tax office). In the test environment, such data does not exist. Therefore, the **auxiliary API is an essential tool**, enabling integrators to independently create test entities on which they can freely implement and verify full scenarios of their applications' operation.
