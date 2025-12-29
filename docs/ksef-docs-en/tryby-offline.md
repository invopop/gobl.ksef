## Offline Modes
10.07.2025

## Introduction

The KSeF system offers two basic modes for issuing invoices:
* ```online``` mode - invoice issued and transmitted to the KSeF system in real time,
* ```offline``` mode - invoice issued and transmitted to KSeF at a later, legally specified date.

In offline mode, invoices are issued electronically, in accordance with the applicable FA(3) structure template. Key technical aspects:
* When sending an invoice – both in interactive and batch mode – the parameter `offlineMode: true` must be set.
* For invoices sent as online (offlineMode: false), the KSeF system may independently assign them offline mode - based on comparing the issue date with the acceptance date. Details of the mechanism: [Automatic determination of offline submission mode](offline/automatyczne-okreslanie-trybu-offline.md).
* The KSeF system accepts only the value contained in the ```P_1``` field of the e-invoice structure as the issue date.
* The invoice receipt date is the date when the KSeF number was assigned, or in case of delivery outside KSeF, the date of actual receipt.
* After issuing an invoice in offline mode, the client application should generate two [QR codes](kody-qr.md) for invoice visualization:
  * **CODE I** – enables invoice verification in the KSeF system,
  * **CODE II** – confirms the issuer's identity.
* A correcting invoice is sent only after the KSeF number has been assigned to the original document.
* If the submitted offline invoice is rejected for technical reasons, the [technical correction](/offline/korekta-techniczna.md) mechanism can be used.


### Comparison of Invoice Issuing Modes in KSeF – offline24, offline and emergency

| Mode          | Responsible party | Activation circumstances                                               | Deadline for submission to KSeF                                                         | Legal basis                                 |
| ------------- | ----------------- | ---------------------------------------------------------------------- | --------------------------------------------------------------------------------------- | ------------------------------------------- |
| **offline24** | client            | No restrictions (taxpayer's discretion)                                | by the next business day after the issue date                                           | Art. 106nda of the VAT Act (KSeF 2.0 draft) |
| **offline**   | KSeF system       | System unavailability (announced in BIP and interface software)        | by the next business day after the end of unavailability                                | Art. 106nh of the VAT Act (from 1 Feb 2026) |
| **emergency** | KSeF system       | KSeF failure (announcement in MF BIP and in interface software)        | up to 7 business days from the end of failure (with next announcement the counter resets) | Art. 106nf of the VAT Act (from 1 Feb 2026) |

### Deadline for Submitting Invoice to KSeF in Case of Subsequent Events
In offline24 and offline modes, if a KSeF failure is announced during the expected invoice submission period (announcement in MF BIP or in interface software), the submission deadline is postponed and is counted from the day the last announced failure ends, but no longer than 7 business days.

In emergency mode, if another failure announcement appears during the seven-day period for submitting the invoice, the deadline counter resets and runs from the day that subsequent failure ends.

Announcement of a total failure during any of the above modes removes the obligation to submit invoices to KSeF.

#### Example: offline24 mode with announced KSeF failure
1. 2025-07-08 (Wednesday)
    * The taxpayer generates an invoice in offline24 mode (offlineMode = true).
    * The submission deadline to KSeF is set for 2025-07-09 (next business day).
2. 2025-07-09 (Thursday)
    * The Ministry of Finance publishes an announcement about KSeF failure (BIP and API interface).
    * According to the rule: the original deadline is postponed, and the new one is counted from the day the failure ends.
3. 2025-07-12 (Saturday)
    * The failure is resolved – the system is available again.
    * A period of 7 business days begins to submit the overdue invoice.
4. 2025-07-22 (Tuesday)
    * The deadline of 7 business days from the end of failure expires.
    * The application has until this date to send the invoice to KSeF with offlineMode = true set.


### Total Failure Mode
In case of announcement of total failure (mass media: TV, radio, press, Internet):
* The invoice can be issued in paper or electronic form, without the obligation to use the FA(3) template.
* There is no obligation to submit the invoice to KSeF after the failure ends.
* Delivery to the buyer is done through any channel (in person, email, other).
* The issue date is always the actual date indicated on the invoice and the receipt date is the actual date of receiving the purchase invoice.
* Invoices from this mode are not marked with QR codes.
* A correcting invoice during an ongoing KSeF failure is issued similarly – outside KSeF, with the actual date.

## Related Documents
- [Technical correction of offline invoice](offline/korekta-techniczna.md)
- [QR Codes](kody-qr.md)
