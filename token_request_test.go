package tango

import "testing"

func TestBuildClientCredentialsTokenRequest(t *testing.T) {
	request := buildClientCredentialsTokenRequest("client-id", "client-secret")

	if request.GrantType != "client_credentials" {
		t.Fatalf("expected client_credentials grant type, got %s", request.GrantType)
	}
	if request.ClientID != "client-id" || request.ClientSecret != "client-secret" {
		t.Fatalf("unexpected client credentials in request")
	}
	if request.Username != "" || request.Password != "" {
		t.Fatalf("expected empty service-account fields for client credentials request")
	}
}

func TestBuildServiceAccountTokenRequest(t *testing.T) {
	request := buildServiceAccountTokenRequest("client-id", "client-secret", "svc-user", "svc-pass")

	if request.GrantType != "password" {
		t.Fatalf("expected password grant type, got %s", request.GrantType)
	}
	if request.Username != "svc-user" || request.Password != "svc-pass" {
		t.Fatalf("unexpected service-account credentials in request")
	}
	if request.ClientID != "client-id" || request.ClientSecret != "client-secret" {
		t.Fatalf("unexpected client credentials in request")
	}
}
