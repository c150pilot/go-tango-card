package tango

import (
	"encoding/json"
	"fmt"
	"time"

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
	DeliveryMethod     string  `json:"deliveryMethod,omitempty"` // "NONE", "EMAIL", "PHONE", "ADDRESS", or "EMBEDDED" (uppercase)
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
	SendEmail          bool    `json:"sendEmail,omitempty"`      // Deprecated: use deliveryMethod instead
	DeliveryMethod     string  `json:"deliveryMethod,omitempty"` // "NONE", "EMAIL", "PHONE", "ADDRESS", or "EMBEDDED" (uppercase)
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

type CreateOrderResponseError struct {
	Timestamp  time.Time `json:"timestamp"`
	RequestId  string    `json:"requestId"`
	Path       string    `json:"path"`
	HttpCode   int       `json:"httpCode"`
	HttpPhrase string    `json:"httpPhrase"`
	Errors     []struct {
		Path         string `json:"path"`
		I18NKey      string `json:"i18nKey,omitempty"`
		Message      string `json:"message"`
		InvalidValue string `json:"invalidValue"`
		Constraint   string `json:"constraint"`
	} `json:"errors"`
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

	// Transfer data to payload
	payload := CreateOrderRequest{
		AccountIdentifier:  c.AccountIdentifier,
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

	// Use DeliveryMethod if provided, otherwise fall back to SendEmail for backward compatibility
	// Note: When deliveryMethod is set, we should NOT include sendEmail as it's deprecated
	if data.DeliveryMethod != "" {
		payload.DeliveryMethod = data.DeliveryMethod
		// Don't set SendEmail when using deliveryMethod
	} else {
		// Legacy behavior: use SendEmail from client config
		payload.SendEmail = c.SendEmail
	}

	// Marshal payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return CreateOrderResponse{}, err
	}

	payloadMap, err := clean(payloadBytes)
	if err != nil {
		return CreateOrderResponse{}, err
	}

	// If recipient[address] is empty, remove it from the payload
	if payloadMap["recipient"] != nil {
		recipient, ok := payloadMap["recipient"].(map[string]interface{})
		if !ok {
			return CreateOrderResponse{}, fmt.Errorf("recipient is not a map[string]interface{}")
		}
		if recipient["address"] != nil {
			address, ok := recipient["address"].(map[string]interface{})
			if !ok {
				return CreateOrderResponse{}, fmt.Errorf("address is not a map[string]interface{}")
			}
			streetLine1, ok := address["streetLine1"].(string)
			if ok && streetLine1 == "" {
				delete(recipient, "address")
			}
		}
	}

	payloadJSON, err := json.Marshal(payloadMap)
	if err != nil {
		return CreateOrderResponse{}, err
	}

	// Create HTTP Post Request with payload
	client := resty.New()

	// POST JSON string
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+c.Token).
		SetBody(payloadJSON).
		Post(url)
	if err != nil {
		return CreateOrderResponse{}, fmt.Errorf("HTTP request failed: %w", err)
	}

	// If status is not 2xx, check for errors
	if err := ensureSuccessStatus(resp, "create order"); err != nil {
		// Try to parse as error response
		var responseError CreateOrderResponseError
		err = json.Unmarshal(resp.Body(), &responseError)
		if err == nil && len(responseError.Errors) > 0 {
			return CreateOrderResponse{}, fmt.Errorf("Tango API error (status %d): %v", resp.StatusCode(), responseError.Errors)
		}

		// If not parseable as structured error, return raw body
		return CreateOrderResponse{}, fmt.Errorf("%w", err)
	}

	// Check JSON response for errors (even on 2xx status)
	var responseError CreateOrderResponseError
	err = json.Unmarshal(resp.Body(), &responseError)
	if err == nil && len(responseError.Errors) > 0 {
		return CreateOrderResponse{}, fmt.Errorf("Tango API error: %v", responseError.Errors)
	}

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

// GetOrder retrieves order details including credentials from Tango API
// https://developers.tangocard.com/reference/get-details-for-a-specific-order
func (c *TangoClient) GetOrder(referenceOrderID string) (CreateOrderResponse, error) {
	if referenceOrderID == "" {
		return CreateOrderResponse{}, fmt.Errorf("referenceOrderID is required")
	}

	url := fmt.Sprintf("%s/orders/%s", ApiURL, referenceOrderID)

	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+c.Token).
		Get(url)

	if err != nil {
		return CreateOrderResponse{}, fmt.Errorf("HTTP request failed: %w", err)
	}

	if err := ensureSuccessStatus(resp, "get order"); err != nil {
		// Try to parse as error response
		var responseError CreateOrderResponseError
		err = json.Unmarshal(resp.Body(), &responseError)
		if err == nil && len(responseError.Errors) > 0 {
			return CreateOrderResponse{}, fmt.Errorf("Tango API error (status %d): %v", resp.StatusCode(), responseError.Errors)
		}
		return CreateOrderResponse{}, fmt.Errorf("%w", err)
	}

	// Check JSON response for errors (even on 2xx status)
	var responseError CreateOrderResponseError
	err = json.Unmarshal(resp.Body(), &responseError)
	if err == nil && len(responseError.Errors) > 0 {
		return CreateOrderResponse{}, fmt.Errorf("Tango API error: %v", responseError.Errors)
	}

	// Parse response
	var responseData CreateOrderResponse
	err = json.Unmarshal(resp.Body(), &responseData)
	if err != nil {
		return responseData, err
	}

	return responseData, nil
}

// ResendOrder resends an order email using Tango's resend API
func (c *TangoClient) ResendOrder(referenceOrderID string) error {
	if referenceOrderID == "" {
		return fmt.Errorf("referenceOrderID is required")
	}

	url := fmt.Sprintf("%s/orders/%s/resends", ApiURL, referenceOrderID)

	// Create HTTP Post Request
	client := resty.New()

	// POST request to resend order
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+c.Token).
		Post(url)

	if err != nil {
		return fmt.Errorf("failed to resend order: %w", err)
	}

	if resp.StatusCode() != 200 && resp.StatusCode() != 201 && resp.StatusCode() != 204 {
		return fmt.Errorf("resend order failed with status %d (%s): %s", resp.StatusCode(), resp.Status(), string(resp.Body()))
	}

	return nil
}

func clean(payload []byte) (map[string]interface{}, error) {
	var dataMap map[string]interface{}
	if err := json.Unmarshal(payload, &dataMap); err != nil {
		return nil, err
	}

	// Remove empty and nil values from nested maps
	// But preserve deliveryMethod even if empty string (it's a valid value)
	for k, v := range dataMap {
		if k == "deliveryMethod" {
			// Keep deliveryMethod field even if empty - it's a valid API field
			continue
		}
		if v == nil || v == "" || v == 0 {
			delete(dataMap, k)
		}
	}

	return dataMap, nil
}
