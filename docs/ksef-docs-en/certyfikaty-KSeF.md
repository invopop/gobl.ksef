## KSeF Certificates
31.08.2025

### Introduction
A KSeF certificate is a digital identity credential issued by the KSeF system upon user request.

A request for a KSeF certificate can only be submitted for data that is contained in the certificate used for [authentication](uwierzytelnianie.md). Based on this data, the endpoint [/certificates/enrollments/data](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Certyfikaty/paths/~1api~1v2~1certificates~1enrollments~1data/get)
 returns the identification data that must be used in the certificate request.

>The system does not allow requesting a certificate on behalf of another entity.

Two types of certificates are available – each certificate can have **only one type** (`Authentication` or `Offline`). It is not possible to issue a certificate combining both functions.

| Type             | Description |
| ---------------- | ---- |
| `Authentication` | Certificate intended for authentication in the KSeF system.<br/>**keyUsage:** Digital Signature (80) |
| `Offline`        | Certificate intended exclusively for issuing invoices in offline mode. Used to confirm the authenticity of the issuer and the integrity of the invoice through [QR code II](kody-qr.md). Does not enable authentication.<br/>**keyUsage:** Non-Repudiation (40) |

#### Certificate Acquisition Process
The certificate application process consists of several stages:
1. Checking available limits,
2. Retrieving data for the certificate request,
3. Submitting the request,
4. Downloading the issued certificate,


### 1. Checking Limits

Before an API client submits a request for a new certificate, it is recommended to verify the certificate limits.

The API provides information about:
* the maximum number of certificates that can be held,
* the number of currently active certificates,
* the possibility of submitting another request.

GET [/certificates/limits](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Certyfikaty/paths/~1api~1v2~1certificates~1limits/get)

Example in C#:
[KSeF.Client.Tests.Core\E2E\Certificates\CertificatesE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/Certificates/CertificatesE2ETests.cs)
```csharp
CertificateLimitResponse certificateLimitResponse = await KsefClient
    .GetCertificateLimitsAsync(accessToken, CancellationToken);
```

Example in Java:
[CertificateIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/CertificateIntegrationTest.java)

```java
CertificateLimitsResponse response = ksefClient.getCertificateLimits(accessToken);
```

### 2. Retrieving Data for the Certificate Request

To begin the process of applying for a KSeF certificate, you need to retrieve a set of identification data that the system will return in response to calling the endpoint
GET [/certificates/enrollments/data](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Certyfikaty/paths/~1api~1v2~1certificates~1enrollments~1data/get).

This data is read from the certificate used for authentication, which can be:
- a qualified certificate of a natural person – containing a PESEL or NIP number,
- a qualified organization certificate (so-called company seal) – containing a NIP number,
- Trusted Profile (ePUAP) – used by natural persons, contains a PESEL number,
- KSeF internal certificate – issued by the KSeF system, it is not a qualified certificate but is accepted in the authentication process.

Based on this, the system returns a complete set of DN attributes (X.500 Distinguished Name) that must be used when building the certificate signing request (CSR). Modification of this data will result in rejection of the request.

**Note**: Retrieving certificate data is only possible after authentication using a signature (XAdES). Authentication using a KSeF system token does not allow submitting a certificate request.


Example in C#:
[KSeF.Client.Tests.Core\E2E\Certificates\CertificatesE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/Certificates/CertificatesE2ETests.cs)
```csharp
CertificateEnrollmentsInfoResponse certificateEnrollmentsInfoResponse =
    await KsefClient.GetCertificateEnrollmentDataAsync(accessToken, CancellationToken).ConfigureAwait(false);
```

Example in Java:
[CertificateIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/CertificateIntegrationTest.java)

```java
CertificateEnrollmentsInfoResponse response = ksefClient.getCertificateEnrollmentInfo(accessToken);
```

Here is the complete list of fields that may be returned, presented in a table containing the OID:

| OID      | Name                  | Description                            | Natural person | Company seal |
|----------|-----------------------|----------------------------------------|----------------|-----------------|
| 2.5.4.3  | commonName            | Common name                            | ✔️             | ✔️              |
| 2.5.4.4  | surname               | Surname                                | ✔️             | ❌              |
| 2.5.4.5  | serialNumber          | Serial number (e.g., PESEL, NIP)       | ✔️             | ❌              |
| 2.5.4.6  | countryName           | Country code (e.g., PL)                | ✔️             | ✔️              |
| 2.5.4.10 | organizationName      | Organization name / company            | ❌             | ✔️              |
| 2.5.4.42 | givenName             | First name or names                    | ✔️             | ❌              |
| 2.5.4.45 | uniqueIdentifier      | Unique identifier (optional)           | ✔️             | ✔️              |
| 2.5.4.97 | organizationIdentifier| Organization identifier (e.g., NIP)    | ❌             | ✔️              |


The `givenName` attribute may appear multiple times and is returned as a list of values.

### 3. Preparing CSR (Certificate Signing Request)
To submit a request for a KSeF certificate, you need to prepare a so-called certificate signing request (CSR) in the PKCS#10 standard, in DER format, encoded in Base64. The CSR contains:
* information identifying the entity (DN – Distinguished Name),
* the public key that will be associated with the certificate.

Requirements for the private key used to sign the CSR:
* Allowed types:
  * RSA (OID: 1.2.840.113549.1.1.1), key length: 2048 bits,
  * EC (elliptic curve keys, OID: 1.2.840.10045.2.1), NIST P-256 curve (secp256r1).
* EC keys are recommended.

* Allowed signature algorithms:
  * RSA PKCS#1 v1.5,
  * RSA PSS,
  * ECDSA (signature format compliant with RFC 3279).

* Allowed hash functions used for CSR signing:
  * SHA1,
  * SHA256,
  * SHA384,
  * SHA512.

All identification data (X.509 attributes) should match the values returned by the system in the previous step (/certificates/enrollments/data). Modifying this data will result in rejection of the request.

Example in C# (using ```ICryptographyService```):
[KSeF.Client.Tests.Core\E2E\Certificates\CertificatesE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/Certificates/CertificatesE2ETests.cs)

```csharp
var (csr, key) = CryptographyService.GenerateCsrWithRSA(TestFixture.EnrollmentInfo);
```


Example in Java:
[CertificateIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/CertificateIntegrationTest.java)

```java
CsrResult csr = defaultCryptographyService.generateCsrWithRsa(enrollmentInfo);
```

* ```csrBase64Encoded``` – contains the CSR request encoded in Base64 format, ready to be sent to KSeF
* ```privateKeyBase64Encoded``` – contains the private key associated with the generated CSR, encoded in Base64. This key will be needed for signing operations using the certificate.

**Note**: The private key should be stored securely and in accordance with the security policy of the given organization.

### 4. Submitting the Certificate Request
After preparing the certificate signing request (CSR), it should be sent to the KSeF system via the call

POST [/certificates/enrollments](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Certyfikaty/paths/~1api~1v2~1certificates~1enrollments/post)

In the submitted request, you need to provide:
* **certificate name** – visible later in the certificate metadata, facilitating identification,
* **certificate type** – `Authentication` or `Offline`,
* **CSR** in PKCS#10 format (DER), encoded as a Base64 string,
* (optionally) **validFrom** – the validity start date. If not specified, the certificate will be valid from the moment of its issuance.

Make sure the CSR contains exactly the same data that was returned by the /certificates/enrollments/data endpoint.

Example in C#:
[KSeF.Client.Tests.Core\E2E\Certificates\CertificatesE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/Certificates/CertificatesE2ETests.cs)

```csharp
SendCertificateEnrollmentRequest sendCertificateEnrollmentRequest = SendCertificateEnrollmentRequestBuilder
    .Create()
    .WithCertificateName(TestCertificateName)
    .WithCertificateType(CertificateType.Authentication)
    .WithCsr(csr)
    .WithValidFrom(DateTimeOffset.UtcNow.AddDays(CertificateValidityDays))
    .Build();

CertificateEnrollmentResponse certificateEnrollmentResponse = await KsefClient
    .SendCertificateEnrollmentAsync(sendCertificateEnrollmentRequest, accessToken, CancellationToken);
```

Example in Java:
[CertificateIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/CertificateIntegrationTest.java)

```java
SendCertificateEnrollmentRequest request = new SendCertificateEnrollmentRequestBuilder()
        .withValidFrom(OffsetDateTime.now().toString())
        .withCsr(csr.csr())
        .withCertificateName("certificate")
        .withCertificateType(CertificateType.AUTHENTICATION)
        .build();

CertificateEnrollmentResponse response = ksefClient.sendCertificateEnrollment(request, accessToken);
```

In the response, you will receive a ```referenceNumber```, which allows you to monitor the status of the request and later download the issued certificate.

### 5. Checking Request Status

The certificate issuance process is asynchronous. This means that the system does not return the certificate immediately after submitting the request, but allows it to be downloaded later after processing is complete.
The request status should be checked periodically using the reference number (```referenceNumber```) that was returned in the response to the request submission (/certificates/enrollments).

GET [/certificates/enrollments/\{referenceNumber\}](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Certyfikaty/paths/~1api~1v2~1certificates~1enrollments~1%7BreferenceNumber%7D/get)

If the certificate request is rejected, the response will contain error information.

Example in C#:
[KSeF.Client.Tests.Core\E2E\Certificates\CertificatesE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/Certificates/CertificatesE2ETests.cs)

```csharp
CertificateEnrollmentStatusResponse certificateEnrollmentStatusResponse = await KsefClient
    .GetCertificateEnrollmentStatusAsync(TestFixture.EnrollmentReference, accessToken, CancellationToken);
```

Example in Java:
[CertificateIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/CertificateIntegrationTest.java)

```java
CertificateEnrollmentStatusResponse response = ksefClient.getCertificateEnrollmentStatus(referenceNumber, accessToken);

```

After obtaining the certificate serial number (```certificateSerialNumber```), it is possible to download its content and metadata in the subsequent steps of the process.

### 6. Retrieving Certificate List

The KSeF system allows downloading the content of previously issued internal certificates based on a list of serial numbers. Each certificate is returned in DER format, encoded as a Base64 string.

POST [/certificates/retrieve](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Certyfikaty/paths/~1api~1v2~1certificates~1retrieve/post)

Example in C#:
[KSeF.Client.Tests.Core\E2E\Certificates\CertificatesE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/Certificates/CertificatesE2ETests.cs)

```csharp
CertificateListRequest certificateListRequest = new CertificateListRequest { CertificateSerialNumbers = TestFixture.SerialNumbers };

CertificateListResponse certificateListResponse = await KsefClient
    .GetCertificateListAsync(certificateListRequest, accessToken, CancellationToken);
```

Example in Java:
[CertificateIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/CertificateIntegrationTest.java)

```java
CertificateListResponse certificateResponse = ksefClient.getCertificateList(new CertificateListRequest(List.of(certificateSerialNumber)), accessToken);
```

Each element of the response contains:

| Field                     | Description    |
|---------------------------|------------------------|
| `certificateSerialNumber` | Certificate serial number          |
| `certificateName` | Certificate name assigned during registration          |
| `certificate` | Certificate content encoded in Base64 (DER format)          |
| `certificateType` | Certificate type (`Authentication`, `Offline`).          |

### 7. Retrieving Certificate Metadata List

It is possible to retrieve a list of internal certificates submitted by a given entity. This data includes both active and historical certificates, along with their status, validity range, and identifiers.

POST [/certificates/query](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Certyfikaty/paths/~1api~1v2~1certificates~1query/post)

Filtering parameters (optional):
* `status` - certificate status (`Active`, `Blocked`, `Revoked`, `Expired`)
* `expiresAfter` - certificate expiration date (optional)
* `name` - certificate name (optional)
* `type` - certificate type (`Authentication`, `Offline`) (optional)
* `certificateSerialNumber` - certificate serial number (optional)
* `pageSize` - number of elements per page (default 10)
* `pageOffset` - page number of results (default 0)

Example in C#:
[KSeF.Client.Tests.Core\E2E\Certificates\CertificateMetadataListE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core\E2E\Certificates/CertificateMetadataListE2ETests.cs)

```csharp
var request = GetCertificateMetadataListRequestBuilder
    .Create()
    .WithCertificateSerialNumber(serialNumber)
    .WithName(name)
    .Build();
    CertificateMetadataListResponse certificateMetadataListResponse = await KsefClient
            .GetCertificateMetadataListAsync(accessToken, requestPayload, pageSize, pageOffset, CancellationToken);
```
Example in Java:
[CertificateIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/CertificateIntegrationTest.java)

```java
QueryCertificatesRequest request = new CertificateMetadataListRequestBuilder().build();

CertificateMetadataListResponse response = ksefClient.getCertificateMetadataList(request, pageSize, pageOffset, accessToken);


```

In the response, you will receive certificate metadata.



### 8. Revoking Certificates

A KSeF certificate can only be revoked by its owner in case of private key compromise, end of use, or organizational change. After revocation, the certificate cannot be used for further authentication or operations in the KSeF system.
Revocation is performed based on the certificate serial number (```certificateSerialNumber```) and an optional revocation reason.

POST [/certificates/\{certificateSerialNumber\}/revoke](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Certyfikaty/paths/~1api~1v2~1certificates~1%7BcertificateSerialNumber%7D~1revoke/post)

Example in C#:
[KSeF.Client.Tests.Core\E2E\Certificates\CertificatesE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/Certificates/CertificatesE2ETests.cs)
```csharp
CertificateRevokeRequest certificateRevokeRequest = RevokeCertificateRequestBuilder
        .Create()
        .WithRevocationReason(CertificateRevocationReason.KeyCompromise)
        .Build();

await ksefClient.RevokeCertificateAsync(request, certificateSerialNumber, accessToken, cancellationToken)
     .ConfigureAwait(false);
```

Example in Java:
[CertificateIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/CertificateIntegrationTest.java)

```java
CertificateRevokeRequest request = new CertificateRevokeRequestBuilder()
        .withRevocationReason(CertificateRevocationReason.KEYCOMPROMISE)
        .build();

ksefClient.revokeCertificate(request, serialNumber, accessToken);
```

After revocation, the certificate cannot be reused. If further use is needed, a new certificate must be requested.
