package tango

import (
	"encoding/json"

	"github.com/go-resty/resty/v2"
)

type Account struct {
	AccountIdentifier string `json:"accountIdentifier"`
	AccountNumber     string `json:"accountNumber"`
	DisplayName       string `json:"displayName"`
	CurrencyCode      string `json:"currencyCode"`
	CurrentBalance    int    `json:"currentBalance"`
	CreatedAt         string `json:"createdAt"`
	Status            string `json:"status"`
	ContactEmail      string `json:"contactEmail"`
}

/*
Get details for a specific Account on the Tango Platform
https://developers.tangocard.com/reference/getaccount-1
*/
func (c *TangoClient) GetAccountInfo(accountID string) (Account, error) {
	// https://integration-api.tangocard.com/raas/v2/accounts/{accountIdentifier}
	url := ApiURL + "/accounts/" + accountID

	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+c.Token).
		Get(url)

	if err != nil {
		return Account{}, err
	}
	if err := ensureSuccessStatus(resp, "get account info"); err != nil {
		return Account{}, err
	}

	var responseData Account
	err = json.Unmarshal(resp.Body(), &responseData)
	if err != nil {
		return Account{}, err
	}

	return responseData, nil
}
