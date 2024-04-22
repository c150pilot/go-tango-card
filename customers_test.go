package tango_test

import (
	"fmt"
	"github.com/c150pilot/go-tango-card"
	"github.com/joho/godotenv"
	"os"
	"testing"
)

func TestTangoClient_GetCustomers(t *testing.T) {
	// Initialize Environment Variables from .env using godotenv
	err := godotenv.Load()
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	// Get Token
	clientID := os.Getenv("TANGO_CLIENT_ID")
	clientSecret := os.Getenv("TANGO_CLIENT_SECRET")
	accountID := os.Getenv("TANGO_ACCOUNT_ID")
	environment := os.Getenv("ENVIRONMENT")

	tokenResp, err := tango.GetToken(clientID, clientSecret, environment)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
		return
	}

	client, err := tango.New(tokenResp.AccessToken, accountID, true, environment)
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
		return
	}

	customers, err := client.GetCustomers()
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
		return
	}

	if len(customers) == 0 {
		t.Errorf("Expected non-empty, got %v", customers)
	}

	for _, customer := range customers {
		if customer.DisplayName == "ArcadeAppsLLC" {
			fmt.Println(customer)
			fmt.Println(customer.CustomerIdentifier)
		}
	}

	t.Logf("Testing for GetCustomers() complete")
}
