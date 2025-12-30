## API Request Limits
22.11.2025

Due to the scale of KSeF operations and its public nature, mechanisms limiting the intensity of API requests have been introduced. Their purpose is to protect system stability, guard against cyber threats, and ensure equal access conditions for all users. Limits define the maximum number of requests that can be made within a specified time and enforce an integration approach that aligns with the system's architectural assumptions.

### General Limit Rules

#### 1. How Limits Are Calculated
All requests to the KSeF API are subject to limits. These restrictions apply to every call to a protected endpoint. For traffic accounting purposes, requests are grouped by pair: context and IP address.

- **context** - defined by the `ContextIdentifier` (`Nip`, `InternalId`, or `NipVatUe`) passed during authentication.
- **IP address** - the IP address from which the network connection is established.

Request limits are calculated independently for each pair: context and IP address. This means that traffic within the same context but from different IP addresses is accounted for separately.

Example
Accounting firm A downloads invoices on behalf of company B, using company B's context (NIP) and connecting to KSeF from IP address IP1.
At the same time, company B downloads invoices independently, within the same context (its own NIP), but from a different IP address - IP2. Despite the shared context, different IP addresses cause limits to be calculated independently.
In this situation, the system treats each connection as a separate pair (context + IP address) and calculates limits independently: separately for accounting firm A and separately for company B.

**Limit Units**
The following notations are used in limit tables:
- req/s - number of requests per second,
- req/min - number of requests per minute,
- req/h - number of requests per hour.

**Limit Calculation Model (sliding/rolling window)**
Limits are enforced using a sliding window model. At any moment, requests made within the following periods are counted:

- for the req/h threshold - in the last 60 minutes,
- for the req/min threshold - in the last 60 seconds,
- for the req/s threshold - in the last second.

Windows are not aligned to full hours or minutes (they do not "reset" at :00). All thresholds (req/s, req/min, req/h) apply simultaneously - blocking is triggered upon the first breach of any of them.

#### 2. API Access Is Blocked When Limits Are Exceeded
When request limits are exceeded, the API returns HTTP code **429 Too Many Requests**, and subsequent requests are temporarily blocked.
The blocking duration is **dynamic** and depends on the frequency and scale of violations. The exact blocking time is returned in the `Retry-After` response header (in seconds). Repeated violations may result in significantly extended blocking periods.

Example 429 response:
```json
HTTP/1.1 429 Too Many Requests
Content-Type: application/json
Retry-After: 30

{
  "status": {
    "code": 429,
    "description": "Too Many Requests",
    "details": [ "Exceeded limit of 20 requests per minute. Try again after 30 seconds." ]
  }
}

```

#### 3. Violation Logging
All cases of limit violations are logged and analyzed by security mechanisms. This data is used to monitor API stability and detect potential abuse.
The system pays particular attention to patterns indicating attempts to circumvent limits, e.g., through parallel and systematic use of multiple IP addresses within a single context. Such activities may be considered a security threat.

In case of repeated violations or extreme load, the system may automatically apply protective measures, such as:
- blocking API access to KSeF for a given entity or IP address range,
- limiting availability for the most resource-intensive contexts.

#### 4. Higher Limits During Night Hours
Between 20:00-06:00, higher download limits apply than during the day.
Specific values will be determined during the initial period of KSeF 2.0 operation, after tuning parameters to actual loads.

#### 5. Preliminary Limit Assumptions
API request limits have been defined based on anticipated system usage scenarios and load models.

The actual nature of traffic will depend on how integrations are implemented in external systems and the load patterns they generate. This means that limits established during the design phase may differ from values maintained in the production environment.

For this reason, limits are dynamic and may be adjusted depending on operational conditions and integrator behavior. In particular, temporary reductions are permitted in cases of intensive or inefficient API usage.


### Limits by Environment

**TE Environment (test)**
On the TE environment, limits have been configured to allow integrators to work freely and test integrations without risk of blocking. Default limit values are **ten times higher** than in production, enabling intensive testing.
Additionally, the following endpoints allow simulation of various scenarios:

* [POST /testdata/rate-limits/production](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Limity-i-ograniczenia/paths/~1api~1v2~1testdata~1rate-limits~1production/post) - activates production (PRD) limits,
* [POST /testdata/rate-limits](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Limity-i-ograniczenia/paths/~1api~1v2~1testdata~1rate-limits/post) - allows setting custom values,
* [DELETE /testdata/rate-limits](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Limity-i-ograniczenia/paths/~1api~1v2~1testdata~1rate-limits/delete) - restores default TE environment limits.

**DEMO Environment (pre-production)**
On the DEMO environment, **the same limits as in production** apply for a given context. These values are **replicated from PRD** and serve for final validation of performance and integration stability before production deployment.

**PRD Environment (production)**
On the PRD environment, **default limits** specified in this documentation are applied.
In justified cases - e.g., large-scale invoice processing - the option to **submit a request for limit increase** is provided through a dedicated form (in preparation).

## Invoice Download - Limits

### Architectural Assumptions
The KSeF API in the invoice download area has been designed as a **document synchronization** mechanism between the central repository and local databases of external systems. A key assumption is that business operations such as searching, filtering, or reporting should be performed on locally stored data that has been previously synchronized with KSeF. This approach increases operational stability, minimizes system overload risk, and allows for more flexible data usage by client applications.

The KSeF API for invoice downloads is not intended for handling direct end-user operations in real-time. This means it should not be used for:
- downloading individual invoices on user request, e.g., invoice preview,
- downloading invoice metadata lists or initiating package exports in response to current application actions, except when the user consciously initiates data synchronization.

### Recommended Integration Approach for Downloads
The `/invoices/query/metadata` endpoint is used for incremental synchronization. Detailed rules for incremental synchronization are described in a separate document.

Depending on invoice volume, different approaches to downloading can be applied:
1. **Low volume scenarios** - if the number of invoices is limited and can be handled within available limits in the expected time, they can be downloaded synchronously by calling `/invoices/ksef/{ksefNumber}` for selected documents.
2. **High volume scenarios** - if the number of documents is significant and synchronous handling becomes impractical, the export mechanism (`/invoices/exports`) is recommended. Export works asynchronously, is queued, and therefore does not negatively affect system performance.
3. **Business operations** - regardless of the chosen strategy, all user actions (searching, filtering, reporting) should be performed on the **local database**, previously synchronized with KSeF.

### Synchronization and Invoice Download Modes
Invoice download to an accounting system can be implemented in three modes:
1. **On user request** - incremental synchronization is initiated **manually** by the user, from the last confirmed checkpoint.
2. **Cyclically** - incremental synchronization is performed automatically according to the system schedule.
3. **Mixed mode** - incremental synchronization runs cyclically, and additionally the user can initiate it manually on request.

### Query Frequency
- **High-frequency schedules are not recommended**. In production environments, the cyclic interval should not be shorter than 15 minutes for each entity appearing on the invoice (Entity 1, Entity 2, Entity 3, Authorized Entity).
- **Low volume profiles.** On-demand downloading is recommended, supplemented with a cycle, e.g., once daily during a night window.
- **Invoice receipt date.** The invoice receipt date is the date when the KSeF number was assigned. This number is assigned automatically by the system at the time of invoice processing and does not depend on when it is downloaded to the accounting system.

### Examples of Unrecommended Implementation
Improper integration can lead to API blocking. The most common errors include:
1. Synchronization only through downloading individual invoices (synchronous path), without using invoice package export.
This approach is acceptable only in low-volume profiles; with larger numbers of documents, the `/invoices/exports` mechanism should be used.
2. Handling end-user requests (e.g., displaying full invoice content in the application, downloading XML files) through direct KSeF API calls instead of using the local database.

### Detailed Limits

| Endpoint | | req/s | req/min | req/h |
|----------|---|-------|---------|-------|
| Get invoice metadata list | POST /invoices/query/metadata | 8 | 16 | 20 |
| Export invoice package | POST /invoices/exports | 4 | 8 | 20 |
| Get invoice package export status | /invoices/exports/{referenceNumber} | 10 | 60 | 600 |
| Get invoice by KSeF number | GET /invoices/ksef/{ksefNumber} | 8 | 16 | 64 |

**Note:** If the available invoice download limits are insufficient for your organization's business scenarios, please contact KSeF support for individual analysis and to find an appropriate solution.

## Invoice Submission - Limits

### Architectural Assumptions
- Invoice submission, regardless of submission type, is queued.
- Processing is optimized for the fastest possible invoice validation confirmation and KSeF number return.

#### Batch Submission (invoice packages):

- An invoice package is treated as a single message in the queue (reference to the package instead of separate entries for each invoice) and processed with the same priority as a single document.
- Batch submission reduces network and operational overhead because:
	- fewer HTTP requests are made,
	- operations on content (decryption, validation, saving) are performed in batches, which is the most efficient way to handle multiple documents simultaneously.
- Batch compression. Due to the XML format and high repeatability of elements between invoices (constant structure, similar field names, repeating blocks), the compression ratio achieved is usually very favorable, significantly reducing data volume and shortening transmission time. In practice, it is faster to send one package containing, e.g., 100 invoices than 100 individual invoices in an interactive session.
- Limits. The limit mechanism works independently of the submission mode. Batch submission inherently reduces the number of requests and facilitates efficient use of available limits.
- Application. Batch mode is recommended wherever more than one document is transmitted in a single operational window. It is particularly effective for cyclic customer settlements, e-commerce, and automated invoicing processes.

Example batch mode usage scenarios:
- **Online store (e-commerce).** Orders and payments are processed asynchronously, and invoices are issued automatically by the ERP system or invoicing module. A single invoice does not need to be sent to KSeF immediately after issuance. A dedicated process can aggregate issued invoices and periodically - e.g., every 5 minutes - send them in batch packages to KSeF, significantly reducing the number of HTTP requests and optimizing limit usage.
- **Subscription services / cyclic settlements.** Invoices are generated collectively once daily or once monthly (e.g., in telecommunications or utilities) and sent in a single package within a scheduled batch session.
- **Automated invoicing processes in enterprises.** These occur, e.g., in the distribution, logistics, manufacturing, or B2B outsourced services sectors. Invoices are generated automatically based on system events (deliveries, order completions) and sent collectively, e.g., after operations are completed.

**Recommendation:** To ensure maximum integration efficiency, it is recommended to aggregate documents in a single batch session wherever possible from a business process perspective. This limits the number of API requests and optimizes the use of available limits.

**Detailed Limits**

| Endpoint | | req/s | req/min | req/h |
|----------|---|-------|---------|-------|
| Open batch session * | POST /sessions/batch | 10 | 20 | 60 |
| Close batch session | POST /sessions/batch/{referenceNumber}/close | 10 | 20 | 60 |

**Package part submission** - requests transmitting package parts within a single batch session are not subject to API limits. For packages divided into multiple parts, parallel (multi-threaded) transmission is recommended, which significantly shortens submission time.

#### Interactive Submission (individual)
Interactive mode has been designed for scenarios requiring fast registration of individual invoices and immediate KSeF number retrieval. Unlike batch sessions, each invoice is sent independently within an active interactive session. Its purpose is to minimize the time needed to obtain a KSeF number for a single document. Application includes low-volume scenarios where individual invoices are sent.

Example interactive mode usage scenarios:
- **Point of sale (POS)**. After completing a transaction, an invoice is issued, and the system registers it immediately in KSeF and returns the KSeF number for printing or presentation to the customer.
- **Mobile applications and lightweight sales systems** that do not have queuing or buffering mechanisms and send invoices immediately after issuance.
- **One-off or irregular events** e.g., a single corrective invoice.

Interactive mode, despite higher network overhead for larger volumes, is an essential complement to batch mode in scenarios requiring immediate response or instant document registration in KSeF. It should be used only where immediate invoice processing is critical to the business process or where the scale of operations does not justify using a batch session.

**Detailed Limits**

| Endpoint | | req/s | req/min | req/h |
|----------|---|-------|---------|-------|
| Open interactive session | POST /sessions/online | 10 | 30 | 120 |
| Send invoice * | POST /sessions/online/{referenceNumber}/invoices | 10 | 30 | 180 |
| Close interactive session | POST /sessions/online/{referenceNumber}/close | 10 | 30 | 120 |

\* **Note:** If your organization's business scenarios regularly reach interactive session submission limits, first consider using batch mode, which allows more efficient use of available resources and limits.
In situations where using an interactive session is necessary and the limits remain insufficient, please contact KSeF support for individual analysis and assistance in selecting a solution.

### Session and Invoice Status

**Detailed Limits**

| Endpoint | | req/s | req/min | req/h |
|----------|---|-------|---------|-------|
| Get invoice status from session | GET /sessions/{referenceNumber}/invoices/{invoiceReferenceNumber} | 30 | 120 | 1200 |
| Get session list | GET /sessions | 5 | 10 | 60 |
| Get session invoices | GET /sessions/{referenceNumber}/invoices | 10 | 20 | 200 |
| Get incorrectly processed session invoices | GET /sessions/{referenceNumber}/invoices/failed | 10 | 20 | 200 |
| Other | GET /sessions/* | 10 | 120 | 1200 |

## Other

Default limits apply to all API resources that do not have specific values defined in this documentation. Each such endpoint has its own limit counter, and its requests do not affect other resources.

These limits apply only to protected resources. They do not cover public API resources such as `/auth/challenge`, which do not require authentication and have their own protection mechanisms - the limit is 60 requests per second for a single IP address.

| Endpoint | | req/s | req/min | req/h |
|----------|---|-------|---------|-------|
| Other | POST/GET /* | 10 | 30 | 120 |

Related documents:
- [Limits](limity.md)
