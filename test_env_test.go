//go:build integration
// +build integration

package tango_test

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func loadTestEnv(t *testing.T) {
	t.Helper()

	if err := godotenv.Load(".env"); err != nil {
		t.Skipf("skipping integration test: .env not found (%v)", err)
	}
}

func requireEnv(t *testing.T, key string) string {
	t.Helper()

	value := os.Getenv(key)
	if value == "" {
		t.Fatalf("missing required environment variable: %s", key)
	}
	return value
}
