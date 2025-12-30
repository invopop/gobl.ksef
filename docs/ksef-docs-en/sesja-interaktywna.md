## Interactive Session
10.07.2025

The interactive session is used to send individual structured invoices to the KSeF API. Each invoice must be prepared in XML format according to the current schema published by the Ministry of Finance.

### Prerequisites

To use interactive submission, you must first complete the [authentication](uwierzytelnianie.md) process and have a valid access token (```accessToken```), which authorizes the use of protected KSeF API resources.

Before opening a session and sending invoices, the following is required:
* generating a symmetric key with a length of 256 bits and an initialization vector with a length of 128 bits (IV), attached as a prefix to the ciphertext,
* encrypting the document using the AES-256-CBC algorithm with PKCS#7 padding,
* encrypting the symmetric key using the RSAES-OAEP algorithm (OAEP padding with MGF1 function based on SHA-256 and SHA-256 hash), using the public key of the Ministry of Finance KSeF.

These operations can be performed using the ```CryptographyService``` component, available in the KSeF client.

Example in C#:
[KSeF.Client.Tests.Core\E2E\OnlineSession\OnlineSessionE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/OnlineSession/OnlineSessionE2ETests.cs)

```csharp
EncryptionData encryptionData = CryptographyService.GetEncryptionData();
```
Example in Java:
[OnlineSessionIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/OnlineSessionIntegrationTest.java)

```java
EncryptionData encryptionData = cryptographyService.getEncryptionData();
```

### 1. Opening a Session

Initialization of a new interactive session with the following parameters:
* invoice schema version: [FA(2)](faktury/schemy/FA/schemat_FA(2)_v1-0E.xsd), [FA(3)](faktury/schemy/FA/schemat_FA(3)_v1-0E.xsd) <br>
specifies which XSD version the system will use to validate submitted invoices.
* encrypted symmetric key<br>
the symmetric key for encrypting XML files, encrypted with the Ministry of Finance public key; it is recommended to use a newly generated key for each session.

Opening a session is a lightweight and synchronous operation – multiple interactive sessions can be maintained simultaneously within a single authentication.

POST [sessions/online](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Wysylka-interaktywna/operation/onlineSession.open)

The response returns an object containing:
 - ```referenceNumber``` – a unique identifier of the interactive session, which must be passed in all subsequent API calls.
 - ```validUntil``` – Session validity period. After this time expires, the session will be automatically closed. The lifetime of an interactive session is 12 hours from the moment of its creation.

Example in C#:
[KSeF.Client.Tests.Core\E2E\OnlineSession\OnlineSessionE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/OnlineSession/OnlineSessionE2ETests.cs)

```csharp
OpenOnlineSessionRequest openOnlineSessionRequest = OpenOnlineSessionRequestBuilder
    .Create()
    .WithFormCode(systemCode: SystemCodeHelper.GetValue(systemCode), schemaVersion: DefaultSchemaVersion, value: DefaultFormCodeValue)
    .WithEncryption(
        encryptedSymmetricKey: encryptionData.EncryptionInfo.EncryptedSymmetricKey,
        initializationVector: encryptionData.EncryptionInfo.InitializationVector)
    .Build();

OpenOnlineSessionResponse openOnlineSessionResponse = await KsefClient.OpenOnlineSessionAsync(openOnlineSessionRequest, accessToken, CancellationToken);
```

Example in Java:
[OnlineSessionIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/OnlineSessionIntegrationTest.java)

```java
OpenOnlineSessionRequest request = new OpenOnlineSessionRequestBuilder()
        .withFormCode(new FormCode(systemCode, schemaVersion, value))
        .withEncryptionInfo(encryptionData.encryptionInfo())
        .build();

OpenOnlineSessionResponse openOnlineSessionResponse = ksefClient.openOnlineSession(request, accessToken);
```

### 2. Sending an Invoice

The encrypted invoice should be sent to the endpoint:

POST [sessions/online/{referenceNumber}/invoices/](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Wysylka-interaktywna/paths/~1api~1v2~1sessions~1online~1%7BreferenceNumber%7D~1invoices/post)

The response contains the document's ```referenceNumber``` – used to identify the invoice in subsequent operations (e.g., document lists).

After the invoice is successfully submitted, asynchronous invoice verification begins ([verification details](faktury/weryfikacja-faktury.md)).

Example in C#:
[KSeF.Client.Tests.Core\E2E\OnlineSession\OnlineSessionE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/OnlineSession/OnlineSessionE2ETests.cs)

```csharp
byte[] encryptedInvoice = cryptographyService.EncryptBytesWithAES256(invoice, encryptionData.CipherKey, encryptionData.CipherIv);
FileMetadata invoiceMetadata = cryptographyService.GetMetaData(invoice);
FileMetadata encryptedInvoiceMetadata = cryptographyService.GetMetaData(encryptedInvoice);

SendInvoiceRequest sendOnlineInvoiceRequest = SendInvoiceOnlineSessionRequestBuilder
    .Create()
    .WithInvoiceHash(invoiceMetadata.HashSHA, invoiceMetadata.FileSize)
    .WithEncryptedDocumentHash(encryptedInvoiceMetadata.HashSHA, encryptedInvoiceMetadata.FileSize)
    .WithEncryptedDocumentContent(Convert.ToBase64String(encryptedInvoice))
    .Build();

SendInvoiceResponse sendInvoiceResponse = await KsefClient.SendOnlineSessionInvoiceAsync(sendOnlineInvoiceRequest, referenceNumber, accessToken);
```

Example in Java:
[OnlineSessionIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/OnlineSessionIntegrationTest.java)

```java
byte[] invoice = "";

byte[] encryptedInvoice = defaultCryptographyService.encryptBytesWithAES256(invoice,
        encryptionData.cipherKey(),
        encryptionData.cipherIv());

FileMetadata invoiceMetadata = defaultCryptographyService.getMetaData(invoice);
FileMetadata encryptedInvoiceMetadata = defaultCryptographyService.getMetaData(encryptedInvoice);

SendInvoiceOnlineSessionRequest sendInvoiceOnlineSessionRequest = new SendInvoiceOnlineSessionRequestBuilder()
        .withInvoiceHash(invoiceMetadata.getHashSHA())
        .withInvoiceSize(invoiceMetadata.getFileSize())
        .withEncryptedInvoiceHash(encryptedInvoiceMetadata.getHashSHA())
        .withEncryptedInvoiceSize(encryptedInvoiceMetadata.getFileSize())
        .withEncryptedInvoiceContent(Base64.getEncoder().encodeToString(encryptedInvoice))
        .build();

SendInvoiceResponse sendInvoiceResponse = ksefClient.onlineSessionSendInvoice(sessionReferenceNumber, sendInvoiceOnlineSessionRequest, accessToken);

```

### 3. Closing the Session
After sending all invoices, the session should be closed, which initiates asynchronous generation of the consolidated UPO.

POST [/sessions/online/\{referenceNumber\}/close](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Wysylka-interaktywna/paths/~1api~1v2~1sessions~1online~1%7BreferenceNumber%7D~1close/post)

The consolidated UPO will be available after checking the session status.

Example in C#:
[KSeF.Client.Tests.Core\E2E\OnlineSession\OnlineSessionE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/OnlineSession/OnlineSessionE2ETests.cs)

```csharp
await KsefClient.CloseOnlineSessionAsync(referenceNumber, accessToken, CancellationToken);
```

Example in Java:
[OnlineSessionIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/OnlineSessionIntegrationTest.java)

```java
ksefClient.closeOnlineSession(sessionReferenceNumber, accessToken);
```

Related documents:
- [Checking Status and Downloading UPO](faktury/sesja-sprawdzenie-stanu-i-pobranie-upo.md)
- [Invoice Verification](faktury/weryfikacja-faktury.md)
- [KSeF Number – Structure and Validation](faktury/numer-ksef.md)
