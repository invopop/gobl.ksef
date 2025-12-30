## Example Scenarios
05.08.2025

### Scenario 1 – Bailiff

If you want to use the KSeF system on the test environment as a natural person with bailiff permissions, you need to add such a person using the `/v2/testdata/person` endpoint, setting the *isBailiff* flag to **true**.

Example JSON:
```json
{
  "nip": "7980332920",
  "pesel": "30112206276",
  "description": "Bailiff",
  "isBailiff": true
}
```

As a result of this operation, the person logging in within the context of the specified NIP, using their PESEL or NIP number, will receive owner permissions (**Owner**) and enforcement permissions (**EnforcementOperations**), which will enable them to use the system from the bailiff's perspective.

---

### Scenario 2 – Sole Proprietorship

If you want to use the KSeF system on the test environment as a sole proprietorship, you need to add such a person using the `/v2/testdata/person` endpoint, setting the *isBailiff* flag to **false**.

Example JSON:
```json
{
  "nip": "7980332920",
  "pesel": "30112206276",
  "description": "Sole Proprietorship",
  "isBailiff": false
}
```

As a result of this operation, the person logging in within the context of the specified NIP, using their PESEL or NIP number, will receive owner permissions (**Owner**), which will enable them to use the system from the sole proprietorship's perspective.

---

### Scenario 3 – VAT Group

If you want to create a VAT group structure on the test environment and grant permissions to the group administrator and administrators of its members, the first step is to create the entity structure using the `/v2/testdata/subject` endpoint, specifying the parent entity's NIP and the subordinate entities.

Example JSON:
```json
{
  "subjectNip": "3755747347",
  "subjectType": "VatGroup",
  "description": "VAT Group",
  "subunits": [
    {
      "subjectNip": "4972530874",
      "description": "NIP 4972530874: VAT group member for 3755747347"
    },
    {
      "subjectNip": "8225900795",
      "description": "NIP 8225900795: VAT group member for 3755747347"
    }
  ]
}
```

As a result of this operation, the specified entities and their relationships will be created in the system. Next, you need to grant permissions to a person within the context of the VAT group's NIP, according to the ZAW-FA rules. This operation can be performed using the `/v2/testdata/permissions` method.

Example JSON for a person authorized within the VAT group context:
```json
{
  "contextIdentifier": {
    "value": "3755747347",
    "type": "nip"
  },
  "authorizedIdentifier": {
    "value": "38092277125",
    "type": "pesel"
  },
  "permissions": [
    {
      "permissionType": "InvoiceRead",
      "description": "working in context 3755747347: authorized PESEL: 38092277125, Adam Abacki"
    },
    {
      "permissionType": "InvoiceWrite",
      "description": "working in context 3755747347: authorized PESEL: 38092277125, Adam Abacki"
    },
    {
      "permissionType": "Introspection",
      "description": "working in context 3755747347: authorized PESEL: 38092277125, Adam Abacki"
    },
    {
      "permissionType": "CredentialsRead",
      "description": "working in context 3755747347: authorized PESEL: 38092277125, Adam Abacki"
    },
    {
      "permissionType": "CredentialsManage",
      "description": "working in context 3755747347: authorized PESEL: 38092277125, Adam Abacki"
    },
    {
      "permissionType": "SubunitManage",
      "description": "working in context 3755747347: authorized PESEL: 38092277125, Adam Abacki"
    }
  ]
}
```

This operation can be performed both for the VAT group (as shown above) and for VAT group members. Note that while this is the only way to grant initial permissions for the VAT group itself, it is not necessary for group members. For group members, you can use the standard `/v2/permissions/subunit/grants` endpoint to appoint administrators for VAT group members.

Alternatively, you can use the test data endpoint described above. Example JSON for granting `CredentialsManage` permission to a group member's administrator:
```json
{
  "contextIdentifier": {
    "value": "4972530874",
    "type": "nip"
  },
  "authorizedIdentifier": {
    "value": "3388912629",
    "type": "nip"
  },
  "permissions": [
    {
      "permissionType": "CredentialsManage",
      "description": "working in context 4972530874: authorized NIP: 3388912629, Bogdan Babacki"
    }
  ]
}
```

Thanks to this operation, the VAT group member's representative gains the ability to grant permissions to themselves or other persons (e.g., employees) in the standard way, through the KSeF system.

### Scenario 4 – Enabling the Ability to Send Invoices with Attachments
On the test environment, you can simulate an entity that has the ability to send invoices with attachments enabled. The operation should be performed using the /testdata/attachment endpoint.

```json
{
  "nip": "4972530874"
}
```

As a result, the entity with NIP 4972530874 will receive the ability to send invoices containing attachments.

### Scenario 5 – Disabling the Ability to Send Invoices with Attachments
To test a situation where a given entity no longer has the ability to send invoices with attachments, use the /testdata/attachment/revoke endpoint.

```json
{
  "nip": "4972530874"
}
```

As a result, the entity with NIP 4972530874 loses the ability to send invoices containing attachments.
