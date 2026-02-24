# go-tango-card

Go client for the Tango RaaS v2 API.

This package provides a lightweight client for common Tango operations:
- OAuth token acquisition
- Catalog retrieval
- Order create/get/resend
- Customers and customer accounts
- Accounts
- Exchange rates
- Line items

## Install

```bash
go get github.com/c150pilot/go-tango-card
```

## Quick start

```go
package main

import (
	"log"

	tango "github.com/c150pilot/go-tango-card"
)

func main() {
	tokenResp, err := tango.GetToken("client-id", "client-secret", "sandbox")
	if err != nil {
		log.Fatal(err)
	}

	client, err := tango.New(tokenResp.AccessToken, "account-id", true, "sandbox")
	if err != nil {
		log.Fatal(err)
	}

	catalog, err := client.GetCatalogItems()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("brands: %d", len(catalog.Brands))
}
```

## Authentication

### Client credentials

Use `GetToken(clientID, clientSecret, env)` for standard client-credentials flow.

### Service account with fallback

Use `GetTokenWithServiceAccount(clientID, clientSecret, username, password, env)` to try service-account auth first and fall back to client credentials if service-account auth fails.

The function returns `TokenAuthMode`:
- `service_account`: service-account auth succeeded
- `client_credentials`: service-account credentials not provided
- `client_credentials_fallback`: service-account attempt failed, fallback succeeded

## Environments

Supported values:
- `sandbox`
- `production`

Invalid values return an error.

## Error behavior

For non-2xx responses, methods return errors including HTTP status and response body.

## Testing

### Unit tests (default)

Runs offline tests only:

```bash
go test ./...
```

### Integration tests (opt-in)

Integration tests are tagged with `integration` and require real Tango credentials.

```bash
go test -tags=integration ./...
```

Required env vars for integration tests:
- `TANGO_CLIENT_ID`
- `TANGO_CLIENT_SECRET`
- `TANGO_ACCOUNT_ID`
- `ENVIRONMENT` (`sandbox` or `production`)

For service-account integration tests:
- `TANGO_SERVICE_ACCOUNT_USERNAME`
- `TANGO_SERVICE_ACCOUNT_PASSWORD`

## Notes

- Keep credentials out of source control.
- `.env` is ignored by default.
- Tango OAuth/service-account behavior depends on your Tango OAuth client configuration.
