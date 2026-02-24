package tango

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
)

var tokenURLResolver = getTokenURL

/*
GetToken returns a token for the Tango Card API
https://developers.tangocard.com/reference/acquiretoken
*/
func GetToken(clientID, clientSecret, env string) (TokenResponse, error) {
	if err := validateTokenInputs(clientID, clientSecret, env); err != nil {
		return TokenResponse{}, err
	}
	request := buildClientCredentialsTokenRequest(clientID, clientSecret)
	return getTokenFromRequest(request, env)
}

// GetTokenWithServiceAccount attempts service-account OAuth first and falls back to
// client-credentials if the service-account call fails.
//
// Return mode semantics:
// - TokenAuthModeServiceAccount: service account request succeeded
// - TokenAuthModeClientCredentials: service account username/password were not provided
// - TokenAuthModeClientCredentialsFallback: service account request failed, fallback succeeded
func GetTokenWithServiceAccount(clientID, clientSecret, serviceAccountUsername, serviceAccountPassword, env string) (TokenResponse, TokenAuthMode, error) {
	if err := validateTokenInputs(clientID, clientSecret, env); err != nil {
		return TokenResponse{}, TokenAuthModeUnknown, err
	}

	if serviceAccountUsername != "" && serviceAccountPassword != "" {
		request := buildServiceAccountTokenRequest(clientID, clientSecret, serviceAccountUsername, serviceAccountPassword)
		responseData, err := getTokenFromRequest(request, env)
		if err == nil {
			return responseData, TokenAuthModeServiceAccount, nil
		}

		fallbackResponse, fallbackErr := GetToken(clientID, clientSecret, env)
		if fallbackErr != nil {
			return TokenResponse{}, TokenAuthModeUnknown, fmt.Errorf("service-account request failed: %w; fallback request failed: %v", err, fallbackErr)
		}
		return fallbackResponse, TokenAuthModeClientCredentialsFallback, nil
	}

	responseData, err := GetToken(clientID, clientSecret, env)
	if err != nil {
		return TokenResponse{}, TokenAuthModeUnknown, err
	}
	return responseData, TokenAuthModeClientCredentials, nil
}

func getTokenFromRequest(request TokenRequest, env string) (TokenResponse, error) {
	url := tokenURLResolver(env)

	client := resty.New()

	var responseData TokenResponse
	formData := map[string]string{
		"client_id":     request.ClientID,
		"client_secret": request.ClientSecret,
		"scope":         request.Scope,
		"audience":      request.Audience,
		"grant_type":    request.GrantType,
	}
	if request.Username != "" {
		formData["username"] = request.Username
	}
	if request.Password != "" {
		formData["password"] = request.Password
	}

	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetFormData(formData).
		Post(url)
	if err != nil {
		return responseData, err
	}

	err = json.Unmarshal(resp.Body(), &responseData)
	if err != nil {
		return responseData, err
	}

	if resp.StatusCode() != 200 {
		return responseData, fmt.Errorf("token request failed with status %d (%s): %s", resp.StatusCode(), resp.Status(), strings.TrimSpace(string(resp.Body())))
	}

	return responseData, nil
}

func getTokenURL(env string) string {
	if env == "sandbox" {
		return "https://sandbox-auth.tangocard.com/oauth/token"
	}
	return "https://auth.tangocard.com/oauth/token"
}

func validateTokenInputs(clientID, clientSecret, env string) error {
	if strings.TrimSpace(clientID) == "" {
		return fmt.Errorf("clientID is required")
	}
	if strings.TrimSpace(clientSecret) == "" {
		return fmt.Errorf("clientSecret is required")
	}
	if strings.TrimSpace(env) == "" {
		return fmt.Errorf("env is required")
	}
	if env != "production" && env != "sandbox" {
		return fmt.Errorf("env must be either production or sandbox")
	}
	return nil
}

func buildClientCredentialsTokenRequest(clientID, clientSecret string) TokenRequest {
	return TokenRequest{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scope:        "raas.all",
		Audience:     "https://api.tangocard.com/",
		GrantType:    "client_credentials",
	}
}

func buildServiceAccountTokenRequest(clientID, clientSecret, serviceAccountUsername, serviceAccountPassword string) TokenRequest {
	return TokenRequest{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scope:        "raas.all",
		Audience:     "https://api.tangocard.com/",
		GrantType:    "password",
		Username:     serviceAccountUsername,
		Password:     serviceAccountPassword,
	}
}

type TokenAuthMode string

const (
	TokenAuthModeUnknown                   TokenAuthMode = "unknown"
	TokenAuthModeServiceAccount            TokenAuthMode = "service_account"
	TokenAuthModeClientCredentials         TokenAuthMode = "client_credentials"
	TokenAuthModeClientCredentialsFallback TokenAuthMode = "client_credentials_fallback"
)

type TokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Scope        string `json:"scope"`
	Audience     string `json:"audience"`
	GrantType    string `json:"grant_type"`
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}
