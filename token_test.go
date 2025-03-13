package tango

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"testing"
)

func TestTangoClient_GetToken(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		return
	}

	clientID := os.Getenv("TANGO_CLIENT_ID")
	clientSecret := os.Getenv("TANGO_CLIENT_SECRET")

	result, err := GetToken(clientID, clientSecret, os.Getenv("ENVIRONMENT"))
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	if result.AccessToken == "" {
		t.Errorf("Expected non-empty string, got %v", result)
	}

	fmt.Println(result)
}
