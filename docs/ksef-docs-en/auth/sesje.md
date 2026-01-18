## Authentication Session Management
10.07.2025

### Retrieving List of Active Authentication Sessions

Returns a list of active authentication sessions.

GET [/auth/sessions](https://api-test.ksef.mf.gov.pl/docs/v2/index.html#tag/Aktywne-sesje/paths/~1api~1v2~1auth~1sessions/get)

Example in ```C#```:
[KSeF.Client.Tests.Core/E2E/Authorization/Sessions/SessionE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/Authorization/Sessions/SessionE2ETests.cs)
```csharp
const int pageSize = 20;
string continuationToken = string.Empty;
List<AuthenticationListItem> authenticationListItems = [];

do
{
    AuthenticationListResponse page = await ActiveSessionsClient.GetActiveSessions(accessToken, pageSize, continuationToken, CancellationToken.None);
    continuationToken = page.ContinuationToken;
    if (page.Items != null)
    {
        authenticationListItems.AddRange(page.Items);
    }
}
while (!string.IsNullOrWhiteSpace(continuationToken));
```

Example in ```Java```:
[SessionIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/SessionIntegrationTest.java)

```java
int pageSize = 10;
AuthenticationListResponse activeSessions = createKSeFClient().getActiveSessions(10, null, accessToken);
while (Strings.isNotBlank(activeSessions.getContinuationToken())) {
    activeSessions = createKSeFClient().getActiveSessions(10, activeSessions.getContinuationToken(), accessToken);
}
```

### Revoking Current Session

DELETE [`/auth/sessions/current`](https://api-test.ksef.mf.gov.pl/docs/v2/index.html#tag/Aktywne-sesje/paths/~1api~1v2~1auth~1sessions~1current/delete)

Revokes the session associated with the token used to call this endpoint. After the operation:
- the associated ```refreshToken``` is revoked,
- active ```accessTokens``` remain valid until their expiration time.

Example in ```C#```:
[KSeF.Client.Tests.Core/E2E/Authorization/Sessions/SessionE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/Authorization/Sessions/SessionE2ETests.cs)
```csharp
await ksefClient.RevokeCurrentSessionAsync(token, cancellationToken);
```

Example in ```Java```:
[SessionIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/SessionIntegrationTest.java)

```java
createKSeFClient().revokeCurrentSession(accessToken);
```

### Revoking Selected Session

DELETE [`/auth/sessions/{referenceNumber}`](https://api-test.ksef.mf.gov.pl/docs/v2/index.html#tag/Aktywne-sesje/paths/~1api~1v2~1auth~1sessions~1%7BreferenceNumber%7D/delete)

Revokes the session with the specified reference number. After the operation:
- the associated ```refreshToken``` is revoked,
- active ```accessTokens``` remain valid until their expiration time.

Example in ```C#```:
[KSeF.Client.Tests.Core/E2E/Authorization/Sessions/SessionE2ETests.cs](https://github.com/CIRFMF/ksef-client-csharp/blob/main/KSeF.Client.Tests.Core/E2E/Authorization/Sessions/SessionE2ETests.cs)
```csharp
await ksefClient.RevokeSessionAsync(referenceNumber, accessToken, cancellationToken);
```

Example in ```Java```:
[SessionIntegrationTest.java](https://github.com/CIRFMF/ksef-client-java/blob/main/demo-web-app/src/integrationTest/java/pl/akmf/ksef/sdk/SessionIntegrationTest.java)

```java
createKSeFClient().revokeSession(secondSessionReferenceNumber, firstAccessTokensPair.accessToken());
```
