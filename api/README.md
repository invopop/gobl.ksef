## API Client

The `api` package provides a client for interacting with KSeF. As of now, it supports:
- Authentication, using a certificate
- Uploading invoices

### Example

```go
client := api.NewClient(
    &api.ContextIdentifier{Nip: "8126178616"}, // The login session will be on the behalf of business entity specified here
    "./test/cert-20260102-131809.pfx", // Path to certificate
)

ctx := context.Background()
if err := client.Authenticate(ctx); err != nil {
  log.Fatalf("ksef auth failed: %v", err)
}

// Start invoice upload session.
session, err := client.CreateSession(ctx)
if err != nil {
  log.Fatalf("session creation failed: %v", err)
}

// In a single session, it's possible to upload multiple invoices.
// invoiceBytes should hold a serialized GOBL document.
if err := session.UploadInvoice(ctx, invoiceBytes); err != nil {
  log.Fatalf("invoice upload failed: %v", err)
}

// After uploading all invoices, finish the session. KSeF system will start processing the uploaded invoices.
if err := session.FinishUpload(ctx); err != nil {
  log.Fatalf("closing session failed: %v", err)
}

// Wait until KSeF finishes processing.
if _, err := session.PollSessionStatus(ctx); err != nil {
  log.Printf("polling failed: %v", err)
}

// If any invoices failed to process, get the details.
failed, err := session.GetFailedUploadData(ctx)
if err != nil {
  log.Printf("failed uploads lookup failed: %v", err)
}
for _, inv := range failed {
  log.Printf("failed invoice %s (ordinal %d): %+v", inv.ReferenceNumber, inv.OrdinalNumber, inv.Status)
}
```
