# KSeF Number - Structure and Validation

The KSeF number is a unique invoice identifier assigned by the system. It is always **35 characters** long and is globally unique - it unambiguously identifies each invoice in KSeF.

## General Structure of the Number
```
9999999999-RRRRMMDD-FFFFFFFFFFFF-FF
```
Where:
- `9999999999` - Seller's NIP (10 digits),
- `RRRRMMDD` - Date of invoice acceptance (year, month, day) for further processing,
- `FFFFFFFFFFFF` - Technical part consisting of 12 characters in hexadecimal notation, only [0-9 A-F], uppercase letters,
- `FF` - CRC-8 checksum - 2 characters in hexadecimal notation, only [0-9 A-F], uppercase letters.

## Example
```
5265877635-20250826-0100001AF629-AF
```
- `5265877635` - Seller's NIP,
- `20250826` - Date of invoice acceptance for further processing,
- `0100001AF629` - Technical part,
- `AF` - CRC-8 checksum.

## KSeF Number Validation

The validation process includes:
1. Checking if the number has **exactly 35 characters**.
2. Separating the data part (32 characters) and the checksum (2 characters).
3. Calculating the checksum from the data part using the **CRC-8 algorithm**.
4. Comparing the calculated checksum with the value in the number.

## CRC-8 Algorithm

The **CRC-8** algorithm is used to calculate the checksum with the following parameters:

- **Polynomial:** `0x07`
- **Initial value:** `0x00`
- **Result format:** 2-character hexadecimal notation (HEX, uppercase letters)

Example: if the calculated checksum is `0x46`, `"46"` will be added to the KSeF number.

## Example in C#:
```csharp
using KSeF.Client.Core;

bool isValid = KsefNumberValidator.IsValid(ksefNumber, out string message);
```

## Example in Java:
[OnlineSessionIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/OnlineSessionIntegrationTest.java)

```java
KSeFNumberValidator.ValidationResult result = KSeFNumberValidator.isValid(ksefNumber);

```
