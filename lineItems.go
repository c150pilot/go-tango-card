package tango

import (
	"encoding/json"

	"github.com/go-resty/resty/v2"
)

/*
{
  "keysetPage": {
    "nextPageKeys": [
      "string"
    ],
    "previousPageKeys": [
      "string"
    ],
    "resultCount": 0,
    "totalCount": 0
  },
  "lineItems": [
    {
      "referenceLineItemID": "string",
      "referenceOrderID": "string",
      "orderSource": "string",
      "status": "string",
      "orderStatus": "string",
      "emailStatus": "string",
      "lineNumber": 0,
      "rewardName": "string",
      "amountIssued": {
        "value": 0,
        "currencyCode": "string",
        "exchangeRate": 0,
        "fee": 0,
        "total": 0
      },
      "dateIssued": "2023-09-01T17:39:23.758Z",
      "expirationDate": "2023-09-01T17:39:23.758Z",
      "accountNumber": "string",
      "accountIdentifier": "string",
      "etid": "string",
      "utid": "string",
      "customerIdentifier": "string",
      "recipient": {
        "email": "string",
        "firstName": "string",
        "lastName": "string",
        "address": {
          "streetLine1": "string",
          "streetLine2": "string",
          "city": "string",
          "stateOrProvince": "string",
          "postalCode": "string",
          "country": "string"
        }
      },
      "sender": {
        "firstName": "string",
        "lastName": "string",
        "email": "string"
      }
    }
  ]
}
*/

type LineItemsResponse struct {
	KeysetPage KeysetPage `json:"keysetPage"`
	LineItems  []LineItem `json:"lineItems"`
}

type KeysetPage struct {
	NextPageKeys     []string `json:"nextPageKeys"`
	PreviousPageKeys []string `json:"previousPageKeys"`
	ResultCount      int      `json:"resultCount"`
	TotalCount       int      `json:"totalCount"`
}

type LineItem struct {
	ReferenceLineItemID string `json:"referenceLineItemID"`
	ReferenceOrderID    string `json:"referenceOrderID"`
	OrderSource         string `json:"orderSource"`
	Status              string `json:"status"`
	OrderStatus         string `json:"orderStatus"`
	EmailStatus         string `json:"emailStatus"`
	LineNumber          int    `json:"lineNumber"`
	RewardName          string `json:"rewardName"`
	AmountIssued        Amount `json:"amountIssued"`
	DateIssued          string `json:"dateIssued"`
	ExpirationDate      string `json:"expirationDate"`
	AccountNumber       string `json:"accountNumber"`
	AccountIdentifier   string `json:"accountIdentifier"`
	Etid                string `json:"etid"`
	Utid                string `json:"utid"`
	CustomerIdentifier  string `json:"customerIdentifier"`
	Recipient           Person `json:"recipient"`
	Sender              Sender `json:"sender"`
}

/*
Get a list of Line Items placed under this Platform.
https://developers.tangocard.com/reference/listlineitems
*/
func (c *TangoClient) GetLineItems() (LineItemsResponse, error) {
	url := ApiURL + "/lineItems"

	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+c.Token).
		Get(url)

	if err != nil {
		return LineItemsResponse{}, err
	}
	if err := ensureSuccessStatus(resp, "get line items"); err != nil {
		return LineItemsResponse{}, err
	}

	var responseData LineItemsResponse
	err = json.Unmarshal(resp.Body(), &responseData)
	if err != nil {
		return LineItemsResponse{}, err
	}

	return responseData, nil
}

/*
Get details for a specific line item.
https://developers.tangocard.com/reference/getlineitem
*/
func (c *TangoClient) GetLineItem(lineItemID string) (LineItem, error) {
	url := ApiURL + "/lineItems/" + lineItemID

	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+c.Token).
		Get(url)

	if err != nil {
		return LineItem{}, err
	}
	if err := ensureSuccessStatus(resp, "get line item"); err != nil {
		return LineItem{}, err
	}

	var responseData LineItem
	err = json.Unmarshal(resp.Body(), &responseData)
	if err != nil {
		return LineItem{}, err
	}

	return responseData, nil
}

/*
Resend a specific Line Item.
https://developers.tangocard.com/reference/resendlineitem
*/
func (c *TangoClient) ResendLineItem(lineItemID string) (ResendResponse, error) {
	url := ApiURL + "/lineItems/" + lineItemID + "/resends"

	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+c.Token).
		Post(url)

	if err != nil {
		return ResendResponse{}, err
	}
	if err := ensureSuccessStatus(resp, "resend line item"); err != nil {
		return ResendResponse{}, err
	}

	var responseData ResendResponse
	err = json.Unmarshal(resp.Body(), &responseData)
	if err != nil {
		return ResendResponse{}, err
	}

	return responseData, nil
}

type ResendResponse struct {
	Id        string `json:"id"`
	LegacyID  string `json:"legacyId"`
	CreatedAt string `json:"createdAt"`
	Email     string `json:"email"`
}
