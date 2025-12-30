# Technical Correction of an Offline Invoice
20.08.2025

## Feature Description
Technical correction allows resubmitting an invoice issued in [offline mode](../tryby-offline.md) that was **rejected** after submission to the KSeF system due to technical errors, such as:
- schema non-compliance,
- exceeding the allowed file size,
- duplicate invoice,
- other **technical validation errors** preventing the assignment of a ```KSeF number```.


> **Note**!
1. Technical correction **does not apply** to situations related to the lack of authorization of entities appearing on the invoice (e.g., self-invoicing, validation of relations for local government units or VAT groups).
2. In this mode, **correcting the invoice content is not allowed** - technical correction applies only to technical issues preventing its acceptance in the KSeF system.
3. Technical correction can only be submitted in an [interactive session](../sesja-interaktywna.md), but it can apply to offline invoices rejected in both [interactive sessions](../sesja-interaktywna.md) and [batch sessions](../sesja-wsadowa.md).
4. It is not allowed to technically correct an offline invoice for which another valid correction has already been accepted.

## Example Workflow of Technical Correction for an Offline Invoice

1. **The seller issues an invoice in offline mode.**
   - The invoice contains two QR codes:
     - **QR Code I** - allows verification of the invoice in the KSeF system,
     - **QR Code II** - allows confirmation of the issuer's authenticity based on the [KSeF certificate](../certyfikaty-KSeF.md).

2. **The customer receives a visualization of the invoice (e.g., as a printout).**
   - After scanning **QR Code I**, the customer receives information that the invoice **has not yet been submitted to the KSeF system**.
   - After scanning **QR Code II**, the customer receives information about the KSeF certificate that confirms the issuer's authenticity.

3. **The seller submits the offline invoice to the KSeF system.**
   - The KSeF system verifies the document.
   - The invoice is **rejected** due to a technical error (e.g., invalid XSD schema compliance).

4. **The seller updates their software** and regenerates an invoice with the same content but compliant with the schema.
   - Because the XML content differs from the original version, **the SHA-256 hash of the invoice file is different**.

5. **The seller sends the corrected invoice as a technical correction.**
   - They specify in the `hashOfCorrectedInvoice` field the SHA-256 hash of the original, rejected offline invoice.
   - The `offlineMode` parameter is set to `true`.

6. **The KSeF system correctly accepts the invoice.**
   - The document receives a KSeF number.
   - The invoice is **linked to the original offline invoice** whose hash was specified in the `hashOfCorrectedInvoice` field.
   - This allows redirecting the customer from the "old" QR Code I to the corrected invoice.

7. **The customer uses QR Code I placed on the original invoice.**
   - The KSeF system informs that **the original invoice has been technically corrected**.
   - The customer receives metadata of the new, correctly processed invoice and has the ability to download it from the system.

## Sending the Correction

The correction is submitted according to the rules described in the [interactive session](../sesja-interaktywna.md) document, with additional settings:
- `offlineMode: true`,
- `hashOfCorrectedInvoice` - hash of the original invoice.

Example in C#:
[KSeF.Client.Tests.Core\E2E\OnlineSession\OnlineSessionE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/OnlineSession/OnlineSessionE2ETests.cs)
```csharp
var sendOnlineInvoiceRequest = SendInvoiceOnlineSessionRequestBuilder
    .Create()
    .WithInvoiceHash(invoiceMetadata.HashSHA, invoiceMetadata.FileSize)
    .WithEncryptedDocumentHash(
        encryptedInvoiceMetadata.HashSHA, encryptedInvoiceMetadata.FileSize)
    .WithEncryptedDocumentContent(Convert.ToBase64String(encryptedInvoice))
    .WithOfflineMode(true)
    .WithHashOfCorrectedInvoice(hashOfCorrectedInvoice)
    .Build();
```

Example in Java:
[OnlineSessionController#sendTechnicalCorrection.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/main/java/pl/akmf/ksef/sdk/api/OnlineSessionController.java#L120)
```java
SendInvoiceOnlineSessionRequest sendInvoiceOnlineSessionRequest = new SendInvoiceOnlineSessionRequestBuilder()
           .withInvoiceHash(invoiceMetadata.getHashSHA())
           .withInvoiceSize(invoiceMetadata.getFileSize())
           .withEncryptedInvoiceHash(encryptedInvoiceMetadata.getHashSHA())
           .withEncryptedInvoiceSize(encryptedInvoiceMetadata.getFileSize())
           .withEncryptedInvoiceContent(Base64.getEncoder().encodeToString(encryptedInvoice))
           .withOfflineMode(true)
           .withHashOfCorrectedInvoice(hashOfCorrectedInvoice)
        .build();
```

## Related Documents
- [Offline Modes](../tryby-offline.md)
- [QR Codes](../kody-qr.md)
