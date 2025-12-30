## API 2.0 Changes

### Version 2.0.0

- **UPO**
  As announced in RC6.0, from `2025-12-22` UPO version v4-3 is returned by default.

- **Session status** (GET `/sessions/{referenceNumber}`)
  - Extended the response model with `dateCreated` ("Session creation date") and `dateUpdated` ("Date of last activity within the session") properties.

- **Batch session close (POST `/sessions/batch/{referenceNumber}/close`)**
  - Added error code `21208` ("Timeout waiting for upload or finish requests was exceeded").

- **Invoice/UPO download**
  - Added `x-ms-meta-hash` header (`SHA-256` hash, `Base64`) in `200` responses for endpoints:
    - GET `/invoices/ksef/{ksefNumber}`,
    - GET `/sessions/{referenceNumber}/invoices/ksef/{ksefNumber}/upo`,
    - GET `/sessions/{referenceNumber}/invoices/{invoiceReferenceNumber}/upo`,
    - GET `/sessions/{referenceNumber}/upo/{upoReferenceNumber}`.

- **Get authentication status** (GET `/auth/{referenceNumber}`)
  - Updated HTTP 400 (Bad Request) documentation with error code `21304` ("Not authenticated") - authentication operation with reference number {`referenceNumber`} was not found.
  - Extended status `450` ("Authentication failed due to invalid token") with additional cause: "Invalid authorization challenge".

- **Get access tokens** (POST `/auth/token/redeem`)
  Updated HTTP 400 (Bad Request) documentation with error codes:
    - `21301` - "Not authorized":
      - Tokens for operation {`referenceNumber`} have already been retrieved,
      - Authentication status ({`operation.Status`}) does not allow token retrieval,
      - KSeF token has been revoked.
    - `21304` - "Not authenticated" - Authentication operation {`referenceNumber`} was not found,
    - `21308` - "Attempt to use authorization methods of a deceased person".

- **Refresh access token** (POST `/auth/token/refresh`)
  Updated HTTP 400 (Bad Request) documentation with error codes:
    - `21301` - "Not authorized":
      - Authentication status ({`operation.Status`}) does not allow token retrieval,
      - KSeF token has been revoked.
    - `21304` - "Not authenticated" - Authentication operation {`referenceNumber`} was not found,
    - `21308` - "Attempt to use authorization methods of a deceased person".

- **Interactive sending** (POST `/sessions/online/{referenceNumber}/invoices`)
  Updated error code documentation with:
    - `21402` "Invalid file size" - content length does not match file size,
    - `21403` "Invalid file hash" - content hash does not match file hash.

- **Invoice package export (POST `/invoices/exports`). Get invoice metadata list (POST `/invoices/query/metadata`)**
  Reduced maximum allowed `dateRange` range from 2 years to 3 months.

- **Permissions**
  - Added `required` attribute for `subjectDetails` property ("Details of the subject being granted permissions") in all permission granting endpoints (`/permissions/.../grants`).
  - Added `required` attribute for `euEntityDetails` property ("EU entity details in the context of which permissions are granted") in endpoint POST `/permissions/eu-entities/administration/grants` ("Grant EU entity administrator permissions").
  - Added `PersonByFingerprintWithIdentifier` value ("Natural person using a certificate without NIP or PESEL identifier, but having NIP or PESEL") to `EuEntityPermissionSubjectDetailsType` enum in endpoint POST `/permissions/eu-entities/administration/grants` ("Grant EU entity administrator permissions").
  - Changed `subjectEntityDetails` property type to `PermissionsSubjectEntityByIdentifierDetails` ("Authorized subject details") in response model for POST `/permissions/query/authorizations/grants` ("Get list of entity permissions for invoice handling").
  - Changed `subjectEntityDetails` property type to `PermissionsSubjectEntityByFingerprintDetails` ("Authorized subject details") in response model for POST `/permissions/query/eu-entities/grants` ("Get list of EU entity administrator or representative permissions for self-billing").
  - Changed `subjectPersonDetails` property type to `PermissionsSubjectPersonByFingerprintDetails` ("Authorized person details") in response model for POST `/permissions/query/eu-entities/grants` ("Get list of EU entity administrator or representative permissions for self-billing").
  - Introduced checksum validation for `InternalId` identifier in POST `/permissions/subunits/grants` ("Grant subordinate entity administrator permissions").
  - Clarified property descriptions.

- **OpenAPI**
  - Updated `429` response documentation with returned `Retry-After` header and `TooManyRequestsResponse` response body.
  - Clarified `byte` type property descriptions - values are transmitted as binary data encoded in `Base64` format.
  - Fixed typos in specification.

### Version 2.0.0 RC6.1

- **New environment addressing**
  New addresses released. Changes in [KSeF API 2.0 environments](srodowiska.md) section.

- **Authentication - status retrieval (GET `/auth/{referenceNumber}`)**
  Added code `480` - Authentication blocked: "Suspected security incident. Contact the Ministry of Finance through the submission form."

- **Permissions**
  - Extended access rules for session operations (GET/POST `/sessions/...`): added `EnforcementOperations` (enforcement authority) to the list of accepted permissions.
  - Added length restrictions for string type properties: `minLength` and `maxLength`.
  - Added `id` (`Asc`) as second sorting key in `x-sort` metadata for permission search queries. Default order: `dateCreated` (`Desc`), then `id` (`Asc`) - ordering change increasing pagination determinism.
  - Added validation for `IdDocument.country` property in endpoint POST `/permissions/persons/grants` ("Grant permissions to natural persons for KSeF work") - requires compliance with **ISO 3166-1 alpha-2** (e.g., `PL`, `DE`, `US`).
  - "Get list of EU entity administrator or representative permissions for self-billing" (POST `/permissions/query/eu-entities/grants`):
    - removed pattern (regex) validation and clarified description of `EuEntityPermissionsQueryRequest.authorizedFingerprintIdentifier` property.
    - clarified description of `EuEntityPermissionsQueryRequest.vatUeIdentifier` property.

- **Interactive session**
  Added new error codes for POST `/sessions/online/{referenceNumber}/invoices` ("Send invoice"):
    - `21166` - Technical correction unavailable.
    - `21167` - Invoice status does not allow technical correction.

- **API Limits**
  - Increased hourly limit for `invoiceStatus` group (get invoice status from session) from 720 to 1200 req/h:
    - GET /sessions/{referenceNumber}/invoices/{invoiceReferenceNumber}.
  - Increased hourly limit for `sessionMisc` group (GET `/sessions/*` resources) from 720 to 1200 req/h:
    - GET `/sessions/{referenceNumber}`,
    - GET `/sessions/{referenceNumber}/invoices/ksef/{ksefNumber}/upo`,
    - GET `/sessions/{referenceNumber}/invoices/{invoiceReferenceNumber}/upo`,
    - GET `/sessions/{referenceNumber}/upo/{upoReferenceNumber}`.
  - Decreased hourly limit for `batchSession` group (batch session open/close) from 120 to 60 req/h:
    - POST `/sessions/batch`,
    - POST `/sessions/batch/{referenceNumber}/close`.
  - Increased limits for endpoint `/invoices/exports/{referenceNumber}` ("Get invoice package export status") by adding new group `invoiceExportStatus` with parameters: 10 req/s, 60 req/min, 600 req/h.

- **Batch session open (POST `/sessions/batch`)**
  Removed `fileName` property from `BatchFilePartInfo` model (previously marked as deprecated; x-removal-date: 2025-12-07).

- **Authentication initialization (POST `/auth/challenge`)**
  Added `timestampMs` property (int64) in response model - challenge generation time in milliseconds since 1.01.1970 (Unix).

- **Test data**
  - Clarified `expectedEndDate` property type: format: `date` in (POST `/testdata/attachment/revoke`).
  - Removed `Token` value from `SubjectIdentifierType` enum in endpoint POST `/testdata/limits/subject/certificate`. The value was unused: in KSeF a subject cannot be a "token" - identity always derives from `NIP/PESEL` or certificate fingerprint, which carries the identity of the subject who created it.

- **OpenAPI**
  Increased maximum `pageSize` value from 500 to 1000 for endpoints:
  - GET `/sessions`
  - GET `/sessions/{referenceNumber}/invoices`
  - GET `/sessions/{referenceNumber}/invoices/failed`

### Version 2.0.0 RC6.0

- **API Limits**
  - On **TE** (test) environment, enabled and defined [api limits](limity/limity-api.md) policy with values 10x higher than on **PRD**; details: ["Limits on environments"](/limity/limity-api.md#limity-na-środowiskach).
  - On **TR** (DEMO) environment, enabled [api limits](limity/limity-api.md) with values identical to **PRD**. Values are replicated from production; details: ["Limits on environments"](/limity/limity-api.md#limity-na-środowiskach).
  - Added endpoint POST `/testdata/rate-limits/production` - sets api limit values matching production profile in current context. Available only on **TE** environment.

- **Invoice package export (POST `/invoices/exports`). Get invoice metadata list (POST `/invoices/query/metadata`)**
  - Added [High Water Mark (HWM)](pobieranie-faktur/hwm.md) document describing the mechanism for managing data completeness over time.
  - Updated [Incremental invoice retrieval](pobieranie-faktur/przyrostowe-pobieranie-faktur.md) with recommendations for using the `HWM` mechanism.
  - Extended response model with `permanentStorageHwmDate` property (string, date-time, nullable). Applies only to queries with `dateType = PermanentStorage` and indicates the point below which data is complete; for `dateType = Issue/Invoicing` - null.
  - Added `restrictToPermanentStorageHwmDate` property (boolean, nullable) in `dateRange` object, which enables the High Water Mark (`HWM`) mechanism and restricts the date range to the current `HWM` value. Applies only to queries with `dateType = PermanentStorage`. Using this parameter significantly reduces duplicates between consecutive exports and ensures consistency during long-running incremental synchronization.

- **UPO - XSD update to v4-3**
  - Changed `NumerKSeFDokumentu` element pattern to also allow KSeF numbers generated for invoices from KSeF 1.0 (36 characters).
  - Added `TrybWysylki` element - document submission mode to KSeF: `Online` or `Offline`.
  - Changed `NazwaStrukturyLogicznej` value to format: Schemat_{systemCode}_v{schemaVersion}.xsd (e.g., Schemat_FA(3)_v1-0E.xsd).
  - Changed `NazwaPodmiotuPrzyjmujacego` value on test environments by adding environment name suffix:
    - `TE`: Ministry of Finance - test environment (TE),
    - `TR`: Ministry of Finance - pre-production environment (TR).

    `PRD`: unchanged - Ministry of Finance.
  - Currently UPO v4-2 is returned by default. To receive UPO v4-3, add header: `X-KSeF-Feature: upo-v4-3` when opening session (online/batch).
  - From `2025-12-22` UPO v4-3 will be the default version.
  - XSD UPO v4-3: [schema](/faktury/upo/schemy/upo-v4-3.xsd).

- **Session status** (GET `/sessions/{referenceNumber}`)
    Clarified description of code `440` - Session cancelled: possible causes are "Sending time exceeded" or "No invoices sent".

- **Invoice status** (GET `/sessions/{referenceNumber}/invoices/{invoiceReferenceNumber}`)
    Added `InvoiceStatusInfo` type (extends `StatusInfo`) with `extensions` field - object with structured status details. The `details` field remains unchanged. Example (duplicate invoice):

    ```json
    "status": {
      "code": 440,
      "description": "Duplicate invoice",
      "details": [
        "Duplicate invoice. Invoice with KSeF number: 5265877635-20250626-010080DD2B5E-26 has already been correctly sent to the system in session: 20250626-SO-2F14610000-242991F8C9-B4"
      ],
      "extensions": {
        "originalSessionReferenceNumber": "20250626-SO-2F14610000-242991F8C9-B4",
        "originalKsefNumber": "5265877635-20250626-010080DD2B5E-26"
      }
    }
    ```

- **Permissions**
    Added `subjectDetails` property - "Details of the subject being granted permissions" to all permission granting endpoints (/permissions/.../grants). In RC6.0 the field is optional; from 2025-12-19 it will be required.

- **Granted permissions search** (POST `/permissions/query/authorizations/grants`)
    Extended access rules with `PefInvoiceWrite`.

- **Test data - attachments (POST /testdata/attachment/revoke)**
  Extended request model `AttachmentPermissionRevokeRequest` with `expectedEndDate` field (optional) - date of consent withdrawal for sending invoices with attachments.

- **OpenAPI**
  - Added HTTP `429` - "Too Many Requests" response to all endpoints. The `description` property of this response includes tabular presentation of limits (`req/s`, `req/min`, `req/h`) and the limit group name assigned to the endpoint. The mechanism and semantics of `429` remain consistent with the description in the [limits](/limity/limity-api.md) documentation.
  - Added `x-rate-limits` metadata with limit values (`req/s`, `req/min`, `req/h`) to each endpoint.
  - Removed explicit `exclusiveMaximum`: `false` and `exclusiveMinimum`: `false` properties from numeric definitions (only minimum/maximum retained). Cleanup change - no impact on validation (in OpenAPI default values of these properties are `false`).
  - Added length restrictions for string type properties: `minLength`.
  - Removed explicit `style`: `form` settings for in: query parameters.
  - Changed order of `BuyerIdentifierType` enum values (now: `None`, `Other`, `Nip`, `VatUe`). Ordering change - no impact on functionality.
  - Fixed typo in `KsefNumber` property description.
  - Clarified format of `PublicKeyCertificate` property representing `Base64` encoded binary data, set format: `byte`.
  - Made minor linguistic and punctuation corrections in `description` fields.

### Version 2.0.0 RC5.7

- **Batch session open (POST `/sessions/batch`)**
  Marked `BatchFilePartInfo.fileName` as `deprecated` in request model (planned removal: 2025-12-05).

- **Asynchronous operation statuses**
  Added status `550` - "Operation was cancelled by the system". Description: "Processing was interrupted due to internal system reasons. Please try again."

- **OpenAPI**
  - Added array element count restrictions: `minItems`, `maxItems`.
  - Added length restrictions for string type properties: `minLength` and `maxLength`.
  - Updated property descriptions (`invoiceMetadataAuthorizedSubject.role`, `invoiceMetadataBuyer`, `invoiceMetadataThirdSubject.role`, `buyerIdentifier`).
  - Updated regex patterns for `vatUeIdentifier`, `authorizedFingerprintIdentifier`, `internalId`, `nipVatUe`, `peppolId`.

### Version 2.0.0 RC5.6

- **Get session status (GET `/sessions/{referenceNumber}`)**
  Added `UpoPageResponse.downloadUrlExpirationDate` field in response - date and time of UPO download URL expiration; after this moment `downloadUrl` is no longer active.

- **Get certificate metadata list (POST `/certificates/query`)**
  Extended response (`CertificateListItem`) with `requestDate` property - certification request submission date.

- **Get Peppol service providers list (GET `/peppol/query`)**
  - Extended response model with `dateCreated` field - Peppol service provider registration date in the system.
  - Marked `dateCreated`, `id`, `name` properties in response model as always returned.
  - Defined `PeppolI` schema (string, 9 characters) and applied in `PeppolProvider`.

- **OpenAPI**
  - Added `x-sort` metadata to all endpoints returning lists. Added Sorting section in endpoint descriptions with default order (e.g., "requestDate (Desc)").
  - Added length restrictions for string type properties: `minLength` and `maxLength`.
  - Clarified format of properties representing `Base64` encoded binary data: set format: `byte` (`encryptedInvoiceContent`, `encryptedSymmetricKey`, `initializationVector`, `encryptedToken`).
  - Defined common `Sha256HashBase64` schema and applied it to all properties representing `SHA-256` hash in `Base64` (including `invoiceHash`).
  - Defined common `ReferenceNumber` schema (string, length 36) and applied it to all parameters and properties representing asynchronous operation reference number (in paths, queries and responses).
  - Defined common `Nip` schema (string, 10 characters, regex pattern) and applied it to all properties representing NIP.
  - Defined `Pesel` schema (string, 11 characters, regex pattern) and applied it in property representing PESEL.
  - Defined common `KsefNumber` schema (string, 35-36 characters, regex pattern) and applied it to all properties representing KSeF number.
  - Defined `Challenge` schema (string, 36 characters) and applied in `AuthenticationChallengeResponse`.`challenge`.
  - Defined common `PermissionId` schema (string, 36 characters) and applied it everywhere: in parameters and response properties.
  - Added regular expressions for selected text fields.

### Version 2.0.0 RC5.5

- **Get current API limits (GET `/api/v2/rate-limits`)**
  Added endpoint returning effective API call limits in `perSecond`/`perMinute`/`perHour` layout for individual areas (including `onlineSession`, `batchSession`, `invoiceSend`, `invoiceStatus`, `invoiceExport`, `invoiceDownload`, `other`).

- **Invoice status in session**
  Extended response for GET `/sessions/{referenceNumber}/invoices` ("Get session invoices") and GET `/sessions/{referenceNumber}/invoices/{invoiceReferenceNumber}` ("Get invoice status from session") with properties: `upoDownloadUrlExpirationDate` - "date and time of URL expiration. After this date the `UpoDownloadUrl` link will no longer be active". Extended `upoDownloadUrl` description.

- **Removal of \*InMib fields (change consistent with 5.3 announcement)**
  Removed `maxInvoiceSizeInMib` and `maxInvoiceWithAttachmentSizeInMib` properties.
  Change affects:
    - GET `/limits/context` - responses (`onlineSession`, `batchSession`),
    - POST `/testdata/limits/context/session` - request model (`onlineSession`, `batchSession`),
    - Models: `BatchSessionContextLimitsOverride`, `BatchSessionEffectiveContextLimits`, `OnlineSessionContextLimitsOverride`, `OnlineSessionEffectiveContextLimits`.
  Only *InMB fields (1 MB = 1,000,000 B) are used for indicating sizes.

- **Removal of `operationReferenceNumber` (change consistent with 5.3 announcement)**
  Removed `operationReferenceNumber` property from response model; the only valid name is `referenceNumber`. Change covers:
  - GET `/invoices/exports/{referenceNumber}` - "Invoice package export status",
  - POST `/permissions/operations/{referenceNumber}` - "Get permission operation status".

- **Invoice package export (POST `/invoices/exports`)**
  - Added new error code: `21182` - "Reached concurrent exports limit. For authenticated subject in current context, maximum limit of {count} concurrent invoice exports has been reached. Please try again later".
  - Extended response model with `packageExpirationDate` property indicating package expiration date. After this date the package will not be available for download.
  - Added error code `210` - "Invoice export has expired and is no longer available for download".

- **Invoice package export status (GET `/invoices/exports/{referenceNumber}`)**
  Clarified descriptions of package part download link fields:
  - `url` - "URL to which the package part download request should be sent. The link is generated dynamically at the time of querying the export operation status. Not subject to API limits and does not require access token when downloading".
  - `expirationDate` - "Date and time of link expiration for downloading the package part. After this moment the link becomes inactive".

- **Authorization**
  - Extended access rules with `SubunitManage` for POST `/permissions/query/persons/grants`: operation can be performed if subject has `CredentialsManage`, `CredentialsRead`, `SubunitManage`.
  - Indirect permission granting (POST `/permissions/indirect/grants`)
    Updated `targetIdentifier.description` property description: clarified that absence of context identifier means granting general indirect permission.

- **OpenAPI**
  Increased maximum `pageSize` value from 100 to 500 for endpoints:
  - GET `/sessions`
  - GET `/sessions/{referenceNumber}/invoices`
  - GET `/sessions/{referenceNumber}/invoices/failed`

### Version 2.0.0 RC5.4

- **Get invoice metadata list (POST /invoices/query/metadata)**
  - Added `sortOrder` parameter, allowing specification of result sorting direction.

- **Session status**
  Fixed bug preventing this property from being populated in API responses for invoices (field was not previously returned). Value is populated asynchronously at time of permanent storage and may be temporarily null.

- **Test data (test environments only)**
  - Change API limits for current context (POST `testdata/rate-limits`)
  Added endpoint allowing temporary override of effective API limits for current context. Change prepares limit simulator launch on TE environment.
  - Restore default limits (DELETE `/testdata/rate-limits`)
  Added endpoint restoring default limit values for current context.

- **OpenAPI**
  - Clarified array parameter definitions in query; applied `style: form`. Multiple values should be passed by repeating the parameter, e.g., `?statuses=InProgress&statuses=Succeeded`. Documentation change, no impact on API behavior.
  - Updated property descriptions (`partUploadRequests`, `encryptedSymmetricKey`, `initializationVector`).

### Version 2.0.0 RC5.3

- **Invoice package export (POST `/invoices/exports`)**
  Added ability to include `_metadata.json` file in export package. File has JSON object format with `invoices` array containing `InvoiceMetadata` objects (model returned by POST `/invoices/query/metadata`).
  Enable (preview): add `X-KSeF-Feature`: `include-metadata` to request header.
  From 2025-10-27 default endpoint behavior changes - export package will always contain `_metadata.json` file (header will not be required).

- **Invoice status**
  - In case of processing with error, when invoice number could be read (e.g., error code `440` - duplicate invoice), response contains `invoiceNumber` property with the read number.
  - Marked `invoiceHash`, `referenceNumber` properties in response model as always returned.

- **Size unit standardization (MB, SI)**
  Unified limit notation in documentation and API: values presented in MB (SI), where 1 MB = 1,000,000 B.

- **Get limits for current context (GET `/limits/context`)**
  Added `maxInvoiceSizeInMB`, `maxInvoiceWithAttachmentSizeInMB` in response model for `onlineSession` and `batchSession` properties.
  `maxInvoiceSizeInMib`, `maxInvoiceWithAttachmentSizeInMib` properties marked as deprecated (planned removal: 2025-10-27).

- **Change session limits for current context (POST `/testdata/limits/context/session`)**
  Added `maxInvoiceSizeInMB`, `maxInvoiceWithAttachmentSizeInMB` in request model for `onlineSession` and `batchSession` properties.
  `maxInvoiceSizeInMib`, `maxInvoiceWithAttachmentSizeInMib` properties marked as deprecated (planned removal: 2025-10-27).

- **Invoice package export status (GET `/invoices/exports/{referenceNumber}`)**
  Changed path parameter name from `operationReferenceNumber` to `referenceNumber`.
  Change does not affect HTTP contract (path and value meaning unchanged) nor endpoint behavior.

- **Permissions**
  - Updated endpoint descriptions and examples in permissions/* area. Change affects documentation only (clarification of descriptions, formats and examples); no changes to API behavior or contract.
  - Changed path parameter name from `operationReferenceNumber` to `referenceNumber` in "Get operation status" (POST `/permissions/operations/{referenceNumber}`).
  Change does not affect HTTP contract (path and value meaning unchanged) nor endpoint behavior.
  - "Grant indirect permissions" (POST `permissions/indirect/grants`)
    Added internal identifier support - extended `targetIdentifier` property with `InternalId` value.
  - "Get own permissions list" (POST `/permissions/query/personal/grants`)
      - Extended `targetIdentifier` property in request model with `InternalId` value (ability to specify internal identifier).
      - Removed `PersonalPermissionScope.Owner` value from response model. Owner permissions (granted by ZAW-FA or NIP/PESEL association) are not returned.

- **Authentication status (GET `/auth/{referenceNumber}`)**
  Extended error codes table with `470` - "Authentication failed" with clarification: "Attempt to use authorization methods of a deceased person".

- **PEF invoice handling**
  Changed enum values (`FormCode`):
    - `FA_PEF (3)` to `PEF (3)`,
    - `FA_KOR_PEF (3)` to `PEF_KOR (3)`.

- **Generate new token (POST `/tokens`)**
  - In request model (`GenerateTokenRequest`) marked `description` and `permissions` fields as required.
  - In response model (`GenerateTokenResponse`) marked `referenceNumber` and `token` fields as always returned.

- **KSeF token status (GET /tokens/{referenceNumber})**
  - Marked `authorIdentifier`, `contextIdentifier`, `dateCreated`, `description`, `referenceNumber`, `requestedPermissions`, `status` properties in response model as always returned.

- **Get generated tokens list (GET /tokens)**
  - Marked `authorIdentifier`, `contextIdentifier`, `dateCreated`, `description`, `referenceNumber`, `requestedPermissions`, `status` properties in response model as always returned.

- **Test data - create natural person (POST `/testdata/person`)**
  Extended request with `isDeceased` property (boolean) enabling creation of test deceased person (e.g., for scenarios verifying status code `470`).

- **OpenAPI**
  - Clarified restrictions for integer type properties in requests by adding `minimum` / `exclusiveMinimum`, `maximum` / `exclusiveMaximum` attributes.
  - Extended response with `referenceNumber` field (contains same value as existing `operationReferenceNumber`). Marked `operationReferenceNumber` as `deprecated` and will be removed from response on 2025-10-27; migrate to `referenceNumber`. Change nature: transitional rename with backward compatibility (both properties returned in parallel until removal date).
  Affects endpoints:
    - POST `/permissions/persons/grants`,
    - POST `/permissions/entities/grants`,
    - POST `/permissions/authorizations/grants`,
    - POST `/permissions/indirect/grants`,
    - POST `/permissions/subunits/grants`,
    - POST `/permissions/eu-entities/administration/grants`,
    - POST `/permissions/eu-entities/grants`,
    - DELETE `/permissions/common/grants/{permissionId}`,
    - DELETE `/permissions/authorizations/grants/{permissionId}`,
    - POST `/invoices/exports`.
  - Removed `required` attribute from `pageSize` property in GET `/sessions` ("Get sessions list") request.
  - Updated examples in endpoint definitions.

### Version 2.0.0 RC5.2
- **Permissions**
  - "Grant subordinate entity administrator permissions" (POST `/permissions/subunits/grants`)
  Added `subunitName` property ("Subordinate unit name") in request. Field is required when subordinate unit is identified by internal identifier.
  - "Get subordinate units and entities administrator permissions list" (POST `/permissions/query/subunits/grants`)
  Added `subunitName` property ("Subordinate unit name") in response.
  - "Get permissions list for KSeF work granted to natural persons or entities" (POST `permissions/query/persons/grants`)
    Removed `Owner` permission type from results. `Owner` permission is assigned systemically to a natural person and is not subject to granting, so should not appear on the list of granted permissions.
  - "Get own permissions list" (POST `/permissions/query/personal/grants`)
    Extended `PersonalPermissionType` filter enum with `VatUeManage` value.

- **Limits**
  - Added endpoints for checking configured limits (context, authenticated subject):
    - GET `/limits/context`
    - GET `/limits/subject`
  - Added endpoints for managing limits (context, authenticated subject) in test environment:
    - POST/DELETE `/testdata/limits/context/session`
    - POST/DELETE `/testdata/limits/subject/certificate`
  - Updated [Limits](limity/limity.md).

- **Invoice status**
  Added `invoicingMode` property in response model. Updated documentation: [Automatic offline mode determination](offline/automatyczne-okreslanie-trybu-offline.md).

- **OpenAPI**
  - Clarified restrictions for integer type properties in requests by adding `minimum` / `exclusiveMinimum`, `maximum` / `exclusiveMaximum` attributes and `default` values.
  - Updated examples in endpoint definitions.
  - Clarified endpoint descriptions.
  - Added `required` attribute for required properties in requests and responses.

### Version 2.0.0 RC5.1

- **Get certificate metadata list (POST /certificates/query)**
  Changed subject identifier representation from property pair `subjectIdentifier` + `subjectIdentifierType` to composite object `subjectIdentifier` { `type`, `value` }.

- **Get invoice metadata list (POST /invoices/query/metadata)**
  - Changed selected identifier representations from type + value property pairs to composite objects { type, value }:
    - `invoiceMetadataBuyer.identifier` + `invoiceMetadataBuyer.identifierType` to composite object `invoiceMetadataBuyerIdentifier` { `type`, `value` },
    - `invoiceMetadataThirdSubject.identifier` + `invoiceMetadataThirdSubject.identifierType` to composite object `InvoiceMetadataThirdSubjectIdentifier` { `type`, `value` }.
  - Removed `obsoleted` `Identitifer` properties from `InvoiceMetadataSeller` and `InvoiceMetadataAuthorizedSubject` objects.
  - Changed `invoiceQuerySeller` property to `sellerNip` in request filter.
  - Changed `invoiceQueryBuyer` property to `invoiceQueryBuyerIdentifier` with properties { `type`, `value` } in request filter.

- **Permissions**
  Changed selected identifier representations from type + value property pairs to composite objects { type, value }:
    - "Get own permissions list" (POST `/permissions/query/personal/grants`):
      - `contextIdentifier` + `contextIdentifierType` -> `contextIdentifier` { `type`, `value` },
      - `authorizedIdentifier` + `authorizedIdentifierType` -> `authorizedIdentifier` { `type`, `value` },
      - `targetIdentifier` + `targetIdentifierType` -> `targetIdentifier` { type, value }.
    - "Get permissions list for KSeF work granted to natural persons or entities" (POST `/permissions/query/persons/grants`),
      - `contextIdentifier` + `contextIdentifierType` -> `contextIdentifier` { `type`, `value` },
      - `authorizedIdentifier` + `authorizedIdentifierType` -> `authorizedIdentifier` { `type`, `value` },
      - `targetIdentifier` + `targetIdentifierType` -> `targetIdentifier` { `type`, `value` },
      - `authorIdentifier` + `authorIdentifierType` -> `authorIdentifier` { `type`, `value` }.
    - "Get subordinate units and entities administrator permissions list" (POST `/permissions/query/subunits/grants`):
      - `authorizedIdentifier` + `authorizedIdentifierType` -> `authorizedIdentifier` { `type`, `value` },
      - `subunitIdentifier` + `subunitIdentifierType` -> `subunitIdentifier` { `type`, `value` },
      - `authorIdentifier` + `authorIdentifierType` -> `authorIdentifier` { `type`, `value` }.
    - "Get entity roles list" (POST `/permissions/query/entities/roles`):
      - `parentEntityIdentifier` + `parentEntityIdentifierType` -> `parentEntityIdentifier` { `type`, `value` }.
    - "Get subordinate entities list" (POST `/permissions/query/subordinate-entities/roles`):
      - `subordinateEntityIdentifier` + `subordinateEntityIdentifierType` -> `subordinateEntityIdentifier` { `type`, `value` }.
    - "Get entity permissions list for invoice handling" (POST `/permissions/query/authorizations/grants`):
      - `authorizedEntityIdentifier` + `authorizedEntityIdentifierType` -> `authorizedEntityIdentifier` { `type`, `value` },
      - `authorizingEntityIdentifier` + `authorizingEntityIdentifierType` -> `authorizingEntityIdentifier` { `type`, `value` },
      - `authorIdentifier` + `authorIdentifierType` -> `authorIdentifier` { `type`, `value` }
    - "Get EU entity administrator or representative permissions list for self-billing" (POST `/permissions/query/eu-entities/grants`):
      - `authorIdentifier` + `authorIdentifierType` -> `authorIdentifier` { `type`, `value` }

- **Grant EU entity administrator permissions (POST permissions/eu-entities/administration/grants)**
  Changed property name in request from `subjectName` to `euEntityName`.

- **Authentication using KSeF token**
  Removed redundant enum values `None`, `AllPartners` in `contextIdentifier.type` property of POST `/auth/ksef-token` request.

- **KSeF Tokens**
  - Unified GET `/tokens` response model: `authorIdentifier.type`, `authorIdentifier.value`, `contextIdentifier.type`, `contextIdentifier.value` properties are always returned (required, non-nullable),
  - Removed redundant enum values `None`, `AllPartners` in `authorIdentifier.type` and `contextIdentifier.type` properties in GET `/tokens` ("Get generated tokens list") response model.

- **Batch session**
  Removed redundant error code `21401` - "Document does not conform to schema (json)".

- **Get session status (GET /sessions/{referenceNumber})**
  - Added error code `420` - "Session invoice limit exceeded".

- **Get invoice metadata (GET `/invoices/query/metadata`)**
  - Added (always returned) `isTruncated` property (boolean) in response - "Indicates whether result was truncated due to invoice count limit (10,000) being exceeded",
  - Marked `amount.type` property in request filter as required.

- **Invoice package export: initiate (POST `/invoices/exports`)**
  - Marked `operationReferenceNumber` property in response model as always returned,
  - Marked `amount.type` property in request filter as required.

- **Get permissions list for KSeF work granted to natural persons or entities (POST /permissions/query/persons/grants)**
  - Added `contextIdentifier` in request filter and response model.

- **OpenAPI**
  Removed unused `operationId` from specification. Cleanup change.

### Version 2.0.0 RC5

- **PEF invoice and Peppol service provider handling**
  - Added support for `PEF` invoices sent by Peppol service provider. New capabilities do not change existing KSeF behavior for other formats, they are API extensions.
  - Introduced new authentication context type: `PeppolId`, enabling work in Peppol service provider context.
  - Automatic provider registration: on first Peppol service provider authentication (using dedicated certificate), automatic system registration occurs.
  - Added GET `/peppol/query` ("Peppol service providers list") endpoint returning registered providers.
  - Updated access rules for session opening and closing, invoice sending requires `PefInvoiceWrite` permission.
  - Added new invoice schemas: `FA_PEF (3)`, `FA_KOR_PEF (3)`,
  - Extended `ContextIdentifier` with `PeppolId` in xsd `AuthTokenRequest`.

- **UPO**
  - Added `Uwierzytelnienie` (Authentication) element, which organizes UPO header data and extends it with additional information; replaces previous `IdentyfikatorPodatkowyPodmiotu` and `SkrotZlozonejStruktury`.
  - `Uwierzytelnienie` contains:
    - `IdKontekstu` - authentication context identifier,
    - authentication proof (depending on method):
      - `NumerReferencyjnyTokenaKSeF` - KSeF authenticating token identifier in system,
      - `SkrotDokumentuUwierzytelniajacego` - hash value of authenticating document in form received by system (including electronic signature).
  - Added to `Dokument` element:
    - NipSprzedawcy (SellerNIP),
    - DataWystawieniaFaktury (InvoiceIssueDate),
    - DataNadaniaNumeruKSeF (KSeFNumberAssignmentDate).
  - Unified UPO schema. Invoice UPO and session UPO use common schema [upo-v4-2.xsd](/faktury/upo/schemy/upo-v4-2.xsd). Replaces previous upo-faktura-v3-0.xsd and upo-sesja-v4-1.xsd.

- **API request limits**
  Added [API request limits](limity/limity-api.md) specification.

- **Authentication**
  - Clarified status codes in GET `/auth/{referenceNumber}`, `/auth/sessions`:
    - `415` (no permissions),
    - `425` (authentication revoked),
    - `450` (invalid token: incorrect token, incorrect time, revoked, inactive),
    - `460` (certificate error: invalid, chain verification error, untrusted chain, revoked, incorrect).
  - Update of optional IP policy in XSD `AuthTokenRequest`:
    Replaced `IpAddressPolicy` with new `AuthorizationPolicy`/`AllowedIps` structure. Updated [Authentication](uwierzytelnianie.md) document.

- **Authorization**
  - Extended access rules with `VatUeManage`, `SubunitManage` for DELETE `/permissions/common/grants/{permissionId}`: operation can be performed if subject has `CredentialsManage`, `VatUeManage` or `SubunitManage`.
  - Extended access rules with `Introspection` for GET `/sessions/{referenceNumber}/...`: each of these endpoints can now be called with `InvoiceWrite` or `Introspection`.
  - Extended access rules with `InvoiceWrite` for GET `/sessions` ("Get sessions list"): with `InvoiceWrite` permission, only sessions created by authenticating subject can be retrieved; with `Introspection` permission, all sessions can be retrieved.
  - Changed access rules for DELETE `/tokens/{referenceNumber}`: removed `CredentialsManage` permission requirement.

- **Get certification request data (GET `certificates/enrollments/data`)**
  - Response structure change:
    - Removed: givenNames (string array).
    - Added: givenName (string).
    - Change nature: breaking (name and type change from array to text).
  - Added error code `25011` - "Invalid CSR signature algorithm".
  - Clarified requirements for private key used for CSR signing in [KSeF Certificates](certyfikaty-KSeF.md).

- **KSeF Tokens**
  - Added error code for POST `/tokens` ("Generate new token") response: `26002` - "Cannot generate token for current context type". Token can only be generated in `Nip` or `InternalId` context.
  - Extended permission catalog assignable to token: added `SubunitManage` and `EnforcementOperations`.
  - Added query parameters for filtering GET `/tokens` results:
    - `description` - search in token description (case-insensitive), min. 3 characters,
    - `authorIdentifier` - search by creator identifier (case-insensitive), min. 3 characters,
    - `authorIdentifierType` - creator identifier type used with authorIdentifier (Nip, Pesel, Fingerprint).
  - Added properties
    - `lastUseDate` - "Token last use date",
    - `statusDetails` - "Additional status information, returned in case of errors"
    in responses for:
    - GET `/tokens` ("token list"),
    - GET `/tokens/{referenceNumber}` ("token status").

- **Get invoice metadata (GET `/invoices/query/metadata`)**
  - Filters:
    - pagination: increased maximum page size to 250 records,
    - removed `schemaType` property (with values `FA1`, `FA2`, `FA3`), previously marked as deprecated,
    - added `seller.nip`; `seller.identifier` marked as deprecated (will be removed in next release),
    - added `authorizedSubject.nip`; `authorizedSubject.identifier` marked as deprecated (will be removed in next release),
    - clarified description: missing value in `dateRange.to` means current date and time (UTC) is used,
    - clarified maximum allowed `DateRange` range as 2 years.
  - Sorting:
    - results are sorted ascending by date type specified in `DateRange`; recommended type for incremental retrieval is `PermanentStorage`,
  - Response model:
    - removed `totalCount` property,
    - changed name from `fileHash` to `invoiceHash`,
    - added `seller.nip`; `seller.identifier` marked as deprecated (will be removed in next release),
    - added `authorizedSubject.nip`; `authorizedSubject.identifier` marked as deprecated (will be removed in next release),
    - marked `invoiceHash` as always returned,
    - marked `invoicingMode` as always returned,
    - marked `authorizedSubject.role` ("Authorized subject") as always returned,
    - marked `invoiceMetadataAuthorizedSubject.role` ("Authorized subject NIP") as always returned,
    - marked `invoiceMetadataThirdSubject.role` ("Third subjects list") as always returned.
  - Removed [Mock] labels from property descriptions.

- **Invoice package export: initiate (POST `/invoices/exports`)**
  - Filters:
    - added `seller.nip`; `seller.identifier` marked as deprecated (will be removed in next release),
  - Removed [Mock] labels.
  - Changed error code: from `21180` to `21181` ("Invalid invoice export request").
  - Clarified sorting rules. Invoices in package are sorted ascending by date type specified in `DateRange` during export initialization.

  - **Invoice package export: status (GET `/invoices/exports/{operationReferenceNumber}`)**
    - Status descriptions: updated export status documentation:
      - `100` - "Invoice export in progress"
      - `200` - "Invoice export completed successfully"
      - `415` - "Delivered key decryption error"
      - `500` - "Unknown error ({statusCode})"
    - Response model `package`:
      - added:
        - `invoiceCount` - "Total invoice count in package. Maximum invoice count in package is 10,000",
        - `size` - "Package size in bytes. Maximum package size is 1 GiB (1,073,741,824 bytes)",
        - `isTruncated` - "Indicates whether export result was truncated due to invoice count limit or package size being exceeded",
        - `lastIssueDate` - "Issue date of last invoice included in package.\nField appears only when package was truncated and export was filtered by `Issue` date type",
        - `lastInvoicingDate` - "Acceptance date of last invoice included in package.\nField appears only when package was truncated and export was filtered by `Invoicing` date type",
        - `lastPermanentStorageDate` - "Permanent storage date of last invoice included in package.\nField appears only when package was truncated and export was filtered by `PermanentStorage` date type".
    - Response model `package.parts`
      - removed `fileName`, `headers`,
      - added:
        - `partName` - "Package part file name",
        - `partSize` - "Package part size in bytes. Maximum part size is 50MiB (52,428,800 bytes)",
        - `partHash` - "SHA256 hash of package part file, encoded in Base64 format",
        - `encryptedPartSize` - "Encrypted package part size in bytes",
        - `encryptedPartHash` - "SHA256 hash of encrypted package part, encoded in Base64 format",
        - `expirationDate` - "Part download link expiration moment",
      - marked all properties in `package` as always returned,
    - Removed [Mock] labels.

- **Permissions**
  - Extended POST `/permissions/eu-entities/administration/grants` ("Grant EU entity administrator permissions") request with "Subject name" `subjectName`.
  - Extended POST `/permissions/query/persons/grants` request with new `System` value for granting subject identifier filter `authorIdentifier` and removed requirement from `authorIdentifier.value` field.
  - Extended POST `/permissions/query/persons/grants` request with new `AllPartners` value for target subject identifier filter `targetIdentifier` and removed requirement from `targetIdentifier.value` field.
  - Added POST `/permissions/query/personal/grants` request for getting own permissions list.
  - Added new `AllPartners` value to POST `/permissions/indirect/grants` ("Grant indirect permissions") request for "target subject identifier", meaning general permissions

- **Get invoice (GET `/invoices/ksef/{ksefNumber}`)**
   Added error code for 400 response: `21165` - "Invoice with given KSeF number is not yet available".

- **Invoice attachments**
  Added GET `/permissions/attachments/status` endpoint for checking consent status for issuing invoices with attachments.

- **Get sessions list**
  Extended permissions for GET `/sessions`: added `InvoiceWrite`. With `InvoiceWrite` permission, only sessions created by authenticating subject can be retrieved; with `Introspection` permission, all sessions can be retrieved.

- **Interactive session**
  - Updated error codes for POST `/sessions/online/{referenceNumber}/invoices` ("Send invoice"):
    - removed `21154` - "Interactive session ended",
    - added `21180` - "Session status does not allow operation execution".
  - Added error `21180` - "Session status does not allow operation execution" for POST `/sessions/online/{referenceNumber}/close` ("Close interactive session").

- **Batch session**
  - Added error `21180` - "Session status does not allow operation execution" for POST `/sessions/batch/{referenceNumber}/close` ("Close batch session").

- **Invoice status in session**
  Extended response for GET `/sessions/{referenceNumber}/invoices` ("Get session invoices") and GET `/sessions/{referenceNumber}/invoices/{invoiceReferenceNumber}` ("Get invoice status from session") with properties:
  - `permanentStorageDate` - date of permanent invoice storage in KSeF repository (from this moment invoice is available for download),
  - `upoDownloadUrl` - UPO download URL.

- **OpenAPI**
  - Added universal input data validation error code `21405` to all endpoints. Validator error message returned in response.
  - Added 400 response with validation returning error code `30001` ("Subject or permission already exists.") for POST `/testdata/subject` and POST `/testdata/person`.
  - Updated examples in endpoint definitions.

- **Documentation**
  - Clarified signature algorithms and examples in [QR Codes](kody-qr.md).
  - Updated code examples in C# and Java.

### Version 2.0.0 RC4

- **KSeF Certificates**
  - Added new `type` property in KSeF certificates.
  - Available certificate types:
    - `Authentication` - certificate for KSeF system authentication,
    - `Offline` - certificate limited exclusively to confirming issuer authenticity and invoice integrity in offline mode (CODE II).
  - Updated documentation for `/certificates/enrollments`, `/certificates/query`, `/certificates/retrieve` processes.

- **QR Codes**
  - Clarified that CODE II can only be generated based on `Offline` type certificate.
  - Added security warning that `Authentication` certificates cannot be used for issuing offline invoices.

- **Session status**
  - Authorization update - retrieving session, invoice and UPO information requires permission: ```InvoiceWrite```.
  - Changed *processing in progress* status code: from `300` to `150` for batch session.

- **Get invoice metadata (`/invoices/query/metadata`)**
Extended response model with fields:
  - `fileHash` - SHA256 invoice hash,
  - `hashOfCorrectedInvoice` - SHA256 hash of corrected offline invoice,
  - `thirdSubjects` - third subjects list,
  - `authorizedSubject` - authorized subject (new `InvoiceMetadataAuthorizedSubject` object containing `identifier`, `name`, `role`),
 - Added filtering by document type (`InvoiceQueryFormType`), available values: `FA`, `PEF`, `RR`.
 - `schemaType` field marked as deprecated - planned for removal in future API versions.


- **Documentation**
  - Added document describing [KSeF number](faktury/numer-ksef.md).
  - Added document describing [technical correction](offline/korekta-techniczna.md) for invoices issued in offline mode.
  - Clarified [duplicate detection](faktury/weryfikacja-faktury.md) method

- **OpenAPI**
  - Get invoice metadata list
    - Added property: `hasMore` (boolean) - indicates availability of next page of results. `totalCount` property marked as deprecated (temporarily remains in response for backward compatibility).
    - In `dateRange` filtering, `to` property (end date of range) is no longer required.
  - Granted permissions search - added `hasMore` property in response, removed `pageSize`, `pageOffset`.
  - Get authentication status - removed redundant `referenceNumber`, `isCurrent` from response.
  - Pagination unification - endpoint `/sessions/{referenceNumber}/invoices` (get session invoices) transitions to pagination based on `x-continuation-token` request header; removed `pageOffset` parameter, `pageSize` remains unchanged. First page without header; subsequent pages retrieved by passing token value returned by API. Change consistent with other resources using `x-continuation-token` (e.g., `/auth/sessions`, `/sessions/{referenceNumber}/invoices/failed`).
  - Removed `InternalId` identifier support in `targetIdentifier` field when granting indirect permissions (`/permissions/indirect/grants`). From now on only `Nip` identifier is allowed.
  - Permission granting operation status - extended list of possible status codes in response:
    - 410 - Given identifiers are inconsistent or in incorrect relationship.
    - 420 - Used credentials do not have permissions to perform this operation.
    - 430 - Identifier context does not match required role or permissions.
    - 440 - Operation not allowed for indicated identifier relationships.
    - 450 - Operation not allowed for indicated identifier or its type.
  - Added error **21418** handling - "Provided continuation token has invalid format" in all endpoints using pagination mechanism with `continuationToken` (`/auth/sessions`, `/sessions`, `/sessions/{referenceNumber}/invoices`, `/sessions/{referenceNumber}/invoices/failed`, `/tokens`).
  - Clarified invoice package download process:
    - `/invoices/exports` - initiate invoice package creation process,
    - `/invoices/async-query/{operationReferenceNumber}` - check status and receive ready package.
  - Changed model name from `InvoiceMetadataQueryRequest` to `QueryInvoicesMetadataResponse`.
  - Extended `PersonPermissionsAuthorIdentifier` type with new `System` value (System identifier). This value is used to mark permissions granted by KSeF based on submitted ZAW-FA application. Change affects endpoint: `/permissions/query/persons/grants`.

### Version 2.0.0 RC3

- **Added endpoint for getting invoice metadata list**
  - `/invoices/query` (mock) replaced by `/invoices/query/metadata` - production endpoint for getting invoice metadata
  - Update of related data models.

- **Update of mock endpoint `invoices/async-query` for initializing invoice download query**
  Updated related data models.

- **OpenAPI**
  - Supplemented endpoint specifications with required permissions (`x-required-permissions`).
  - Added `403 Forbidden` and `401 Unauthorized` responses in endpoint specifications.
  - Added `required` attribute in permission query responses.
  - Updated ```/tokens``` endpoint description
  - Removed duplicate ```enum``` definitions
  - Unified SessionInvoiceStatusResponse response model in ```/sessions/{referenceNumber}/invoices``` and ```/sessions/{referenceNumber}/invoices/{invoiceReferenceNumber}```.
  - Added 400 validation status: "Authentication failed | No assigned permissions".

- **Session status**
  - Added status ```Cancelled``` - "Session cancelled. Batch session sending time was exceeded, or no invoices were sent in interactive session."
  - Added new error codes:
    - 415 - "Unable to send invoice with attachment"
    - 440 - "Session cancelled, sending time exceeded"
    - 445 - "Verification error, no valid invoices"

- **Invoice sending status**
  - Added ```AcquisitionDate``` date - KSeF number assignment date.
  - ```ReceiveDate``` field replaced with ```InvoicingDate``` - invoice acceptance date in KSeF system.

- **Invoice sending in session**
  - Added [validation](faktury/weryfikacja-faktury.md#ograniczenia-ilościowe) of zip package size (100 MB) and package count (50) in batch session
  - Added [validation](faktury/weryfikacja-faktury.md#ograniczenia-ilościowe) of invoice count in interactive and batch session.
  - Changed "Processing in progress" status code from 300 to 150.

- **Authentication using XAdES signature**
  - ContextIdentifier fix in xsd AuthTokenRequest. Use corrected [XSD schema](https://ksef-test.mf.gov.pl/docs/v2/schemas/authv2.xsd). [XML document preparation](uwierzytelnianie.md#1-przygotowanie-dokumentu-xml-authtokenrequest)
  - Added error code `21117` - "Invalid subject identifier for indicated context type".

- **Removal of anonymous invoice download endpoint ```invoices/download```**
  Invoice download functionality without authentication has been removed; available only in KSeF web tool for invoice verification and download.

- **Test data - invoice attachment handling**
  Added new endpoints enabling testing of invoice sending with attachments.

- **KSeF Certificates - Key type and length validation in CSR**
  - Supplemented POST ```/certificates/enrollments``` endpoint description with private key type requirements in CSR (RSA, EC),
  - Added new error code 25010 in 400 response: "Invalid key type or length."

- **Public certificate format update**
  `/security/public-key-certificates` - returns certificates in DER format encoded in Base64.


### Version 2.0.0 RC2
- **New endpoints for authentication session management**
  Enable viewing and revoking active authentication sessions.
  [Authentication session management](auth/sesje.md)

- **New endpoint for getting invoice sending sessions list**\
  `/sessions` - enables retrieval of metadata for sending sessions (interactive and batch), with filtering by status, close date and session type among others.\
  [Getting sessions list](faktury/sesja-sprawdzenie-stanu-i-pobranie-upo.md#1-pobranie-listy-sesji)


- **Change in permissions listing**
  `/permissions/query/authorizations/grants` - added query type (queryType) in [entity permissions](uprawnienia.md#pobranie-listy-uprawnień-podmiotowych-do-obsługi-faktur) filtering.

- **Support for new invoice schema version FA(3)**
  When opening interactive and batch sessions, FA(3) schema can be selected.

- **Added invoiceFileName field in batch session response**\
  `/sessions/{referenceNumber}/invoices` - added invoiceFileName field containing invoice file name. Field appears only for batch sessions.
   [Getting information about sent invoices](faktury/sesja-sprawdzenie-stanu-i-pobranie-upo.md#3-pobranie-informacji-na-temat-przesłanych-faktur)
