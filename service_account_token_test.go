//go:build integration
// +build integration

package tango_test

import (
	"encoding/json"
	"testing"

	"github.com/go-resty/resty/v2"
)

func TestServiceAccount_TokenEndpointResponse(t *testing.T) {
	loadTestEnv(t)
	clientID := requireEnv(t, "TANGO_CLIENT_ID")
	clientSecret := requireEnv(t, "TANGO_CLIENT_SECRET")
	serviceAccountUsername := requireEnv(t, "TANGO_SERVICE_ACCOUNT_USERNAME")
	serviceAccountPassword := requireEnv(t, "TANGO_SERVICE_ACCOUNT_PASSWORD")
	environment := requireEnv(t, "ENVIRONMENT")

	url := "https://auth.tangocard.com/oauth/token"
	if environment == "sandbox" {
		url = "https://sandbox-auth.tangocard.com/oauth/token"
	}

	form := map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"username":      serviceAccountUsername,
		"password":      serviceAccountPassword,
		"scope":         "raas.all",
		"audience":      "https://api.tangocard.com/",
		"grant_type":    "password",
	}

	resp, err := resty.New().R().
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(form).
		Post(url)
	if err != nil {
		t.Fatalf("token endpoint request failed: %v", err)
	}

	t.Logf("Tango token endpoint status: %d %s", resp.StatusCode(), resp.Status())

	if resp.StatusCode() != 200 {
		t.Logf("Tango token endpoint response body: %s", string(resp.Body()))
		t.Fatalf("expected 200 from Tango token endpoint, got %d", resp.StatusCode())
	}

	var parsed map[string]any
	if err := json.Unmarshal(resp.Body(), &parsed); err != nil {
		t.Fatalf("failed to parse token response JSON: %v", err)
	}
	if parsed["access_token"] == nil || parsed["access_token"] == "" {
		t.Fatalf("token response missing access_token: %v", parsed)
	}

	accessToken, _ := parsed["access_token"].(string)
	tokenPreview := accessToken
	if len(accessToken) > 12 {
		tokenPreview = accessToken[:6] + "..." + accessToken[len(accessToken)-6:]
	}
	t.Logf("Tango token endpoint response: token_type=%v scope=%v expires_in=%v access_token=%s",
		parsed["token_type"], parsed["scope"], parsed["expires_in"], tokenPreview)
}
