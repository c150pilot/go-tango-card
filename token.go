package tango

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
)

/*
GetToken returns a token for the Tango Card API
https://developers.tangocard.com/reference/acquiretoken
*/
func GetToken(clientID, clientSecret, env string) (TokenResponse, error) {
	url := "https://auth.tangocard.com/oauth/token"
	if env == "sandbox" {
		url = "https://sandbox-auth.tangocard.com/oauth/token"
	}

	client := resty.New()

	request := TokenRequest{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Scope:        "raas.all",
		Audience:     "https://api.tangocard.com/",
		GrantType:    "client_credentials",
	}

	// Response Data
	var responseData TokenResponse
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		Post(url)
	if err != nil {
		return responseData, err
	}

	err = json.Unmarshal(resp.Body(), &responseData)
	if err != nil {
		return responseData, err
	}

	if resp.StatusCode() != 200 {
		return responseData, fmt.Errorf(resp.Status())
	}

	return responseData, nil
}

type TokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Scope        string `json:"scope"`
	Audience     string `json:"audience"`
	GrantType    string `json:"grant_type"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}
