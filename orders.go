package tango

import (
	"encoding/json"

	"github.com/go-resty/resty/v2"
)

type CreateOrderData struct {
	ExternalRefID      string  `json:"externalRefID"`
	CustomerIdentifier string  `json:"customerIdentifier"`
	Utid               string  `json:"utid"`
	Amount             float64 `json:"amount"`
	EmailSubject       string  `json:"emailSubject"`
	Message            string  `json:"message"`
	Etid               string  `json:"etid"`
	Campaign           string  `json:"campaign"`
	Notes              string  `json:"notes"`
	Sender             Sender  `json:"sender"`
	Recipient          Person
}

type CreateOrderRequest struct {
	ExternalRefID      string  `json:"externalRefID"`
	CustomerIdentifier string  `json:"customerIdentifier"`
	AccountIdentifier  string  `json:"accountIdentifier"`
	Utid               string  `json:"utid"`
	Amount             float64 `json:"amount"`
	EmailSubject       string  `json:"emailSubject"`
	Message            string  `json:"message"`
	SendEmail          bool    `json:"sendEmail"`
	Etid               string  `json:"etid"`
	Campaign           string  `json:"campaign"`
	Notes              string  `json:"notes"`
	Sender             Sender  `json:"sender"`
	Recipient          Person  `json:"recipient"`
}

type Sender struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type Address struct {
	StreetLine1     string `json:"streetLine1"`
	StreetLine2     string `json:"streetLine2"`
	City            string `json:"city"`
	StateOrProvince string `json:"stateOrProvince"`
	PostalCode      string `json:"postalCode"`
	Country         string `json:"country"`
}

type CreateOrderResponse struct {
	ReferenceOrderID       string `json:"referenceOrderID"`
	ExternalRefID          string `json:"externalRefID"`
	CustomerIdentifier     string `json:"customerIdentifier"`
	AccountIdentifier      string `json:"accountIdentifier"`
	AmountCharged          Amount `json:"amountCharged"`
	Denomination           Amount `json:"denomination"`
	UTID                   string `json:"utid"`
	RewardName             string `json:"rewardName"`
	Reward                 Reward `json:"reward"`
	Sender                 Person `json:"sender"`
	Recipient              Person `json:"recipient"`
	EmailSubject           string `json:"emailSubject"`
	Message                string `json:"message"`
	SendEmail              bool   `json:"sendEmail"`
	Status                 string `json:"status"`
	Campaign               string `json:"campaign"`
	CreatedAt              string `json:"createdAt"`
	RedemptionInstructions string `json:"redemptionInstructions"`
}

type Amount struct {
	Value        float64 `json:"value"`
	CurrencyCode string  `json:"currencyCode"`
	ExchangeRate float64 `json:"exchangeRate"`
	Fee          float64 `json:"fee"`
	Total        float64 `json:"total"`
}

type Reward struct {
	Credentials            map[string]string `json:"credentials"`
	CredentialList         []CredentialList  `json:"credentialList"`
	RedemptionInstructions string            `json:"redemptionInstructions"`
}

type CredentialList struct {
	Label          string `json:"label"`
	Value          string `json:"value"`
	Type           string `json:"type"`
	CredentialType string `json:"credentialType"`
}

type Person struct {
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Email     string  `json:"email"`
	Address   Address `json:"address"`
}

func (c *TangoClient) Order(data CreateOrderData) (CreateOrderResponse, error) {
	url := ApiURL + "/orders"

	// TO-DO Validate Data

	// Transfer data to payload
	payload := CreateOrderRequest{
		AccountIdentifier:  c.AccountIdentifier,
		SendEmail:          c.SendEmail,
		ExternalRefID:      data.ExternalRefID,
		CustomerIdentifier: data.CustomerIdentifier,
		Utid:               data.Utid,
		Amount:             data.Amount,
		EmailSubject:       data.EmailSubject,
		Message:            data.Message,
		Etid:               data.Etid,
		Campaign:           data.Campaign,
		Notes:              data.Notes,
		Sender:             data.Sender,
		Recipient:          data.Recipient,
	}

	// Create HTTP Post Request with payload
	client := resty.New()

	// POST JSON string
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+c.Token).
		SetBody(payload).
		Post(url)

	if err != nil {
		return CreateOrderResponse{}, err
	}

	// Parse response
	var responseData CreateOrderResponse
	err = json.Unmarshal(resp.Body(), &responseData)
	if err != nil {
		return responseData, err
	}

	return responseData, nil
}
