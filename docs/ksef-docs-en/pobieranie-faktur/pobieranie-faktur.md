## Invoice retrieval
### Retrieving invoices by KSeF number
21.08.2025

Returns the invoice with the specified KSeF number.

GET [/invoices/ksef/\{ksefReferenceNumber\}](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Pobieranie-faktur/paths/~1api~1v2~1invoices~1ksef~1%7BksefNumber%7D/get)

Example in C#:
[KSeF.Client.Tests.Core\E2E\Invoice\InvoiceE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/Invoice/InvoiceE2ETests.cs)

```csharp
string invoice = await ksefClient.GetInvoiceAsync(ksefReferenceNumber, accessToken, cancellationToken);
```

Example in Java:

[OnlineSessionIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/OnlineSessionIntegrationTest.java)

```java
byte[] invoice = ksefClient.getInvoice(ksefNumber, accessToken);
```

### Retrieving the list of invoice metadata
Returns a list of invoice metadata matching the specified search criteria.

**The _metadata.json file in the export package**
The export package contains a `_metadata.json` file containing an array of `InvoiceMetadata` objects (the model returned by POST `/invoices/query/metadata` - "Retrieving invoice metadata").

POST [/invoices/query/metadata](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Pobieranie-faktur/paths/~1api~1v2~1invoices~1query~1metadata/post)

Example in C#:
[KSeF.Client.Tests.Core\E2E\Invoice\InvoiceE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/Invoice/InvoiceE2ETests.cs)

```csharp
InvoiceQueryFilters invoiceQueryFilters = new InvoiceQueryFilters
{
    SubjectType = InvoiceSubjectType.Subject1,
    DateRange = new DateRange
    {
        From = DateTime.UtcNow.AddDays(-30),
        To = DateTime.UtcNow,
        DateType = DateType.Issue
    }
};

PagedInvoiceResponse pagedInvoicesResponse = await ksefClient.QueryInvoiceMetadataAsync(
    invoiceQueryFilters,
    accessToken,
    pageOffset: 0,
    pageSize: 10,
    cancellationToken);
```

Example in Java:

[QueryInvoiceIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/QueryInvoiceIntegrationTest.java)

```java
InvoiceQueryFilters request = new InvoiceQueryFiltersBuilder()
        .withSubjectType(InvoiceQuerySubjectType.SUBJECT1)
        .withDateRange(
                new InvoiceQueryDateRange(InvoiceQueryDateType.INVOICING, OffsetDateTime.now().minusYears(1),
                        OffsetDateTime.now()))
        .build();

QueryInvoiceMetadataResponse response = ksefClient.queryInvoiceMetadata(pageOffset, pageSize, SortOrder.ASC, request, accessToken);


```

### Initialize asynchronous invoice retrieval query

Starts an asynchronous invoice search process in the KSeF system based on the provided filters. It is required to provide encryption information in the `encryption` field, which is used to encrypt the generated invoice packages.

POST [/invoices/exports](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Pobieranie-faktur/paths/~1api~1v2~1invoices~1exports/post)

Example in C#:
[KSeF.Client.Tests.Core\E2E\Invoice\InvoiceE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/Invoice/InvoiceE2ETests.cs)

```csharp
EncryptionData encryptionData = CryptographyService.GetEncryptionData();

InvoiceQueryFilters query = new InvoiceQueryFilters
{
    DateRange = new DateRange
    {
        From = DateTime.Now.AddDays(-1),
        To = DateTime.Now.AddDays(1),
        DateType = DateType.Invoicing
    },
    SubjectType = InvoiceSubjectType.Subject1
};

InvoiceExportRequest invoiceExportRequest = new InvoiceExportRequest
{
    Encryption = encryptionData.EncryptionInfo,
    Filters = query
};

OperationResponse invoicesForSellerResponse = await KsefClient.ExportInvoicesAsync(
    invoiceExportRequest,
    _accessToken,
    CancellationToken);
```

Available values for `DateType` and `SubjectType` are described [here](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Pobieranie-faktur/paths/~1api~1v2~1invoices~1query~1metadata/post).

Invoices in the package are sorted in ascending order by the date type specified in `DateRange` during export initialization.

Example in Java:

[QueryInvoiceIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/QueryInvoiceIntegrationTest.java)

```java
InvoiceExportFilters filters = new InvoicesAsyncQueryFiltersBuilder()
        .withSubjectType(InvoiceQuerySubjectType.SUBJECT1)
        .withDateRange(
                new InvoiceQueryDateRange(InvoiceQueryDateType.INVOICING, OffsetDateTime.now().minusDays(10), OffsetDateTime.now().plusDays(10)))
        .build();

InvoiceExportRequest request = new InvoiceExportRequest(
        new EncryptionInfo(encryptionData.encryptionInfo().getEncryptedSymmetricKey(),
                encryptionData.encryptionInfo().getInitializationVector()), filters);

InitAsyncInvoicesQueryResponse response = ksefClient.initAsyncQueryInvoice(request, accessToken);

```

### Check status of asynchronous invoice retrieval query

Retrieves the status of a previously initialized asynchronous query based on the operation identifier. It allows tracking the processing progress of the query and downloading ready packages with results, if already available.

GET [/invoices/exports/{operationReferenceNumber}](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Pobieranie-faktur/paths/~1api~1v2~1invoices~1exports~1%7BoperationReferenceNumber%7D/get)

Example in C#:
[KSeF.Client.Tests.Core\E2E\Invoice\InvoiceE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/Invoice/InvoiceE2ETests.cs)

```csharp
InvoiceExportStatusResponse exportStatus = await KsefClient.GetInvoiceExportStatusAsync(
    exportInvoicesResponse.OperationReferenceNumber,
    accessToken);
```

Example in Java:

[QueryInvoiceIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/QueryInvoiceIntegrationTest.java)

```java
InvoiceExportStatus response = ksefClient.checkStatusAsyncQueryInvoice(operationReferenceNumber, accessToken);

```

## Related documents

- [Incremental invoice retrieval](przyrostowe-pobieranie-faktur.md)
