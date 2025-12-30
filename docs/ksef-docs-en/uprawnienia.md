## Permissions
10.07.2025

### Introduction - Business Context
The KSeF system introduces an advanced permission management mechanism that forms the foundation for secure and compliant use of the system by various entities. Permissions determine who can perform specific operations in KSeF - such as issuing invoices, viewing documents, granting further permissions, or managing subordinate units.

### Permission Management Objectives:
- Data security - limiting access to information only to persons and entities that are formally authorized.
- Regulatory compliance - ensuring that operations are performed by appropriate entities in accordance with statutory requirements (e.g., VAT Act).
- Auditability - every operation related to granting or revoking permissions is recorded and can be analyzed.

### Who Grants Permissions?
Permissions can be granted by:

- the entity owner - role (Owner),
- administrator of a subordinate entity,
- administrator of a subordinate unit,
- administrator of an EU entity,
- entity administrator, i.e., another entity or person with the CredentialsManage permission.

In practice, this means that each organization must manage the permissions of its employees, e.g., granting permissions to an accounting department employee when hiring a new worker or revoking permissions when such an employee terminates employment.

### When Are Permissions Granted?
#### Examples:
- when starting cooperation with a new employee,
- when a company enters into cooperation, e.g., with an accounting office, it should grant invoice reading permissions to that accounting office so that the office can process the company's invoices,
- due to changes in relationships between entities.

### Structure of Granted Permissions:
Permissions are granted to:
1) Natural persons identified by PESEL, NIP, or certificate fingerprint - for working in KSeF:
    - in the context of the entity granting the permission (directly granted permissions) or
    - in the context of another entity or other entities:
        - in the context of a subordinate entity identified by NIP (subordinate local government unit or VAT group member),
        - in the context of a subordinate unit identified by an internal identifier,
        - in a composite context NIP-VAT UE linking a Polish entity with an EU entity authorized for self-invoicing on behalf of that Polish entity,
        - in the context of a specified entity identified by NIP - a client of the entity granting permissions (selective permissions granted indirectly),
        - in the context of all entities - clients of the entity granting permissions (general permissions granted indirectly).
2) Other entities - identified by NIP:
    - as end recipients of permissions to issue or view invoices,
    - as intermediaries - with the option to allow further delegation of permissions enabled, so that the authorized entity can grant permissions indirectly (see points IV and V above).

3) Other entities to act in their own context on behalf of the authorizing entity (entity permissions):
    - tax representatives,
    - entities authorized for self-invoicing,
    - entities authorized to issue VAT RR invoices.

Access to system functions depends on the context in which authentication occurred and on the scope of permissions granted to the authenticated entity/person in that context.

##  Glossary of Terms (regarding KSeF permissions)

| Term                          | Definition |
|---------------------------------|-----------|
| **Permission**                 | Authorization to perform specific operations in KSeF, e.g., `InvoiceWrite`, `CredentialsManage`. |
| **Owner**                       | Entity owner - a person who by default has full access to operations in the context of an entity with the same NIP identifier as recorded in the authentication method used; for the owner, the NIP-PESEL association also applies, so they can also authenticate with a method containing the associated PESEL number while retaining all owner permissions. |
| **Subordinate Entity Administrator**              | A person with permissions to manage permissions (`CredentialsManage`) in the context of a subordinate entity. Can grant permissions (e.g., `InvoiceWrite`). A subordinate entity can be, for example, a VAT group member. |
| **Subordinate Unit Administrator**              | A person with permissions to manage permissions (`CredentialsManage`) in a subordinate unit. Can grant permissions (e.g., `InvoiceWrite`). |
| **EU Entity Administrator**              | A person with permissions to manage permissions (`CredentialsManage`) in a composite context identified by NipVatUe. Can grant permissions (e.g., `InvoiceRead`). |
| **Intermediary Entity**   | An entity that received a permission with the flag `canDelegate = true` and can pass this permission further, i.e., grant permissions indirectly. These can only be `InvoiceWrite` and `InvoiceRead` permissions. |
| **Target Entity**  | The entity in whose context a given permission applies - e.g., a company serviced by an accounting office. |
| **Directly Granted**       | Permission granted directly to a given user or entity by the owner or administrator. |
| **Indirectly Granted**          | Permission granted by an intermediary for servicing another entity - only for `InvoiceRead` and `InvoiceWrite`. |
| **`canDelegate`**              | Technical flag (`true` / `false`) allowing delegation of permissions. Only `InvoiceRead` and `InvoiceWrite` can have `canDelegate = true`. Can only be used when granting permission to an entity for invoice handling |
| **`subjectIdentifier`**        | Data identifying the permission recipient (person or entity): `Nip`, `Pesel`, `Fingerprint`. |
| **`targetIdentifier` / `contextIdentifier`** | Data identifying the context in which the granted permission operates - e.g., client's NIP, internal identifier of an organizational unit. |
| **Fingerprint**                | The result of calculating the SHA-256 hash function on a qualified certificate. Allows recognition of the certificate of an entity possessing a permission granted to the certificate fingerprint. Used, among others, in identifying foreign persons or entities. |
| **InternalId**                 | Internal identifier of a subordinate unit in the KSeF system - a two-part identifier consisting of the NIP number and five digits `nip-5_digits`.  |
| **NipVatUe**                   | Composite identifier, i.e., a two-part identifier consisting of the Polish entity's NIP number and the EU entity's VAT UE number, separated by a separator `nip-vat_ue`. |
| **Revocation**                     | Operation of revoking a previously granted permission. |
| **`permissionId`**             | Technical identifier of a granted permission - required, among others, for revocation operations. |
| **`operationReferenceNumber`** | Operation identifier (e.g., granting or revoking permissions), returned by the API, used to check status. |
| **Operation Status**            | Current state of the permission granting/revoking process: `100`, `200`, `400`, etc. |

## Role and Permission Model (Permission Matrix)

The KSeF system allows assigning permissions precisely, taking into account the types of activities performed by users. Permissions can be granted both directly and indirectly - depending on the access delegation mechanism.

### Examples of Roles to Map Using Permissions:

| Role / Entity                          | Role Description                                                                                          | Possible Permissions                                                                 |
|----------------------------------------|-----------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------|
| **Entity Owner**      | Role possessed by default automatically by the owner. To be recognized by the system as an owner, one must authenticate with a vector with the same NIP identifier as the login context NIP or associated PESEL number           | `Owner` role covering all invoice and administrative permissions except `VatUeManage`. |
| **Entity Administrator**            | Natural person with rights to grant and revoke permissions to other users and/or appoint administrators of subordinate units/entities.           | `CredentialsManage`, `SubunitManage`, `Introspection`.                              |
| **Operator (accounting / invoicing)** | Person responsible for issuing or viewing invoices.                                        | `InvoiceWrite`, `InvoiceRead`.                                                      |
| **Authorized Entity**                | Another business entity that has been granted specific permissions to issue invoices on behalf of the entity, e.g., Tax Representative.             | `SelfInvoicing`, `RRInvoicing`, `TaxRepresentative`                             |
| **Intermediary Entity**              | Entity that received permissions with the delegation option (`canDelegate`) and can pass them on.       | `InvoiceRead`, `InvoiceWrite` with flag `canDelegate = true`.
| **EU Entity Administrator**     | Person identifying with a certificate having rights to grant and revoke permissions to other users within an EU entity associated with a given Polish entity.                                     | `InvoiceWrite`, `InvoiceRead`,                                    `VatUeManage`,  `Introspection`.                      |                      |
| **EU Entity Representative**     | Person identifying with a certificate acting on behalf of an EU entity associated with a given Polish entity.                                     | `InvoiceWrite`, `InvoiceRead`.                                                      |
| **Subordinate Unit Administrator** | User with the ability to appoint administrators in subordinate units or entities.               | `CredentialsManage`.                                                                    |

---

### Permission Classification by Type:

| Permission Type           | Example Values                                       | Can Be Granted Indirectly | Operational Description                                                              |
|--------------------------|------------------------------------------------------------|-------------------------------|------------------------------------------------------------------------------|
| **Invoice**             | `InvoiceWrite`, `InvoiceRead`                              | Yes (if `canDelegate=true`) | Invoice operations: sending, downloading                     |
| **Administrative**       | `CredentialsManage`, `SubunitManage`,  `VatUeManage`.                       | No                            | Managing permissions, subordinate units                      |
| **Entity**        | `SelfInvoicing`, `RRInvoicing`, `TaxRepresentative`        | No                            | Authorization of other entities to act (issue invoices) in one's own context on behalf of the authorizing entity         |
| **Technical**            | `Introspection`                                            | No                            | Access to operation and session history                                         |

---

## General and Selective Permissions

The KSeF system allows granting selected permissions in a **general** or **selective (individual)** manner, enabling flexible management of access to data of many business partners.

###  Selective (Individual) Permissions

Selective permissions are those granted by an intermediary entity (e.g., an accounting office) in relation to a **specific target entity (partner)**. They allow limiting the scope of access only to a selected client or organizational unit.

**Example:**
Accounting office XYZ received from company ABC the `InvoiceRead` permission with the flag `canDelegate = true`. Now it can pass this permission to its employee, but only in the context of company ABC - other companies serviced by XYZ are not covered by this access.

**Selectivity Characteristics:**
- It is necessary to specify `targetIdentifier` (e.g., partner's `Nip`).
- The permission recipient operates only in the context of the specified entity.
- Does not grant access to data of other partners of the intermediary entity.

---

###  General Permissions

General permissions are those granted without specifying a specific partner, meaning the recipient gains access to operations in the context of **all entities whose data the intermediary entity processes**.

**Example:**
Entity A has `InvoiceRead` permission with `canDelegate = true` for many clients. It passes a general `InvoiceRead` permission to employee B - B can now act on behalf of any of A's clients (e.g., view invoices of all contractors).

**Generality Characteristics:**
- Target entity identifier type `targetIdentifier` is `AllPartners`.
- Access covers all entities serviced by the intermediary.
- Used in cases of mass integration, large shared service centers, or accounting systems.

---

### Technical Notes and Limitations

- The mechanism applies only to `InvoiceRead` and `InvoiceWrite` permissions granted indirectly.
- In practice, the difference lies in the presence (selective) or absence (general) of `targetIdentifier` entity in the `POST /permissions/indirect/grants` request body.
- The system does not allow combining general and selective granting in a single call - separate operations must be performed.
- General permissions should be used with caution, especially in production environments, due to their broad scope.

---

### Permission Assignment Structure:

1. **Direct Granting** - e.g., administrator of entity A assigns `InvoiceWrite` permission to a natural person in the context of entity A.
2. **Granting with Further Delegation Option** - e.g., administrator of entity A grants entity B (intermediary) `InvoiceRead` permission with `canDelegate=true`, which allows administrator of entity B to grant `InvoiceRead` to entity/person C.
3. **Indirect Granting** - using the dedicated endpoint /permissions/indirect/grants, where the administrator of intermediary B, who received permission with delegation from entity A, grants permissions on behalf of target entity A to entity/person C.

---

### Example Permission Matrix:

| User / Entity       | InvoiceWrite | InvoiceRead | CredentialsManage | SubunitManage | TaxRepresentative |
|----------------------------|--------------|-------------|--------------------|----------------|--------------------|
| Anna Kowalska (PESEL)      | Yes           | Yes          | No                 | No             | No                 |
| Accounting Office XYZ (NIP) | Yes (with delegation)          | Yes (with delegation) | No                 | No             | No                 |
| Jan Nowak (Identified by certificate)   | Yes           | Yes          | No                 | No             | No                 |
| Accounting Dept Admin (PESEL)           | No           | No          | Yes                 | Yes             | No                 |
| Parent Company i.e. owner (NIP)         | Yes           | Yes          | Yes                 | Yes             | Yes                 |
| VAT Group Admin (PESEL)          | No           | No          | No                 | Yes             | No                 |
| Tax Representative (NIP)          | No           | No          | No                 | No             | Yes                 |

---

### Roles or Permissions Required for Granting Permissions

| Granting Permissions:                        | Required Role or Permission                      |
|-------------------------------------------|---------------------------------------------------|
| to a natural person for working in KSeF      | `Owner` or `CredentialsManage`                   |
| to an entity for invoice handling           | `Owner` or `CredentialsManage`                   |
| entity permissions | `Owner` or `CredentialsManage`                   |
| for invoice handling - indirectly              | `Owner` or `CredentialsManage`    |
| to subordinate unit administrator   | `SubunitManage`                                   |
| to EU entity administrator      | `Owner` or `CredentialsManage`    |
| to EU entity representative     | `VatUeManage`    |
---

### Identifier Limitations (`subjectIdentifier`, `contextIdentifier`)

| Identifier Type | Identifies | Notes |
|--------------------|---------------------|-------|
| `Nip`              | Domestic entity     | For entities registered in Poland and natural persons |
| `Pesel`            | Natural person       | Required, among others, when granting permissions to employees using a trusted profile or qualified certificate with PESEL number  |
| `Fingerprint`      | Certificate owner      | Used when the qualified certificate does not contain NIP or PESEL identifier and when identifying administrators or representatives of EU entities   |
| `NipVatUe`         | EU entities associated with Polish entities       | Required when granting permissions to administrators and representatives of EU entities |
| `InternalId`       | Subordinate units  | Used in entities with a structure composed of subordinate units |

---

### API Functional Limitations

- The same permission cannot be granted twice - the API may return an error or ignore the duplicate.
- Executing a permission granting operation does not result in immediate access - the operation is asynchronous and must be correctly processed by the system (operation status should be checked).

---

### Time Limitations

- A granted permission remains active until it is revoked.
- Implementing time limitations requires logic on the client system side (e.g., permission revocation schedule).


## Granting Permissions


### Granting Permissions to Natural Persons for Working in KSeF.

Within organizations using KSeF, it is possible to grant permissions to specific natural persons - e.g., employees of the accounting or IT department. Permissions are assigned to a person based on their identifier (PESEL, NIP, or Fingerprint). Permissions can include both operational activities (e.g., issuing invoices) and administrative ones (e.g., managing permissions). This section describes how to grant such permissions via the API and the permission requirements on the granting party's side.

POST [/permissions/persons/grants](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Nadawanie-uprawnien/paths/~1api~1v2~1permissions~1persons~1grants/post)


| Field                                       | Value                                         |
| :----------------------------------------- | :---------------------------------------------- |
| `subjectIdentifier`                        | Entity or natural person identifier. `"Nip"`, `"Pesel"`, `"Fingerprint"`             |
| `permissions`                               | Permissions to grant. `"CredentialsManage"`, `"CredentialsRead"`, `"InvoiceWrite"`, `"InvoiceRead"`, `"Introspection"`, `"SubunitManage"`, `"EnforcementOperations"`		   |
| `description`                              | Text value (description)              |


List of permissions that can be granted to a natural person:


| Permission | Description |
| :------------------ | :---------------------------------- |
| `CredentialsManage` | Managing permissions |
| `CredentialsRead` | Viewing permissions |
| `InvoiceWrite` | Issuing invoices |
| `InvoiceRead` | Viewing invoices |
| `Introspection` | Viewing session history |
| `SubunitManage` | Managing subordinate units |
| `EnforcementOperations` | Performing enforcement operations |




Example in C#:
[KSeF.Client.Tests.Core\E2E\Permissions\PersonPermission\PersonPermissionE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/Permissions/PersonPermission/PersonPermissionE2ETests.cs)

```csharp
GrantPermissionsPersonRequest request = GrantPersonPermissionsRequestBuilder
    .Create()
    .WithSubject(subject)
    .WithPermissions(
        StandardPermissionType.InvoiceRead,
        StandardPermissionType.InvoiceWrite)
    .WithDescription(description)
    .Build();

OperationResponse response =
    await KsefClient.GrantsPermissionPersonAsync(request, accessToken, CancellationToken);
```

Example in Java:
[PersonPermissionIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/PersonPermissionIntegrationTest.java)

```java

GrantPersonPermissionsRequest request = new GrantPersonPermissionsRequestBuilder()
        .withSubjectIdentifier(new PersonPermissionsSubjectIdentifier(PersonPermissionsSubjectIdentifier.IdentifierType.PESEL, personValue))
        .withPermissions(List.of(PersonPermissionType.INVOICEWRITE, PersonPermissionType.INVOICEREAD))
        .withDescription("e2e test grant")
        .build();

OperationResponse response = ksefClient.grantsPermissionPerson(request, accessToken);
```

Permissions can be granted by someone who is:
- an owner
- has `CredentialsManage` permission
- a subordinate unit administrator who has `SubunitManage`
- an EU entity administrator who has `VatUeManage`


---
### Granting Entities Permissions for Invoice Handling

KSeF allows granting permissions to entities that will process invoices on behalf of a given organization - e.g., accounting offices, shared service centers, or outsourcing companies. InvoiceRead and InvoiceWrite permissions can be granted directly and, if needed - with the option for further delegation (flag `canDelegate`). This section discusses the mechanism for granting these permissions, required roles, and example implementations.

POST [/permissions/entities/grants](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Nadawanie-uprawnien/paths/~1api~1v2~1permissions~1entities~1grants/post)


* **InvoiceWrite (Issuing invoices)**: This permission allows sending invoice files in XML format to the KSeF system. After successful verification and KSeF number assignment, these files become structured invoices.
* **InvoiceRead (Viewing invoices)**: With this permission, an entity can download invoice lists within a given context, download invoice contents, invoices, report abuse, and generate and view collective payment identifiers.
* **InvoiceWrite** and **InvoiceRead** permissions can be granted directly to entities by the authorizing entity. The API client granting these permissions directly must have **CredentialsManage** permission or the **Owner** role. When granting permissions to entities, it is possible to set the `canDelegate` flag to `true` for **InvoiceRead** and **InvoiceWrite**, which allows further, indirect delegation of this permission.



| Field                                       | Value                                         |
| :----------------------------------------- | :---------------------------------------------- |
| `subjectIdentifier`                        | Entity identifier. `"Nip"`               |
| `permissions`                               | Permissions to grant. `"InvoiceWrite"`, `"InvoiceRead"`			   |
| `description`                              | Text value (description)              |

Example in C#:
[KSeF.Client.Tests.Core\E2E\Permissions\EntityPermission\EntityPermissionE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/Permissions/EntityPermission/EntityPermissionE2ETests.cs)

```csharp
GrantPermissionsEntityRequest request = GrantEntityPermissionsRequestBuilder
    .Create()
    .WithSubject(subject)
    .WithPermissions(
        Permission.New(StandardPermissionType.InvoiceRead, true),
        Permission.New(StandardPermissionType.InvoiceWrite, false)
    )
    .WithDescription(description)
    .Build();

OperationResponse response = await KsefClient.GrantsPermissionEntityAsync(request, accessToken, CancellationToken);
```
Example in Java:
[EntityPermissionIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/EntityPermissionIntegrationTest.java)

```java
GrantEntityPermissionsRequest request = new GrantEntityPermissionsRequestBuilder()
        .withPermissions(List.of(
                new EntityPermission(EntityPermissionType.INVOICE_READ, true),
                new EntityPermission(EntityPermissionType.INVOICE_WRITE, false)))
        .withDescription(DESCRIPTION)
        .withSubjectIdentifier(new SubjectIdentifier(SubjectIdentifier.IdentifierType.NIP, targetNip))
        .build();

OperationResponse response = ksefClient.grantsPermissionEntity(request, accessToken);
```

---
### Granting Entity Permissions

For selected invoicing processes, KSeF provides so-called entity permissions, which apply in the context of invoicing on behalf of another entity (`TaxRepresentative`, `SelfInvoicing`, `RRInvoicing`). These permissions can only be granted by the owner or an administrator with `CredentialsManage`. This section presents how to grant them, their application, and technical limitations.

POST [/permissions/authorizations/grants](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Nadawanie-uprawnien/paths/~1api~1v2~1permissions~1authorizations~1grants/post)

Used to grant so-called entity permissions, such as `SelfInvoicing` (self-invoicing), `RRInvoicing` (RR self-invoicing), or `TaxRepresentative` (tax representative operations).

Permission Characteristics:

These are entity permissions, meaning they are relevant when sending invoice files by an entity and verified during their validation process. The relationship between the entity and invoice data is verified. They can be changed during a session.

Required permissions for granting: ```CredentialsManage``` or ```Owner```.

| Field                                       | Value                                         |
| :----------------------------------------- | :---------------------------------------------- |
| `subjectIdentifier`                        | Entity identifier. `"Nip"`               |
| `permissions`                               | Permissions to grant. `"SelfInvoicing"`, `"RRInvoicing"`, `"TaxRepresentative"`			   |
| `description`                              | Text value (description)              |


Example in C#:
[KSeF.Client.Tests.Core\E2E\Permissions\ProxyPermission\AuthorizationPermissionsE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/Permissions/ProxyPermission/AuthorizationPermissionsE2ETests.cs)

```csharp
GrantPermissionsAuthorizationRequest grantPermissionsAuthorizationRequest = GrantAuthorizationPermissionsRequestBuilder
    .Create()
    .WithSubject(new AuthorizationSubjectIdentifier
    {
        Type = AuthorizationSubjectIdentifierType.PeppolId,
        Value = peppolId
    })
    .WithPermission(AuthorizationPermissionType.PefInvoicing)
    .WithDescription($"E2E: Granting permission to issue PEF invoices for company {companyNip} (at request of {peppolId})")
    .Build();

OperationResponse operationResponse = await KsefClient
    .GrantsAuthorizationPermissionAsync(grantPermissionAuthorizationRequest,
    accessToken, CancellationToken);
```

Example in Java:
[ProxyPermissionIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/ProxyPermissionIntegrationTest.java)

```java
GrantAuthorizationPermissionsRequest request = new GrantAuthorizationPermissionsRequestBuilder()
        .withSubjectIdentifier(new SubjectIdentifier(SubjectIdentifier.IdentifierType.NIP, subjectNip))
        .withPermission(InvoicePermissionType.SELF_INVOICING)
        .withDescription("e2e test grant")
        .build();

OperationResponse response = ksefClient.grantsPermissionsProxyEntity(request, accessToken);
```
---
### Granting Permissions Indirectly

The indirect permission granting mechanism enables the operation of a so-called intermediary entity, which - based on previously obtained delegations - can pass selected permissions further, in the context of another entity. This most often applies to accounting offices that service many clients. This section describes the conditions that must be met to use this functionality and presents the data structure required to perform such granting.

`InvoiceWrite` and `InvoiceRead` are the only permissions that can be granted indirectly. This means that an intermediary entity can grant these permissions to another entity (authorized), which will apply in the context of the target entity (partner). These permissions can be selective (for a specific partner) or general (for all partners of the intermediary entity). For selective granting, the target entity identifier should specify the type `"Nip"` and the value of the specific NIP number. For general permissions, the target entity identifier should specify the type `"AllPartners"`, without a filled `value` field.

POST [/permissions/indirect/grants](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Nadawanie-uprawnien/paths/~1api~1v2~1permissions~1indirect~1grants/post)



| Field                                       | Value                                         |
| :----------------------------------------- | :---------------------------------------------- |
| `subjectIdentifier`                        | Natural person identifier. `"Nip"`, `"Pesel"`, `"Fingerprint"`               |
| `targetIdentifier`                        | Target entity identifier. `"Nip"` or `null`              |
| `permissions`                               | Permissions to grant. `"InvoiceRead"`, `"InvoiceWrite"`			   |
| `description`                              | Text value (description)              |

Example in C#:
[KSeF.Client.Tests.Core\E2E\Permissions\IndirectPermission\IndirectPermissionE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/Permissions/IndirectPermission/IndirectPermissionE2ETests.cs)

```csharp
GrantPermissionsIndirectEntityRequest request = GrantIndirectEntityPermissionsRequestBuilder
    .Create()
    .WithSubject(subject)
    .WithContext(new TargetIdentifier { Type = TargetIdentifierType.Nip, Value = targetNip })
    .WithPermissions(StandardPermissionType.InvoiceRead, StandardPermissionType.InvoiceWrite)
    .WithDescription(description)
    .Build();

OperationResponse grantOperationResponse = await KsefClient.GrantsPermissionIndirectEntityAsync(request, accessToken, CancellationToken);
```

Example in Java:
[IndirectPermissionIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/IndirectPermissionIntegrationTest.java)

```java
GrantIndirectEntityPermissionsRequest request = new GrantIndirectEntityPermissionsRequestBuilder()
        .withSubjectIdentifier(new SubjectIdentifier(SubjectIdentifier.IdentifierType.NIP, subjectNip))
        .withTargetIdentifier(new TargetIdentifier(TargetIdentifier.IdentifierType.NIP, targetNip))
        .withPermissions(List.of(INVOICE_WRITE))
        .withDescription("E2E indirect grantE2E indirect grant")
        .build();

OperationResponse response = ksefClient.grantsPermissionIndirectEntity(request, accessToken);

```
---
### Granting Subordinate Entity Administrator Permissions

The organizational structure of an entity may include subordinate units or entities - e.g., branches, departments, subsidiaries, VAT group members, and local government units. KSeF allows assigning permissions to manage such units. Having the `SubunitManage` permission is required. This section presents how to grant administrative permissions in the context of a subordinate unit or subordinate entity, taking into account the `InternalId` or `Nip` identifier.

POST [/permissions/subunits/grants](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Nadawanie-uprawnien/paths/~1api~1v2~1permissions~1subunits~1grants/post)



Required permissions for granting:

- The user who wants to grant these permissions must have the ```SubunitManage``` permission (Managing subordinate units).

| Field                                       | Value                                         |
| :----------------------------------------- | :---------------------------------------------- |
| `subjectIdentifier`                        | Natural person or entity identifier. `"Nip"`, `"Pesel"`, `"Fingerprint"`               |
| `contextIdentifier`                        | Subordinate entity identifier. `"Nip"`, `InternalId`              |
| `subunitName`                              | Subordinate unit name              |
| `description`                              | Text value (description)              |

Example in C#:
[KSeF.Client.Tests.Core\E2E\Permissions\SubunitPermission\SubunitPermissionsE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/Permissions/SubunitPermission/SubunitPermissionsE2ETests.cs)

```csharp
GrantPermissionsSubunitRequest subunitGrantRequest =
    GrantSubunitPermissionsRequestBuilder
    .Create()
    .WithSubject(_fixture.SubjectIdentifier)
    .WithContext(new SubunitContextIdentifier
    {
        Type = SubunitContextIdentifierType.InternalId,
        Value = Fixture.UnitNipInternal
    })
    .WithSubunitName("E2E Test Subunit")
    .WithDescription("E2E test grant sub-unit")
    .Build();

OperationResponse operationResponse = await KsefClient
    .GrantsPermissionSubUnitAsync(grantPermissionsSubUnitRequest, accessToken, CancellationToken);
```
Example in Java:

[SubUnitPermissionIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/SubUnitPermissionIntegrationTest.java)

```java
SubunitPermissionsGrantRequest request = new SubunitPermissionsGrantRequestBuilder()
        .withSubjectIdentifier(new SubjectIdentifier(SubjectIdentifier.IdentifierType.NIP, subjectNip))
        .withContextIdentifier(new ContextIdentifier(ContextIdentifier.IdentifierType.INTERNALID, contextNip))
        .withDescription("e2e subunit test")
        .withSubunitName("test")
        .build();

OperationResponse response = ksefClient.grantsPermissionSubUnit(request, accessToken);

```
---
### Granting EU Entity Administrator Permissions

Granting EU entity administrator permissions in KSeF allows authorizing an entity or person designated by an EU entity that has the right to self-invoice on behalf of the Polish entity granting the permission. Executing this operation causes the person authorized in this way to gain the ability to log in with a composite context: `NipVatUe`, linking the Polish entity granting permission with the EU entity having the right to self-invoice. After granting EU entity administrator permissions, such a person will be able to perform invoice operations and manage the permissions of other persons (so-called EU entity representatives) within this composite context.

POST [/permissions/eu-entities/administration/grants](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Nadawanie-uprawnien/paths/~1api~1v2~1permissions~1eu-entities~1administration~1grants/post)



| Field                                       | Value                                         |
| :----------------------------------------- | :---------------------------------------------- |
| `subjectIdentifier`                        | Natural person or entity identifier. `"Nip"`, `"Pesel"`, `"Fingerprint"`               |
| `contextIdentifier`                        | Two-part identifier consisting of NIP number and VAT-UE number `{nip}-{vat_ue}`. `"NipVatUe"`              |
| `description`                              | Text value (description)              |

Example in C#:
[KSeF.Client.Tests.Core\E2E\Permissions\EuEntityPermission\EuEntityPermissionE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/Permissions/EuEntityPermission/EuEntityPermissionE2ETests.cs)

```csharp
GrantPermissionsEuEntityRequest grantPermissionsEuEntityRequest = GrantEUEntityPermissionsRequestBuilder
    .Create()
    .WithSubject(TestFixture.EuEntity)
    .WithSubjectName(EuEntitySubjectName)
    .WithContext(contextIdentifier)
    .WithDescription(EuEntityDescription)
    .Build();

OperationResponse operationResponse = await KsefClient
    .GrantsPermissionEUEntityAsync(grantPermissionsRequest, accessToken, CancellationToken);
```
Example in Java:
[EuEntityPermissionIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/EuEntityPermissionIntegrationTest.java)

```java
EuEntityPermissionsGrantRequest request = new GrantEUEntityPermissionsRequestBuilder()
        .withSubject(new SubjectIdentifier(SubjectIdentifier.IdentifierType.FINGERPRINT, euEntity))
        .withEuEntityName("Sample Subject Name")
        .withContext(new ContextIdentifier(ContextIdentifier.IdentifierType.NIP_VAT_UE, nipVatUe))
        .withDescription("E2E EU Entity Permission Test")
        .build();

OperationResponse response = ksefClient.grantsPermissionEUEntity(request, accessToken);

```
---
### Granting EU Entity Representative Permissions

An EU entity representative is a person acting on behalf of an EU-registered entity that needs access to KSeF for viewing or issuing invoices. Such permission can only be granted by a VAT UE administrator. This section presents the data structure and how to call the appropriate endpoint.

POST [/permissions/eu-entities/grants](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Nadawanie-uprawnien/paths/~1api~1v2~1permissions~1eu-entities~1grants/post)



| Field                                       | Value                                         |
| :----------------------------------------- | :---------------------------------------------- |
| `subjectIdentifier`                        | Entity identifier. `"Fingerprint"`               |
| `permissions`                               | Permissions to grant. `"InvoiceRead"`, `"InvoiceWrite"`			   |
| `description`                              | Text value (description)              |

Example in C#:
[KSeF.Client.Tests.Core\E2E\Permissions\EuAdministrationPermission\EuRepresentativePermissionE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/Permissions/EuAdministrationPermission/EuRepresentativePermissionE2ETests.cs)

```csharp
GrantPermissionsEuEntityRepresentativeRequest grantRepresentativePermissionsRequest = GrantEUEntityRepresentativePermissionsRequestBuilder
    .Create()
    .WithSubject(new Client.Core.Models.Permissions.EUEntityRepresentative.SubjectIdentifier
    {
        Type = Client.Core.Models.Permissions.EUEntityRepresentative.SubjectIdentifierType.Fingerprint,
        Value = euRepresentativeEntityCertificateFingerprint
    })
    .WithPermissions(
        StandardPermissionType.InvoiceWrite,
        StandardPermissionType.InvoiceRead
        )
    .WithDescription("Representative for EU Entity")
    .Build();

OperationResponse grantRepresentativePermissionResponse = await KsefClient.GrantsPermissionEUEntityRepresentativeAsync(grantRepresentativePermissionsRequest,
    euAuthInfo.AccessToken.Token,
    CancellationToken.None);
```
Example in Java:
[EuEntityRepresentativePermissionIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/EuEntityRepresentativePermissionIntegrationTest.java)

```java
GrantEUEntityRepresentativePermissionsRequest request = new GrantEUEntityRepresentativePermissionsRequestBuilder()
        .withSubjectIdentifier(new SubjectIdentifier(SubjectIdentifier.IdentifierType.FINGERPRINT, fingerprint))
        .withPermissions(List.of(EuEntityPermissionType.INVOICE_WRITE, EuEntityPermissionType.INVOICE_READ))
        .withDescription("Representative for EU Entity")
        .build();

OperationResponse response = ksefClient.grantsPermissionEUEntityRepresentative(request, accessToken);


```

## Revoking Permissions

The process of revoking permissions in KSeF is equally important as granting them - it ensures access control and enables quick response in situations such as employee role change, termination of cooperation with an external partner, or security policy violation. Permission revocation can be performed for each recipient category: natural person, entity, subordinate unit, EU representative, or EU administrator. This section discusses methods for revoking different types of permissions and required identifiers.

### Revoking Permissions

The standard method of revoking permissions, which can be used for most cases: natural persons, domestic entities, subordinate units, as well as EU representatives or EU administrators. The operation requires knowledge of `permissionId` and having the appropriate permission.

DELETE [/permissions/common/grants/\{permissionId\}](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Odbieranie-uprawnien/paths/~1api~1v2~1permissions~1common~1grants~1%7BpermissionId%7D/delete)

This method is used to revoke permissions such as:

- granted to natural persons for working in KSeF,
- granted to entities for invoice handling,
- granted indirectly,
- subordinate entity administrator,
- EU entity administrator,
- EU entity representative.

Example in C#:
[KSeF.Client.Tests.Core\E2E\Certificates\CertificatesE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/Certificates/CertificatesE2ETests.cs)
```csharp
OperationResponse operationResponse = await KsefClient.RevokeCommonPermissionAsync(permission.Id, accessToken, CancellationToken);
```

Example in Java:
[EntityPermissionIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/EntityPermissionIntegrationTest.java)

```java
OperationResponse response = ksefClient.revokeCommonPermission(permissionId, accessToken);
```
---
### Revoking Entity Permissions

For entity-type permissions (`SelfInvoicing`, `RRInvoicing`, `TaxRepresentative`), a separate revocation method applies - using an endpoint dedicated to authorization operations. These types of permissions are not transferable, so their revocation has an immediate effect and terminates the ability to perform invoice operations in that mode. Knowledge of `permissionId` is required.

DELETE [/permissions/authorizations/grants/\{permissionId\}](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Odbieranie-uprawnien/paths/~1api~1v2~1permissions~1authorizations~1grants~1%7BpermissionId%7D/delete)

This method is used to revoke permissions such as:

- self-invoicing,
- RR self-invoicing,
- tax representative operations.

Example in C#:
[KSeF.Client.Tests.Core\E2E\Permissions\EuEntityPermission\EuEntityPermissionE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/Permissions/EuEntityPermission/EuEntityPermissionE2ETests.cs)

```csharp
await ksefClient.RevokeAuthorizationsPermissionAsync(permissionId, accessToken, cancellationToken);
```

Example in Java:
[ProxyPermissionIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/ProxyPermissionIntegrationTest.java)

```java
OperationResponse response = ksefClient.revokeAuthorizationsPermission(operationId, accessToken);
```


## Searching Granted Permissions

KSeF provides a set of endpoints allowing querying the list of active permissions granted to users and entities. These mechanisms are essential for auditing, reviewing access status, and building administrative interfaces (e.g., for managing access structure in an organization). This section contains an overview of search methods by category of granted permissions.

---
### Retrieving Own Permissions List

The query allows retrieving a list of permissions held by the authenticated entity.
 This list includes permissions:
- granted directly in the current context
- granted by a parent entity
- granted indirectly, where the context is the intermediary or target entity
- granted to an entity for invoice handling (`"InvoiceRead"` and `"InvoiceWrite"`) by another entity, if the authenticated entity has owner permissions (`"Owner"`)

POST [/permissions/query/personal/grants](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Wyszukiwanie-nadanych-uprawnien/paths/~1api~1v2~1permissions~1query~1personal~1grants/post)

Example in C#:
[KSeF.Client.Tests.Core\E2E\Permissions\PersonPermission\PersonalPermissions_AuthorizedPesel_InNipContext_E2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/Permissions/PersonPermission/PersonalPermissions_AuthorizedPesel_InNipContext_E2ETests.cs)
```csharp
PersonalPermissionsQueryRequest query = new PersonalPermissionsQueryRequest
{
    ContextIdentifier = /*...*/,
    TargetIdentifier = /*...*/,
    PermissionTypes = /*...*/,
    PermissionState = /*...*/
};

PagedPermissionsResponse<PersonalPermission> searchedGrantedPersonalPermissions =
    await KsefClient.SearchGrantedPersonalPermissionsAsync(query, entityAuthorizationInfo.AccessToken.Token);
```

Example in Java:
[SearchPersonalGrantPermissionIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/SearchPersonalGrantPermissionIntegrationTest.java)

```java
QueryPersonalGrantResponse response = ksefClient.searchPersonalGrantPermission(request, pageOffset, pageSize, token.accessToken());

```
---
### Retrieving List of KSeF Work Permissions Granted to Natural Persons or Entities

The query allows retrieving a list of permissions granted to natural persons or entities - e.g., company employees. Filtering by permission type, state (`Active` / `Inactive`), as well as grantor and grantee identifiers is possible. This endpoint is used during onboarding, auditing, and monitoring personal permissions.

POST [/permissions/query/persons/grants](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Wyszukiwanie-nadanych-uprawnien/paths/~1api~1v2~1permissions~1query~1persons~1grants/post)

| Field                  | Description                                                                 |
| :-------------------- | :------------------------------------------------------------------- |
| `authorIdentifier`    | Identifier of the entity granting permissions.   ```Nip```, ```Pesel```, ```Fingerprint```, ```System```                      |
| `authorizedIdentifier`| Identifier of the entity that was granted permissions.      ```Nip```, ```Pesel```,```Fingerprint```             |
| `targetIdentifier`    | Target entity identifier (for indirect permissions).  ```Nip```, ```AllPartners```      |
| `permissionTypes`     | Permission types for filtering.   `"CredentialsManage"`, `"CredentialsRead"`, `"InvoiceWrite"`, `"InvoiceRead"`, `"Introspection"`, `"SubunitManage"`, `"EnforcementOperations"`  |
| `permissionState`     | Permission state.  ```Active``` / ```Inactive```                                                  |

Example in C#:
[KSeF.Client.Tests.Core\E2E\Permissions\SubunitPermission\SubunitPermissionsE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/Permissions/SubunitPermission/SubunitPermissionsE2ETests.cs)

```csharp
PagedPermissionsResponse<Client.Core.Models.Permissions.PersonPermission> response =
    await KsefClient
    .SearchGrantedPersonPermissionsAsync(
        personPermissionsQueryRequest,
        accessToken,
        pageOffset: 0,
        pageSize: 10,
        CancellationToken);
```

Example in Java:
[PersonPermissionIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/PersonPermissionIntegrationTest.java)

```java
PersonPermissionsQueryRequest request = new PersonPermissionsQueryRequestBuilder()
        .withQueryType(PersonPermissionQueryType.PERMISSION_GRANTED_IN_CURRENT_CONTEXT)
        .build();

QueryPersonPermissionsResponse response = ksefClient.searchGrantedPersonPermissions(request, pageOffset, pageSize, accessToken);


```
---
### Retrieving List of Subordinate Unit and Entity Administrator Permissions

This endpoint is used to retrieve information about administrators of subordinate units or subordinate entities (e.g., branches, VAT groups). It allows monitoring who has management permissions for a given subordinate structure, identified by `InternalId` or `Nip`.

POST [/permissions/query/subunits/grants](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Wyszukiwanie-nadanych-uprawnien/paths/~1api~1v2~1permissions~1query~1subunits~1grants/post)

| Field                  | Description                                                                 |
| :-------------------- | :------------------------------------------------------------------- |
| `subjectIdentifier`    | Subordinate entity identifier.   ```InternalId``` or `Nip`            |

Example in C#:
[KSeF.Client.Tests.Core\E2E\Permissions\SubunitPermission\SubunitPermissionsE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/Permissions/SubunitPermission/SubunitPermissionsE2ETests.cs)

```csharp
SubunitPermissionsQueryRequest subunitPermissionsQueryRequest = new SubunitPermissionsQueryRequest();
PagedPermissionsResponse<Client.Core.Models.Permissions.SubunitPermission> response =
    await KsefClient
    .SearchSubunitAdminPermissionsAsync(
        subunitPermissionsQueryRequest,
        accessToken,
        pageOffset: 0,
        pageSize: 10,
        CancellationToken);
```

Example in Java:
[SubUnitPermissionIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/SubUnitPermissionIntegrationTest.java)

```java
SubunitPermissionsQueryRequest request = new SubunitPermissionsQueryRequestBuilder()
        .withSubunitIdentifier(new SubunitPermissionsSubunitIdentifier(SubunitPermissionsSubunitIdentifier.IdentifierType.INTERNALID, subUnitNip))
        .build();

QuerySubunitPermissionsResponse response = ksefClient.searchSubunitAdminPermissions(request, pageOffset, pageSize, accessToken);


```
---
### Retrieving Entity Roles List

The endpoint returns the set of roles assigned to the context in which we are authenticated (i.e., on whose behalf the query is executed). The function is mainly used for automatic access verification by client systems.

GET [/permissions/query/entities/roles](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Wyszukiwanie-nadanych-uprawnien/paths/~1api~1v2~1permissions~1query~1entities~1roles/get)

Example in C#:
[KSeF.Client.Tests.Core\E2E\Permissions\SubunitPermission\SubunitPermissionsE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/Permissions/SubunitPermission/SubunitPermissionsE2ETests.cs)

```csharp
PagedRolesResponse<EntityRole> response =
    await KsefClient
    .SearchEntityInvoiceRolesAsync(
        accessToken,
        pageOffset: 0,
        pageSize: 10,
        CancellationToken);
```

Example in Java:
[SearchEntityInvoiceRoleIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/SearchEntityInvoiceRoleIntegrationTest.java)

```java
QueryEntityRolesResponse response = ksefClient.searchEntityInvoiceRoles(0, 10, token);
```
---
### Retrieving Subordinate Entities List

Allows obtaining information about related subordinate entities for the context in which we are authenticated (i.e., on whose behalf the query is executed). The function is mainly used to verify the structure of local government units or VAT groups.

POST [/permissions/query/subordinate-entities/roles](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Wyszukiwanie-nadanych-uprawnien/paths/~1api~1v2~1permissions~1query~1subordinate-entities~1roles/post)

| Field                     | Description                                                                                                              |
| :----------------------- | :---------------------------------------------------------------------------------------------------------------- |
| `subordinateEntityIdentifier`   | Identifier of the entity that was granted permissions. ```Nip```                                                     |                                               |


Example in C#:
[KSeF.Client.Tests.Core\E2E\Permissions\SubunitPermission\SubunitPermissionsE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/Permissions/SubunitPermission/SubunitPermissionsE2ETests.cs)

```csharp
SubordinateEntityRolesQueryRequest subordinateEntityRolesQueryRequest = new SubordinateEntityRolesQueryRequest();
PagedRolesResponse<SubordinateEntityRole> response =
    await KsefClient
    .SearchSubordinateEntityInvoiceRolesAsync(
        subordinateEntityRolesQueryRequest,
        accessToken,
        pageOffset: 0,
        pageSize: 10,
        CancellationToken);
```

Example in Java:
[SearchSubordinateQueryIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/SearchSubordinateQueryIntegrationTest.java)

```java
SubordinateEntityRolesQueryResponse response = ksefClient.searchSubordinateEntityInvoiceRoles(queryRequest, pageOffset, pageSize,accessToken);
```
---
### Retrieving Entity Invoice Handling Permissions List

This endpoint is used to review all entity permissions granted by the context in which we are authenticated or granted to the context in which we are authenticated. It supports filtering by permission type and recipient.



POST [/permissions/query/authorizations/grants](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Wyszukiwanie-nadanych-uprawnien/paths/~1api~1v2~1permissions~1query~1authorizations~1grants/post)

| Field                     | Description                                                                                                              |
| :----------------------- | :---------------------------------------------------------------------------------------------------------------- |
| `authorizingIdentifier`  | Identifier of the entity granting permissions.  ```Nip```                                                     |
| `authorizedIdentifier`   | Identifier of the entity that was granted permissions. ```Nip```                                                     |
| `queryType`              | Query type. Determines whether we query for granted or received permissions. ```Granted``` ```Received```            |
| `permissionTypes`        | Permission types for filtering.   `"SelfInvoicing"`, `"TaxRepresentative"`, `"RRInvoicing"`,                       |


Example in C#:
[KSeF.Client.Tests.Core\E2E\Permissions\SubunitPermission\SubunitPermissionsE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/Permissions/SubunitPermission/SubunitPermissionsE2ETests.cs)

```csharp
PagedAuthorizationsResponse<AuthorizationGrant> response =
        await KsefClient
        .SearchEntityAuthorizationGrantsAsync(
            entityAuthorizationsQueryRequest,
            accessToken,
            pageOffset: 0,
            pageSize: 10,
            CancellationToken);
```

Example in Java:
[ProxyPermissionIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/ProxyPermissionIntegrationTest.java)

```java
        EntityAuthorizationPermissionsQueryRequest request = new EntityAuthorizationPermissionsQueryRequestBuilder()
        .withQueryType(QueryType.GRANTED)
        .build();

QueryEntityAuthorizationPermissionsResponse response = ksefClient.searchEntityAuthorizationGrants(request, pageOffset, pageSize, accessToken);


```
---
### Retrieving Permissions List for Administrators or Representatives of EU Entities Authorized for Self-Invoicing

EU entities can also have permissions assigned for using KSeF. In this section, it is possible to retrieve information about access granted to them, taking into account VAT UE identifiers and certificate fingerprint.

POST [/permissions/query/eu-entities/grants](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Wyszukiwanie-nadanych-uprawnien/paths/~1api~1v2~1permissions~1query~1eu-entities~1grants/post)

| Field                        | Description                                                                 |
| :-------------------------- | :------------------------------------------------------------------- |
| `vatUeIdentifier`           | VAT UE identifier.                                                |
| `authorizedFingerprintIdentifier` | Certificate fingerprint of the authorized entity.                      |
| `permissionTypes`           | Permission types for filtering. Possible values are: `VatUeManage`, `InvoiceWrite`, `InvoiceRead`, `Introspection`. |

Example in C#:
[KSeF.Client.Tests.Core\E2E\Permissions\SubunitPermission\SubunitPermissionsE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/Permissions/SubunitPermission/SubunitPermissionsE2ETests.cs)

```csharp
PagedPermissionsResponse<Client.Core.Models.Permissions.EuEntityPermission> response =
    await KsefClient
    .SearchGrantedEuEntityPermissionsAsync(
        euEntityPermissionsQueryRequest,
        accessToken,
        pageOffset: 0,
        pageSize: 10,
        CancellationToken);
```

Example in Java:
[EuEntityPermissionIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/EuEntityPermissionIntegrationTest.java)

```java
EuEntityPermissionsQueryRequest request = new EuEntityPermissionsQueryRequestBuilder()
   .withAuthorizedFingerprintIdentifier(subjectContext)
   .build();

QueryEuEntityPermissionsResponse response = createKSeFClient().searchGrantedEuEntityPermissions(request, pageOffset, pageSize, accessToken);
```

## Operations

The National e-Invoice System allows tracking and verifying the status of operations related to permission management. Each permission grant or revocation is performed as an asynchronous operation, whose status can be monitored using a unique reference identifier (`operationReferenceNumber`). This section presents the mechanism for retrieving operation status and its interpretation in the context of automation and verification of administrative actions in KSeF.

### Retrieving Operation Status

After granting or revoking a permission, the system returns an operation reference number (`operationReferenceNumber`). This identifier allows checking the current state of request processing: whether it completed successfully, whether an error occurred, or whether processing is still ongoing. This information can be crucial in supervisory systems, automatic retry logic, or administrative action reporting. This section presents an example of an API call used to retrieve operation status.

GET [/permissions/operations/{operationReferenceNumber}](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Operacje/paths/~1api~1v2~1permissions~1operations~1%7BoperationReferenceNumber%7D/get)

Each permission grant operation returns an operation identifier that should be used to check the status of that operation.

Example in C#:
[KSeF.Client.Tests.Core\E2E\Permissions\SubunitPermission\SubunitPermissionsE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/Permissions/SubunitPermission/SubunitPermissionsE2ETests.cs)

```csharp
var operationStatus = await ksefClient.OperationsStatusAsync(referenceNumber, accessToken, cancellationToken);
```

Example in Java:
[EuEntityPermissionIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/EuEntityPermissionIntegrationTest.java)

```java
PermissionStatusInfo status = ksefClient.permissionOperationStatus(referenceNumber, accessToken);
```

### Checking Consent Status for Issuing Invoices with Attachments

Consent is required for issuing invoices containing attachments and applies within the current context (`ContextIdentifier`) used during authentication. Consent is granted outside the API, exclusively in the e-Tax Office (e-Urzad Skarbowy) service, and applications can be submitted from January 1, 2026. The API does not provide a consent submission operation.

GET [/permissions/attachments/status](https://ksef-test.mf.gov.pl/docs/v2/index.html#tag/Operacje/paths/~1api~1v2~1permissions~1attachments~1status/get)

Returns the consent status for the current context. If consent is not active, an invoice with an attachment sent to the KSeF API will be rejected.

Example in C#:
[KSeF.Client.Tests.Core\E2E\TestData\TestDataE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/TestData/TestDataE2ETests.cs)
```csharp
PermissionsAttachmentAllowedResponse attachmentPermissionStatus = await KsefClient.GetAttachmentPermissionStatusAsync(authOperationStatusResponse.AccessToken.Token)
```

Example in Java:
[PermissionAttachmentStatusIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/PermissionAttachmentStatusIntegrationTest.java)

```java
PermissionAttachmentStatusResponse trueResponse = ksefClient.checkPermissionAttachmentInvoiceStatus(token.accessToken());
```

**Test Environment**
On the test environment, the POST `/testdata/attachment` endpoint is available, which grants the ability to send invoices with attachments for the specified entity. This endpoint is only for simulating consent granting in tests and works within the scope of the current context.
