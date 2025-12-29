# **KSeF 2.0 - Integrator's Guide**
22.12.2025

This document serves as a comprehensive knowledge base for developers, analysts, and system integrators implementing integration with the National e-Invoice System (KSeF) version 2.0. The guide focuses on technical and practical aspects related to communication with the KSeF system API.

##  Introduction to KSeF 2.0

The National e-Invoice System (KSeF) is a central ICT system for issuing and retrieving structured invoices in electronic form.

## Table of Contents
The guide has been divided into thematic sections corresponding to key functions and integration areas in the KSeF API:
* [Overview of Key Changes in KSeF 2.0](przeglad-kluczowych-zmian-ksef-api-2-0.md)
* [Changelog](api-changelog.md)
* Authentication
  * [Obtaining Access](uwierzytelnianie.md)
  * [Session Management](auth/sesje.md)
* [Permissions](uprawnienia.md)
* [KSeF Certificates](certyfikaty-KSeF.md)
* Offline Modes
  * [Offline Modes](tryby-offline.md)
  * [Technical Correction](offline/korekta-techniczna.md)
* [QR Codes](kody-qr.md)
* [Interactive Session](sesja-interaktywna.md)
* [Batch Session](sesja-wsadowa.md)
* Downloading Invoices
  * [Downloading Invoices](pobieranie-faktur/pobieranie-faktur.md)
  * [Incremental Invoice Downloading](pobieranie-faktur/przyrostowe-pobieranie-faktur.md)
* [Managing KSeF Tokens](tokeny-ksef.md)
* [Limits](limity/limity.md)
* [Test Data](dane-testowe-scenariusze.md)

For each area, the following are provided:

* detailed description of operation,
* example calls in C# and Java languages,
* references to the [OpenAPI](https://api-test.ksef.mf.gov.pl/docs/v2) specification and reference library code.

The code examples presented in the guide were prepared based on official open source libraries:
* [ksef-client-csharp](https://github.com/CIRFMF/ksef-client-csharp) - library in C#
* [ksef-client-java](https://github.com/CIRFMF/ksef-client-java) - library in Java

Both libraries are maintained and developed by the Ministry of Finance teams and are publicly available under open source terms. They provide full support for KSeF 2.0 API functionality, including support for all endpoints, data models, and example call scenarios. Using these libraries significantly simplifies the integration process and minimizes the risk of incorrect interpretation of API contracts.


## System Environments
The list of KSeF API 2.0 environments is described in the document [KSeF API 2.0 Environments](srodowiska.md)
