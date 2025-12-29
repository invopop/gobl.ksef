# Limits
21.10.2025

## Introduction

In the KSeF 2.0 system, mechanisms limiting the number and size of API operations and parameters related to transmitted data have been implemented. The purpose of these limits is:
- to protect system stability at large scale of operation,
- to counteract abuse and inefficient integrations,
- to prevent abuse and potential cybersecurity threats,
- to ensure equal access conditions for all users.

Limits have been designed with the ability to flexibly adapt to the needs of specific entities requiring higher operation intensity.

## API Request Limits
The KSeF system limits the number of queries that can be sent in a short time to ensure stable system operation and equal access for all users.
More information can be found in [API Request Limits](limity-api.md).

## Context Limits

| Parameter                                                    | Default Value                       |
| ----------------------------------------------------------- | -------------------------------------- |
| Maximum invoice size without attachment                | 1 MB                                  |
| Maximum invoice size with attachment                 | 3 MB                                  |
| Maximum number of invoices in interactive/batch session | 10,000                                 |

## Authenticated Entity Limits

### Applications and Active Certificates

| Certificate Identifier            | KSeF Certificate Applications | Active KSeF Certificates |
| -------------------------------------- | ------------------------- | ------------------------ |
| NIP                                    | 300                       | 100                      |
| PESEL                                  | 6                         | 2                        |
| Certificate fingerprint | 6                         | 2                        |



## Limit Customization

The KSeF system allows individual customization of selected technical limits for:
- API limits - e.g., increasing the number of requests for a selected endpoint,
- context - e.g., increasing the maximum invoice size,
- authenticating entity - e.g., increasing active KSeF certificate limits for a natural person (PESEL).

On the **production environment**, limit increases are only possible based on a justified application supported by a real operational need.
Applications are submitted via the [contact form](https://ksef.podatki.gov.pl/formularz/), along with a detailed description of the use case.

## Checking Individual Limits
The KSeF system provides endpoints for checking the current limit values for the current context or entity:

### Get Limits for Current Context

GET [/limits/context](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Limity-i-ograniczenia/paths/~1api~1v2~1limits~1context/get)

Returns the applicable interactive and batch session limit values for the current context.

Example in C#:
[KSeF.Client.Tests.Core/E2E/Limits/LimitsE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/Limits/LimitsE2ETests.cs)
```csharp
Client.Core.Models.TestData.SessionLimitsInCurrentContextResponse limitsForContext =
    await LimitsClient.GetLimitsForCurrentContextAsync(
        accessToken,
        CancellationToken);
```
Example in Java:

[ContextLimitIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/ContextLimitIntegrationTest.java)

```java
GetContextLimitResponse response = ksefClient.getContextSessionLimit(accessToken);
```

### Get Limits for Current Entity

GET [/limits/subject](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Limity-i-ograniczenia/paths/~1api~1v2~1limits~1subject/get)

Returns the applicable certificate and certificate application limits for the current authenticated entity.

Example in C#:
[KSeF.Client.Tests.Core/E2E/Limits/LimitsE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/Limits/LimitsE2ETests.cs)
```csharp
Client.Core.Models.TestData.CertificatesLimitInCurrentSubjectResponse limitsForSubject =
        await LimitsClient.GetLimitsForCurrentSubjectAsync(
            accessToken,
            CancellationToken);
```

Example in Java:

[SubjectLimitIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/SubjectLimitIntegrationTest.java)

```java
GetSubjectLimitResponse response = ksefClient.getSubjectCertificateLimit(accessToken);
```

## Modifying Limits on Test Environment

On the **test environment**, a set of methods has been made available for changing and restoring limits to default values.
These operations are available only for authenticated entities and do not affect the production environment.

### Change Session Limits for Current Context

POST [/testdata/limits/context/session](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Limity-i-ograniczenia/paths/~1api~1v2~1testdata~1limits~1context~1session/post)

Example in C#:
[KSeF.Client.Tests.Core/E2E/Limits/LimitsE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/Limits/LimitsE2ETests.cs)

```csharp
Client.Core.Models.TestData.ChangeSessionLimitsInCurrentContextRequest newLimits =
    new()
    {
        OnlineSession = new Client.Core.Models.TestData.SessionLimits
        {
            MaxInvoices = newMaxInvoices,
            MaxInvoiceSizeInMB = newMaxInvoiceSizeInMB
            MaxInvoiceWithAttachmentSizeInMB = newMaxInvoiceWithAttachmentSizeInMB
        },

        BatchSession = new Client.Core.Models.TestData.SessionLimits
        {
            MaxInvoices = newBatchSessionMaxInvoices
            MaxInvoiceSizeInMB = newBatchSessionMaxInvoiceSizeInMB,
            MaxInvoiceWithAttachmentSizeInMB = newBatchSessionMaxInvoiceWithAttachmentSizeInMB,
        }
    };

await TestDataClient.ChangeSessionLimitsInCurrentContextAsync(
    newLimits,
    accessToken);
```

Example in Java:

[ContextLimitIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/ContextLimitIntegrationTest.java)

```java
ChangeContextLimitRequest request = new ChangeContextLimitRequest();
OnlineSessionLimit onlineSessionLimit = new OnlineSessionLimit();
onlineSessionLimit.setMaxInvoiceSizeInMB(4);
onlineSessionLimit.setMaxInvoiceWithAttachmentSizeInMB(5);
onlineSessionLimit.setMaxInvoices(6);

BatchSessionLimit batchSessionLimit = new BatchSessionLimit();
batchSessionLimit.setMaxInvoiceSizeInMB(4);
batchSessionLimit.setMaxInvoiceWithAttachmentSizeInMB(5);
batchSessionLimit.setMaxInvoices(6);

request.setOnlineSession(onlineSessionLimit);
request.setBatchSession(batchSessionLimit);

ksefClient.changeContextLimitTest(request, accessToken);
```

### Restore Session Limits for Context to Default Values

DELETE [/testdata/limits/context/session](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Limity-i-ograniczenia/paths/~1api~1v2~1testdata~1limits~1context~1session/delete)

Example in C#:
[KSeF.Client.Tests.Core/E2E/Limits/LimitsE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/Limits/LimitsE2ETests.cs)

```csharp
await TestDataClient.RestoreDefaultSessionLimitsInCurrentContextAsync(accessToken);
```

Example in Java:
[ContextLimitIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/ContextLimitIntegrationTest.java)

```java
ksefClient.resetContextLimitTest(accessToken);
```

### Change Certificate Limits for Current Entity

POST [/testdata/limits/subject/certificate](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Limity-i-ograniczenia/paths/~1api~1v2~1testdata~1limits~1subject~1certificate/post)

Example in C#:
[KSeF.Client.Tests.Core/E2E/Limits/LimitsE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/Limits/LimitsE2ETests.cs)

```csharp
Client.Core.Models.TestData.ChangeCertificatesLimitInCurrentSubjectRequest newCertificateLimitsForSubject = new()
{
    SubjectIdentifierType = Client.Core.Models.TestData.TestDataSubjectIdentifierType.Nip,
    Certificate = new Client.Core.Models.TestData.TestDataCertificate
    {
        MaxCertificates = newMaxCertificatesValue
    },
    Enrollment = new Client.Core.Models.TestData.TestDataEnrollment
    {
        MaxEnrollments = newMaxEnrollmentsValue
    }
};

await TestDataClient.ChangeCertificatesLimitInCurrentSubjectAsync(
    newCertificateLimitsForSubject,
    accessToken);
```

Example in Java:
[SubjectLimitIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/SubjectLimitIntegrationTest.java)

```java
ChangeSubjectCertificateLimitRequest request = new ChangeSubjectCertificateLimitRequest();
request.setCertificate(new CertificateLimit(15));
request.setEnrollment(new EnrollmentLimit(15));
request.setSubjectIdentifierType(ChangeSubjectCertificateLimitRequest.SubjectType.NIP);

ksefClient.changeSubjectLimitTest(request, accessToken);
```

### Restore Certificate Limits for Entity to Default Values ###

DELETE [/testdata/limits/subject/certificate](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Limity-i-ograniczenia/paths/~1api~1v2~1testdata~1limits~1subject~1certificate/delete)

Example in C#:
[KSeF.Client.Tests.Core/E2E/Limits/LimitsE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/Limits/LimitsE2ETests.cs)

```csharp
await TestDataClient.RestoreDefaultCertificatesLimitInCurrentSubjectAsync(accessToken);
```

Example in Java:
[SubjectLimitIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/SubjectLimitIntegrationTest.java)

```java
ksefClient.resetSubjectCertificateLimit(accessToken);
```

Related documents:
- [API Request Limits](limity-api.md)
- [Invoice Verification](../faktury/weryfikacja-faktury.md)
- [KSeF Certificates](../certyfikaty-KSeF.md)
