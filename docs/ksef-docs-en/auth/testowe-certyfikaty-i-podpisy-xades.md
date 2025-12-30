# Test Certificates and XAdES Signatures

This guide shows how to **quickly run** the console demo application [`KSeF.Client.Tests.CertTestApp`](https://github.com/CIRFMF/ksef-client-csharp) to:
- generate a **test (self-signed) certificate** for the KSeF test environment,
- build and **XAdES sign** the `AuthTokenRequest` document,
- send the signed document to KSeF and **obtain access tokens** (JWT).

> **Note**
> - Self-signed certificates are **only allowed** in the **test** environment.
> - The data in the examples (NIP, reference number, tokens) are **fictitious** and are for demonstration purposes only.

---

## Prerequisites
- **.NET 10 SDK**
- Git
- Windows or Linux

---

## What does the application do?
- Retrieves a **challenge** from KSeF.
- Builds the `AuthTokenRequest` XML document.
- **Signs** the `AuthTokenRequest` document in **XAdES** format.
- Sends the signed document to KSeF and receives `referenceNumber` + `authenticationToken`.
- **Polls the authentication operation status** until completion.
- Upon success, retrieves the token pair: `accessToken` and `refreshToken` (JWT).
- Saves artifacts (including **test certificate** and **signed XML**) to files if `file` output is selected.

---

## Windows

1. **Install .NET 10 SDK**:
   ```powershell
   winget install Microsoft.DotNet.SDK.10
   ```
   Alternatively: download the installer from the .NET website.

2. **Open a new terminal window** (PowerShell/CMD).

3. **Verify installation**:
   ```powershell
   dotnet --version
   ```
   Expected version number: `10.x.x`.

4. **Clone the repository and navigate to the project**:
   ```powershell
   git clone https://github.com/CIRFMF/ksef-client-csharp.git
   cd ksef-client-csharp/KSeF.Client.Tests.CertTestApp
   ```

5. **Run (default random NIP, output to screen)**:
   ```powershell
   dotnet run --framework net10.0
   ```

6. **Run with parameters**:
   - `--output` - `screen` (default) or `file` (save results to files),
   - `--nip` {nip_number} - e.g., `--nip 8976111986`,
   - optionally: `--no-startup-warnings`.

   ```powershell
   dotnet run --framework net10.0 --output file --nip 8976111986 --no-startup-warnings
   ```

---

## Linux (Ubuntu/Debian)

1. **Add Microsoft repository and update packages**:
   ```bash
   wget https://packages.microsoft.com/config/ubuntu/$(lsb_release -rs)/packages-microsoft-prod.deb -O packages-microsoft-prod.deb
   sudo dpkg -i packages-microsoft-prod.deb
   sudo apt-get update
   ```

2. **Install .NET 10 SDK**:
   ```bash
   sudo apt-get install -y dotnet-sdk-10.0
   ```

3. **Refresh shell environment or open a new terminal**:
   ```bash
   source ~/.bashrc
   ```

4. **Verify installation**:
   ```bash
   dotnet --version
   ```
   Expected version number: `10.x.x`.

5. **Clone the repository and navigate to the project**:
   ```bash
   git clone https://github.com/CIRFMF/ksef-client-csharp.git
   cd ksef-client-csharp/KSeF.Client.Tests.CertTestApp
   ```

6. **Run (output to screen, random NIP)**:
   ```bash
   dotnet run --framework net10.0
   ```

7. **Run with parameters**:
   - `--output` - `screen` (default) or `file` (save results to files),
   - `--nip` {nip_number} - e.g., `--nip 8976111986`,
   - optionally: `--no-startup-warnings`.

   ```bash
   dotnet run --framework net10.0 --output file --nip 8976111986 --no-startup-warnings
   ```

---

Related documents:
- [KSeF Authentication](../uwierzytelnianie.md)
- [XAdES Signature](podpis-xades.md)

---
