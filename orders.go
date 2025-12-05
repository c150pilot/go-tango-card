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

	// Check Status
	fmt.Printf("[Tango Order] Status: %s\n", resp.Status())
	fmt.Printf("[Tango Order] Request URL: %s\n", url)
	fmt.Printf("[Tango Order] Request Body: %s\n", string(payloadJSON))

	// If status is not 2xx, check for errors
	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		bodyStr := string(resp.Body())
		fmt.Printf("[Tango Order] Error Response Body: %s\n", bodyStr)

		// Try to parse as error response
		var responseError CreateOrderResponseError
		err = json.Unmarshal(resp.Body(), &responseError)
		if err == nil && len(responseError.Errors) > 0 {
			return CreateOrderResponse{}, fmt.Errorf("Tango API error (status %d): %v", resp.StatusCode(), responseError.Errors)
		}

		// If not parseable as structured error, return raw body
		return CreateOrderResponse{}, fmt.Errorf("Tango API error (status %d): %s", resp.StatusCode(), bodyStr)
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

	// Check Status
	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		bodyStr := string(resp.Body())
		// Try to parse as error response
		var responseError CreateOrderResponseError
		err = json.Unmarshal(resp.Body(), &responseError)
		if err == nil && len(responseError.Errors) > 0 {
			return CreateOrderResponse{}, fmt.Errorf("Tango API error (status %d): %v", resp.StatusCode(), responseError.Errors)
		}
		return CreateOrderResponse{}, fmt.Errorf("Tango API error (status %d): %s", resp.StatusCode(), bodyStr)
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
		fmt.Printf("[Tango ResendOrder] ERROR: referenceOrderID is empty\n")
		return fmt.Errorf("referenceOrderID is required")
	}

	url := fmt.Sprintf("%s/orders/%s/resends", ApiURL, referenceOrderID)

	// Log the resend attempt with full details
	fmt.Printf("[Tango ResendOrder] Starting resend process\n")
	fmt.Printf("[Tango ResendOrder] Order ID: %s\n", referenceOrderID)
	fmt.Printf("[Tango ResendOrder] URL: %s\n", url)
	fmt.Printf("[Tango ResendOrder] Environment: %s\n", c.Environment)
	fmt.Printf("[Tango ResendOrder] Account ID: %s\n", c.AccountIdentifier)
	fmt.Printf("[Tango ResendOrder] Token (masked): %s...%s\n", c.Token[:10], c.Token[len(c.Token)-10:])

	// Create HTTP Post Request
	client := resty.New()
	client.SetDebug(false) // Disable resty debug to avoid token exposure

	// POST request to resend order
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+c.Token).
		Post(url)

	if err != nil {
		fmt.Printf("[Tango ResendOrder] HTTP request failed for order %s: %v\n", referenceOrderID, err)
		return fmt.Errorf("failed to resend order: %w", err)
	}

	// Log comprehensive response details
	fmt.Printf("[Tango ResendOrder] === RESPONSE DETAILS ===\n")
	fmt.Printf("[Tango ResendOrder] Status Code: %d\n", resp.StatusCode())
	fmt.Printf("[Tango ResendOrder] Status: %s\n", resp.Status())
	fmt.Printf("[Tango ResendOrder] Headers: %v\n", resp.Header())
	fmt.Printf("[Tango ResendOrder] Body Length: %d bytes\n", len(resp.Body()))
	fmt.Printf("[Tango ResendOrder] Raw Body: %s\n", string(resp.Body()))
	fmt.Printf("[Tango ResendOrder] Response Time: %v\n", resp.Time())

	// Check for success status codes (Tango may return 200, 201, or 204 for success)
	successCodes := []int{200, 201, 204}
	isSuccess := false
	for _, code := range successCodes {
		if resp.StatusCode() == code {
			isSuccess = true
			break
		}
	}

	if !isSuccess {
		fmt.Printf("[Tango ResendOrder] FAILURE: Unexpected status code %d for order %s\n", resp.StatusCode(), referenceOrderID)

		// Try to parse error response if present
		if len(resp.Body()) > 0 {
			var errorResp map[string]interface{}
			if err := json.Unmarshal(resp.Body(), &errorResp); err == nil {
				fmt.Printf("[Tango ResendOrder] Parsed error response: %+v\n", errorResp)
			}
		}

		return fmt.Errorf("resend order failed with status: %s (body: %s)", resp.Status(), string(resp.Body()))
	}

	fmt.Printf("[Tango ResendOrder] SUCCESS: Order %s resent successfully with status %d\n", referenceOrderID, resp.StatusCode())
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
