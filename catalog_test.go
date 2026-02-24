//go:build integration
// +build integration

package tango_test

import (
	"fmt"
	"os"
	"testing"

	tango "github.com/c150pilot/go-tango-card"
)

func TestServiceAccount_GetCatalogItems(t *testing.T) {
	loadTestEnv(t)
	accountID := requireEnv(t, "TANGO_ACCOUNT_ID")
	clientID := requireEnv(t, "TANGO_CLIENT_ID")
	clientSecret := requireEnv(t, "TANGO_CLIENT_SECRET")
	environment := requireEnv(t, "ENVIRONMENT")

	serviceAccountUsername := os.Getenv("TANGO_SERVICE_ACCOUNT_USERNAME")
	serviceAccountPassword := os.Getenv("TANGO_SERVICE_ACCOUNT_PASSWORD")

	token, authMode, err := tango.GetTokenWithServiceAccount(
		clientID,
		clientSecret,
		serviceAccountUsername,
		serviceAccountPassword,
		environment,
	)
	if err != nil {
		t.Fatalf("Failed to get token: %v", err)
	}

	t.Logf("Auth Mode Used: %s", authMode)

	client, err := tango.New(token.AccessToken, accountID, true, environment)
	if err != nil {
		t.Fatalf("Failed to create Tango client: %v", err)
	}

	// Pull list of available gift cards (catalog items)
	resp, err := client.GetCatalogItems()
	if err != nil {
		t.Fatalf("Failed to get catalog items: %v", err)
	}

	if len(resp.Brands) == 0 {
		t.Fatalf("Expected at least one brand in catalog, got 0")
	}

	for i, brand := range resp.Brands {
		// Just print the first 3 brands to avoid massive logs
		if i >= 3 {
			break
		}
		fmt.Printf("Brand: %s (%s)\n", brand.BrandName, brand.BrandKey)
		for _, item := range brand.Items {
			fmt.Printf("  Item: %s (%s) - Face Value: %f\n", item.RewardName, item.Utid, item.FaceValue)
		}
	}
	fmt.Printf("... and %d more brands\n", len(resp.Brands)-3)
}
