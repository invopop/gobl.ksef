## KSeF API 2.0 Environments
15.12.2025

Below is a summary of information about public environments.

| Abbreviation | Environment                       | Description                                                                 | API Documentation                         | Allowed Formats                         |
|-------|----------------------------------|-----------------------------------|----------------------------------|----------------------------------------------|
| **TEST**  | Test <br/> (Release Candidate)        | Environment for testing integration with KSeF API 2.0, contains RC versions. | https://api-test.ksef.mf.gov.pl/docs/v2   | FA(2), FA(3), FA_PEF (3), FA_KOR_PEF (3)     |
| **DEMO**  | Pre-production (Demo)    | Environment matching production configuration, intended for final integration validation under conditions similar to production. | https://api-demo.ksef.mf.gov.pl/docs/v2   | FA(3), FA_PEF (3), FA_KOR_PEF (3)            |
| **PRD** | Production                        | Environment for issuing and receiving invoices with full legal validity, with guaranteed SLA and proper production data.                           | https://api.ksef.mf.gov.pl/docs/v2             | FA(3), FA_PEF (3), FA_KOR_PEF (3)            |



> <font color="red">Warning:</font> Test environments (TE/DEMO) are used exclusively for testing integration with the KSeF API. **Production invoices** or real entity data should not be sent to them.

In the test environment `TE`, authentication using self-signed certificates is allowed, which in practice means that many integrators can [authenticate](uwierzytelnianie.md#proces-uwierzytelniania) in the context of the same company.
For this reason, data entered in the `TE` environment is not isolated and may be shared between integrators.
Random NIP identifiers should be used for testing, avoiding any real data.

### Maintenance Work on Test Environments
In connection with the planned, systematic development of the National e-Invoice System (KSeF 2.0), **from October 1, 2025**, cyclical maintenance work may be carried out on the System's test environments.

This work will be performed between **4:00 PM and 6:00 PM**. During this time, temporary difficulties in accessing test environments may occur.

After the work is completed, **only changes affecting integration** will be published in the [changelog](api-changelog.md), e.g., API behavior changes, contract modifications, XSD schemas, limits, etc. Changes that do not affect the API and are not noticeable from an integration perspective, e.g., internal quality, performance, or security fixes, **may not be communicated** or will be presented collectively.
