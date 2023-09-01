package tango

import "fmt"

var ApiURL = "https://integration-api.tangocard.com/raas/v2"

var TangoClientInstance *TangoClient

type TangoClient struct {
	Environment       string
	Token             string
	SendEmail         bool
	AccountIdentifier string
}

func New(token string, accountIdentifier string, sendEmail bool, env string) (*TangoClient, error) {
	// Validate Inputs
	if token == "" {
		return nil, fmt.Errorf("token is required")
	}

	if accountIdentifier == "" {
		return nil, fmt.Errorf("accountIdentifier is required")
	}

	if env == "" || (env != "production" && env != "sandbox") {
		env = "production"
	}

	// Set Proper API URL
	if env == "production" {
		ApiURL = "https://api.tangocard.com/raas/v2"
	} else {
		ApiURL = "https://integration-api.tangocard.com/raas/v2"
	}

	// Return new Client
	TangoClientInstance = &TangoClient{
		Environment:       env,
		Token:             token,
		SendEmail:         sendEmail,
		AccountIdentifier: accountIdentifier,
	}

	return TangoClientInstance, nil
}
