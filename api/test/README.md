## How to generate a test certificate and reference XML request:

1. Install .NET SDK version 10
2. Download the official C# client: `git clone https://github.com/CIRFMF/ksef-client-csharp.git`
3. Go to directory: `cd ksef-client-csharp/KSeF.Client.Tests.CertTestApp`
4. Run the C# client: `dotnet run --framework net10.0 --output file`

The command above, when ran with `--output file` option, will generate in the current directory:
- `cert-*.pfx` and `cert-*.cer` files containing the test certificate
- `signed-auth-*.xml` file containing the signed XML

The certificate is a self-signed certificate and is usable only in the test environment.

## Example test certificate

A usable test certificate is stored in `cert-20260102-131809.pfx`. The corresponding context identifier is:

```xml
<Nip>8126178616</Nip>
```
