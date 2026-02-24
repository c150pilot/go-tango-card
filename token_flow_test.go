package tango

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetToken_Validation(t *testing.T) {
	_, err := GetToken("", "secret", "sandbox")
	if err == nil {
		t.Fatalf("expected validation error for missing clientID")
	}

	_, err = GetToken("id", "", "sandbox")
	if err == nil {
		t.Fatalf("expected validation error for missing clientSecret")
	}

	_, err = GetToken("id", "secret", "")
	if err == nil {
		t.Fatalf("expected validation error for missing env")
	}

	_, err = GetToken("id", "secret", "dev")
	if err == nil {
		t.Fatalf("expected validation error for invalid env")
	}
}

func TestGetTokenWithServiceAccount_FallbackMode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatalf("parse form failed: %v", err)
		}

		grantType := r.Form.Get("grant_type")
		switch grantType {
		case "password":
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`{"error":"unauthorized_client","error_description":"Grant type 'password' not allowed for the client."}`))
		case "client_credentials":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"access_token":"fallback-token","scope":"raas.all","expires_in":86400,"token_type":"Bearer"}`))
		default:
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":"unsupported_grant_type"}`))
		}
	}))
	defer server.Close()

	originalResolver := tokenURLResolver
	tokenURLResolver = func(_ string) string { return server.URL }
	defer func() { tokenURLResolver = originalResolver }()

	token, mode, err := GetTokenWithServiceAccount("client-id", "client-secret", "svc-user", "svc-pass", "sandbox")
	if err != nil {
		t.Fatalf("expected fallback success, got error: %v", err)
	}
	if mode != TokenAuthModeClientCredentialsFallback {
		t.Fatalf("expected fallback mode, got %s", mode)
	}
	if token.AccessToken != "fallback-token" {
		t.Fatalf("expected fallback token, got %s", token.AccessToken)
	}
}

func TestGetTokenWithServiceAccount_ServiceAccountMode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatalf("parse form failed: %v", err)
		}
		if r.Form.Get("grant_type") != "password" {
			t.Fatalf("expected password grant, got %s", r.Form.Get("grant_type"))
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"access_token":"service-account-token","scope":"raas.all","expires_in":86400,"token_type":"Bearer"}`))
	}))
	defer server.Close()

	originalResolver := tokenURLResolver
	tokenURLResolver = func(_ string) string { return server.URL }
	defer func() { tokenURLResolver = originalResolver }()

	token, mode, err := GetTokenWithServiceAccount("client-id", "client-secret", "svc-user", "svc-pass", "sandbox")
	if err != nil {
		t.Fatalf("expected service account success, got error: %v", err)
	}
	if mode != TokenAuthModeServiceAccount {
		t.Fatalf("expected service account mode, got %s", mode)
	}
	if token.AccessToken != "service-account-token" {
		t.Fatalf("expected service-account token, got %s", token.AccessToken)
	}
}
