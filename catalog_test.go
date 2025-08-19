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
	if err := godotenv.Load(".env"); err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	accountID := os.Getenv("TANGO_ACCOUNT_ID")
	clientID := os.Getenv("TANGO_CLIENT_ID")
	clientSecret := os.Getenv("TANGO_CLIENT_SECRET")

	token, err := tango.GetToken(clientID, clientSecret, os.Getenv("ENVIRONMENT"))
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	client, err := tango.New(token.AccessToken, accountID, true, os.Getenv("ENVIRONMENT"))
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	// Test that Should Pass
	resp, err := client.GetCatalogItems()
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	for _, brand := range resp.Brands {
		fmt.Println("Brand: ", brand.BrandName+" ("+brand.BrandKey+")")
		for _, item := range brand.Items {
			fmt.Println(" Item: ", item.RewardName+" ("+item.Utid+")")
			fmt.Println(" Face Value: ", item.FaceValue)
			if item.MinValue != 0 {
				fmt.Println(" Min Value: ", item.MinValue)
				fmt.Println(" Max Value: ", item.MaxValue)
			}
		}
		fmt.Println("------------------------------------------------")
	}
}
