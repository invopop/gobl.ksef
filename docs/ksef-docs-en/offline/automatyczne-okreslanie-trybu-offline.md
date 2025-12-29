## Automatic Determination of Offline Mode
04.10.2025

For invoices sent as online (`offlineMode: false`), the KSeF system may assign them offline mode - based on comparing the issue date with the date of acceptance for processing.

## Algorithm Mechanism

For invoices sent with `offlineMode: false`, the system compares:
- **issue date** of the invoice (`issueDate`, e.g., `P_1` for an invoice compliant with FA(3)),
- **acceptance date** of the invoice in the KSeF system for further processing (`invoicingDate`).

Rules:
- If the calendar day from `issueDate` is earlier than the calendar day from `invoicingDate` (comparison by date, not by time), the system automatically marks the invoice as **offline**, even if it was not declared as such.
- If the `issueDate` day and the `invoicingDate` day are the same, the invoice remains **online**.

The `invoicingDate` value depends on the submission mode:
- **batch session** - `invoicingDate` is the moment the session was opened (equal to `dateCreated` returned in the session status - GET `/sessions/{referenceNumber}`),
- **interactive session** - `invoicingDate` is the moment the invoice was submitted.

This means that if, for example, an invoice was issued on 2025-10-03 (`P_1`) and submitted on 2025-10-04 at 00:00:01, despite offlineMode: false, it will be marked as an offline invoice.

## Examples
**Batch session** opened at 23:59:59 on October 3:
Even if the package is submitted after midnight, invoices will remain online - because `invoicingDate` is October 3 (the session opening date).

**Interactive session** started at 23:59:59 on October 3, and invoices were submitted after midnight:
If `P_1` = 2025-10-03, the system will mark them as offline - because the `P_1` day is earlier than the submission day.


## Related Documents
- [Offline Modes](../tryby-offline.md)
