package tango_test

import (
	"fmt"
	"os"
	"testing"

	tango "github.com/c150pilot/go-tango-card"

	"github.com/joho/godotenv"
)

func TestTangoClient_GetCatalogItems(t *testing.T) {
	// Initialize Environment Variables from .env using godotenv
	err := godotenv.Load()
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}
	// Setup Passing Test
	token := os.Getenv("TANGO_API_KEY")
	accountID := os.Getenv("TANGO_ACCOUNT_ID")

	client, err := tango.New(token, accountID, true, "sandbox")
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	// Test that Should Pass
	result, err := client.GetCatalogItems()
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	fmt.Println(result)
}
