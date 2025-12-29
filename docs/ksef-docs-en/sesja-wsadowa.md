## Batch Session
10.07.2025

Batch submission allows you to send multiple invoices at once in a single ZIP file, instead of sending each invoice separately.

This solution speeds up and simplifies the processing of large numbers of documents — especially for companies that generate many invoices daily.

Each invoice must be prepared in XML format according to the current schema published by the Ministry of Finance:
* The ZIP package should be divided into parts no larger than 100 MB (before encryption), which are encrypted and sent separately.
* Information about each part of the package must be provided in the ```fileParts``` object.


### Prerequisites
To use batch submission, you must first complete the [authentication](auth/authentication.md) process and have a valid access token (```accessToken```), which authorizes access to protected KSeF API resources.

**Recommendation (status correlation by `invoiceHash`)**
Before creating a package for batch submission, it is recommended to calculate the SHA-256 hash for each invoice XML file (original, before encryption) and save a local mapping. This enables unambiguous correlation of processing statuses from KSeF with local source documents (XML) prepared for submission.

Before opening a session and sending invoices, the following is required:
* generating a symmetric key with a length of 256 bits and an initialization vector with a length of 128 bits (IV), appended as a prefix to the ciphertext,
* encrypting the document with the AES-256-CBC algorithm with PKCS#7 padding,
* encrypting the symmetric key with the RSAES-OAEP algorithm (OAEP padding with MGF1 function based on SHA-256 and SHA-256 hash), using the Ministry of Finance KSeF public key.

These operations can be performed using the ```CryptographyService``` component, available in the official KSeF client. This library provides ready-made methods for generating and encrypting keys, in accordance with system requirements.

Example in C#:
[KSeF.Client.Tests.Core\E2E\BatchSession\BatchSessionE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/BatchSession/BatchSessionE2ETests.cs)

```csharp
EncryptionData encryptionData = cryptographyService.GetEncryptionData();
```
Example in Java:
[BatchIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/BatchIntegrationTest.java)

```java
EncryptionData encryptionData = cryptographyService.getEncryptionData();
```

The generated data is used to encrypt invoices.

### 1. Preparing the ZIP Package
You need to create a ZIP package containing all invoices that will be sent within a single session.

Example in C#:
[KSeF.Client.Tests.Core\E2E\BatchSession\BatchSessionE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/BatchSession/BatchSessionE2ETests.cs)

```csharp
(byte[] zipBytes, Client.Core.Models.Sessions.FileMetadata zipMeta) =
    BatchUtils.BuildZip(invoices, cryptographyService);

//BatchUtils.BuildZip
public static (byte[] ZipBytes, FileMetadata Meta) BuildZip(
    IEnumerable<(string FileName, byte[] Content)> files,
    ICryptographyService crypto)
{
    using MemoryStream zipStream = new MemoryStream();
    using ZipArchive archive = new ZipArchive(zipStream, ZipArchiveMode.Create, leaveOpen: true);

    foreach ((string fileName, byte[] content) in files)
    {
        ZipArchiveEntry entry = archive.CreateEntry(fileName, CompressionLevel.Optimal);
        using Stream entryStream = entry.Open();
        entryStream.Write(content);
    }

    archive.Dispose();

    byte[] zipBytes = zipStream.ToArray();
    List<byte[]> meta = crypto.GetMetaData(zipBytes);

    return (zipBytes, meta);
}
```

Example in Java:
[BatchIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/BatchIntegrationTest.java)

```java
byte[] zipBytes = FilesUtil.createZip(invoicesInMemory);

// get ZIP metadata (before crypto)
FileMetadata zipMetadata = defaultCryptographyService.getMetaData(zipBytes);
```

### 2. Binary Splitting of the ZIP Package into Parts

Due to file size limitations for uploads, the ZIP package should be split binary into smaller parts, which will be sent separately. Each part should have a unique name and ordinal number.

Example in C#:
[KSeF.Client.Tests.Core\E2E\BatchSession\BatchSessionE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/BatchSession/BatchSessionE2ETests.cs)

```csharp

 // Get ZIP metadata (before encryption)
FileMetadata zipMetadata = cryptographyService.GetMetaData(zipBytes);
int maxPartSize = 100 * 1000 * 1000; // 100 MB
int partCount = (int)Math.Ceiling((double)zipBytes.Length / maxPartSize);
int partSize = (int)Math.Ceiling((double)zipBytes.Length / partCount);
List<byte[]> zipParts = new List<byte[]>();
for (int i = 0; i < partCount; i++)
{
    int start = i * partSize;
    int size = Math.Min(partSize, zipBytes.Length - start);
    if (size <= 0) break;
    byte[] part = new byte[size];
    Array.Copy(zipBytes, start, part, 0, size);
    zipParts.Add(part);
}

```

Example in Java:
[BatchIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/BatchIntegrationTest.java)

```java
List<byte[]> zipParts = FilesUtil.splitZip(partsCount, zipBytes);
```

### 3. Encrypting Package Parts
Each part must be encrypted with a newly generated AES-256-CBC key with PKCS#7 padding.

Example in C#:
[KSeF.Client.Tests.Core\E2E\BatchSession\BatchSessionStreamE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/BatchSession/BatchSessionStreamE2ETests.cs)
```csharp
List<BatchPartStreamSendingInfo> encryptedParts = new(rawParts.Count);
for (int i = 0; i < rawParts.Count; i++)
{
    using MemoryStream partInput = new(rawParts[i], writable: false);
    MemoryStream encryptedOutput = new();
    await cryptographyService.EncryptStreamWithAES256Async(partInput, encryptedOutput, encryptionData.CipherKey, encryptionData.CipherIv, CancellationToken).ConfigureAwait(false);

    if (encryptedOutput.CanSeek)
    {
        encryptedOutput.Position = 0;
    }

    FileMetadata partMeta = await cryptographyService.GetMetaDataAsync(encryptedOutput, CancellationToken).ConfigureAwait(false);
    if (encryptedOutput.CanSeek)
    {
        encryptedOutput.Position = 0; // reset after reading for metadata
    }

    encryptedParts.Add(new BatchPartStreamSendingInfo
    {
        DataStream = encryptedOutput,
        OrdinalNumber = i + 1,
        Metadata = partMeta
    });
}
```

Example in Java:
[BatchIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/BatchIntegrationTest.java)

```java
 List<BatchPartSendingInfo> encryptedZipParts = new ArrayList<>();
 for (int i = 0; i < zipParts.size(); i++) {
     byte[] encryptedZipPart = defaultCryptographyService.encryptBytesWithAES256(
             zipParts.get(i),
             cipherKey,
             cipherIv
     );
     FileMetadata zipPartMetadata = defaultCryptographyService.getMetaData(encryptedZipPart);
     encryptedZipParts.add(new BatchPartSendingInfo(encryptedZipPart, zipPartMetadata, (i + 1)));
 }

```

### 4. Opening a Batch Session

Initialize a new batch session by providing:
* invoice schema version: [FA(2)](faktury/schemy/FA/schemat_FA(2)_v1-0E.xsd), [FA(3)](faktury/schemy/FA/schemat_FA(3)_v1-0E.xsd) <br>
specifies which XSD version the system will use to validate submitted invoices.
* encrypted symmetric key<br>
the symmetric key used to encrypt XML files, encrypted with the Ministry of Finance public key; it is recommended to use a newly generated key for each session.
* ZIP package and parts metadata: file name, hash, size, and list of parts (if the package is split)

POST [/sessions/batch](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Wysylka-wsadowa/operation/batch.open)

In response to opening the session, you will receive an object containing `referenceNumber`, which will be used in subsequent steps to identify the batch session.

Example in C#:
[KSeF.Client.Tests.Core\E2E\BatchSession\BatchSessionE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/BatchSession/BatchSessionE2ETests.cs)

```csharp
Client.Core.Models.Sessions.BatchSession.OpenBatchSessionRequest openBatchRequest =
    BatchUtils.BuildOpenBatchRequest(zipMeta, encryptionData, encryptedParts, systemCode);

Client.Core.Models.Sessions.BatchSession.OpenBatchSessionResponse openBatchSessionResponse =
    await BatchUtils.OpenBatchAsync(KsefClient, openBatchRequest, accessToken).ConfigureAwait(false);

//BatchUtils.BuildOpenBatchRequest
public static OpenBatchSessionRequest BuildOpenBatchRequest(
    FileMetadata zipMeta,
    EncryptionData encryption,
    IEnumerable<BatchPartSendingInfo> encryptedParts,
    SystemCode systemCode = DefaultSystemCode,
    string schemaVersion = DefaultSchemaVersion,
    string value = DefaultValue)
{
    IOpenBatchSessionRequestBuilderBatchFile builder = OpenBatchSessionRequestBuilder
        .Create()
        .WithFormCode(systemCode: SystemCodeHelper.GetValue(systemCode), schemaVersion: schemaVersion, value: value)
        .WithBatchFile(fileSize: zipMeta.FileSize, fileHash: zipMeta.HashSHA);

    foreach (BatchPartSendingInfo p in encryptedParts)
    {
        builder = builder.AddBatchFilePart(
            ordinalNumber: p.OrdinalNumber,
            fileName: $"part_{p.OrdinalNumber}.zip.aes",
            fileSize: p.Metadata.FileSize,
            fileHash: p.Metadata.HashSHA);
    }

    return builder
        .EndBatchFile()
        .WithEncryption(
            encryptedSymmetricKey: encryption.EncryptionInfo.EncryptedSymmetricKey,
            initializationVector: encryption.EncryptionInfo.InitializationVector)
        .Build();
}

//BatchUtils.OpenBatchAsync
public static async Task<OpenBatchSessionResponse> OpenBatchAsync(
    IKSeFClient client,
    OpenBatchSessionRequest openReq,
    string accessToken)
    => await client.OpenBatchSessionAsync(openReq, accessToken).ConfigureAwait(false);
```

Example in Java:
[BatchIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/BatchIntegrationTest.java)

```java
OpenBatchSessionRequestBuilder builder = OpenBatchSessionRequestBuilder.create()
        .withFormCode(SystemCode.FA_2, SchemaVersion.VERSION_1_0E, SessionValue.FA)
        .withOfflineMode(false)
        .withBatchFile(zipMetadata.getFileSize(), zipMetadata.getHashSHA());

for (int i = 0; i < encryptedZipParts.size(); i++) {
        BatchPartSendingInfo part = encryptedZipParts.get(i);
        builder = builder.addBatchFilePart(i + 1, "faktura_part" + (i + 1) + ".zip.aes",part.getMetadata().getFileSize(), part.getMetadata().getHashSHA());
}

OpenBatchSessionRequest request = builder.endBatchFile()
        .withEncryption(
                        encryptionData.encryptionInfo().getEncryptedSymmetricKey(),
                        encryptionData.encryptionInfo().getInitializationVector()
                )
        .build();

OpenBatchSessionResponse response = ksefClient.openBatchSession(request, accessToken);
```

The method returns a list of package parts; for each part it provides the upload address (URL), the required HTTP method, and a complete set of headers that must be sent along with that part.

### 5. Sending Declared Package Parts

Using the data returned when opening the session in `partUploadRequests`, i.e., the unique URL with access key, HTTP method (method), and required headers (headers), you must send each declared package part (`fileParts`) to the specified address, using exactly those values for the given part. The link between the declaration and upload instruction is the `ordinalNumber` property.

The request body should contain the bytes of the corresponding file part (without JSON wrapping).

> Note: do not add the access token (`accessToken`) to the headers.

Each part is sent as a separate HTTP request. Returned response codes:
* `201` - file successfully received,
* `400` - invalid data,
* `401` - invalid authentication,
* `403` - no write permission (e.g., write time has expired).

Example in C#:
[KSeF.Client.Tests.Core\E2E\BatchSession\BatchSessionE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/BatchSession/BatchSessionE2ETests.cs)

```csharp
await KsefClient.SendBatchPartsAsync(openBatchSessionResponse, encryptedParts);
```

Example in Java:
[BatchIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/BatchIntegrationTest.java)

```java
ksefClient.sendBatchParts(response, encryptedZipParts);
```

**Time Limit for Uploading Parts in a Batch Session**
File uploads in a batch session are time-limited. This time depends solely on the number of declared parts and is 20 minutes per part. Each additional part proportionally increases the time limit **for each part** in the package.

Total time for uploading each part = number of parts x 20 minutes.
Example: A package contains 2 parts - each part has 40 minutes for upload.

The size of the part does not affect the time limit - the only criterion is the number of parts declared when opening the session.

Authorization is verified at the beginning of each HTTP request. If the address is valid at the time the request is received, the upload operation is completed in full. Expiration during the upload does not interrupt an operation that has already started.

### 6. Closing the Batch Session
After all package parts have been sent, you must close the batch session, which asynchronously initiates the processing of the invoice package ([verification details](faktury/weryfikacja-faktury.md)), and generates a collective UPO.

POST [/sessions/batch/\{referenceNumber\}/close](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Wysylka-wsadowa/paths/~1api~1v2~1sessions~1batch~1%7BreferenceNumber%7D~1close/post)}]

Example in C#:
[KSeF.Client.Tests.Core\E2E\BatchSession\BatchSessionStreamE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/BatchSession/BatchSessionStreamE2ETests.cs)
```csharp
await KsefClient.CloseBatchSessionAsync(referenceNumber, accessToken);
```
Example in Java:
[BatchIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/BatchIntegrationTest.java)

```java
ksefClient.closeBatchSession(referenceNumber, accessToken);
```

See also
- [Checking Status and Downloading UPO](faktury/sesja-sprawdzenie-stanu-i-pobranie-upo.md)
- [Invoice Verification](faktury/weryfikacja-faktury.md)
- [KSeF Number - Structure and Validation](faktury/numer-ksef.md)
