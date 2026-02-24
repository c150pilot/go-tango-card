package tango

import (
	"encoding/json"

	"github.com/go-resty/resty/v2"
)

type Customer struct {
	CustomerIdentifier string `json:"customerIdentifier"`
	DisplayName        string `json:"displayName"`
	Status             string `json:"status"`
	CreatedAt          string `json:"createdAt"`
	Accounts           []UserAccount
}
type UserAccount struct {
	AccountIdentifier string `json:"accountIdentifier"`
	AccountNumber     string `json:"accountNumber"`
	DisplayName       string `json:"displayName"`
	CreatedAt         string `json:"createdAt"`
	Status            string `json:"status"`
}

type CreateCustomerRequest struct {
	CustomerIdentifier string `json:"customerIdentifier"`
	DisplayName        string `json:"displayName"`
}

type CreateCustomerAccountRequest struct {
	AccountIdentifier string `json:"accountIdentifier"`
	DisplayName       string `json:"displayName"`
	ContactEmail      string `json:"contactEmail"`
}

/*
Get a list of all Customers on the Tango Platform.
https://developers.tangocard.com/reference/listcustomers-1
*/
func (c *TangoClient) GetCustomers() ([]Customer, error) {
	url := ApiURL + "/customers"

	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+c.Token).
		Get(url)
	if err != nil {
		return nil, err
	}
	if err := ensureSuccessStatus(resp, "get customers"); err != nil {
		return nil, err
	}

	var responseData []Customer
	err = json.Unmarshal(resp.Body(), &responseData)
	if err != nil {
		return nil, err
	}

	return responseData, nil
}

/*
Get details for a specific Customer on the tango platform.
https://developers.tangocard.com/reference/getcustomer-1
*/
func (c *TangoClient) GetCustomer(customerIdentifier string) (Customer, error) {
	url := ApiURL + "/customers/" + customerIdentifier

	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+c.Token).
		Get(url)
	if err != nil {
		return Customer{}, err
	}
	if err := ensureSuccessStatus(resp, "get customer"); err != nil {
		return Customer{}, err
	}

	var responseData Customer
	err = json.Unmarshal(resp.Body(), &responseData)
	if err != nil {
		return Customer{}, err
	}

	return responseData, nil
}

/*
Get a list of all Accounts created for a specific Customer on the tango platform
https://developers.tangocard.com/reference/listcustomeraccounts-1
*/
func (c *TangoClient) GetCustomerAccounts(customerIdentifier string) ([]UserAccount, error) {
	url := ApiURL + "/customers/" + customerIdentifier + "/accounts"

	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+c.Token).
		Get(url)

	if err != nil {
		return nil, err
	}
	if err := ensureSuccessStatus(resp, "get customer accounts"); err != nil {
		return nil, err
	}

	var responseData []UserAccount
	err = json.Unmarshal(resp.Body(), &responseData)
	if err != nil {
		return nil, err
	}

	return responseData, nil
}

/*
Create a Customer on the Tango Platform
https://developers.tangocard.com/reference/createcustomer-1
*/
func (c *TangoClient) CreateCustomer(customerIdentifier string, displayName string) (CreateCustomerRequest, error) {
	url := ApiURL + "/customers"

	payload := CreateCustomerRequest{
		CustomerIdentifier: customerIdentifier,
		DisplayName:        displayName,
	}

	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+c.Token).
		SetBody(payload).
		Post(url)

	if err != nil {
		return CreateCustomerRequest{}, err
	}
	if err := ensureSuccessStatus(resp, "create customer"); err != nil {
		return CreateCustomerRequest{}, err
	}

	var responseData CreateCustomerRequest
	err = json.Unmarshal(resp.Body(), &responseData)
	if err != nil {
		return CreateCustomerRequest{}, err
	}

	return responseData, nil
}

/*
Create an Account under a specific Customer on this Platform.
https://developers.tangocard.com/reference/createcustomeraccount-1
*/
func (c *TangoClient) CreateCustomerAccount(customerIdentifier string, accountIdentifier string, displayName string, contactEmail string) (CreateCustomerAccountRequest, error) {
	url := ApiURL + "/customers/" + customerIdentifier + "/accounts"

	payload := CreateCustomerAccountRequest{
		AccountIdentifier: accountIdentifier,
		DisplayName:       displayName,
		ContactEmail:      contactEmail,
	}

	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+c.Token).
		SetBody(payload).
		Post(url)

	if err != nil {
		return CreateCustomerAccountRequest{}, err
	}
	if err := ensureSuccessStatus(resp, "create customer account"); err != nil {
		return CreateCustomerAccountRequest{}, err
	}

	var responseData CreateCustomerAccountRequest
	err = json.Unmarshal(resp.Body(), &responseData)
	if err != nil {
		return CreateCustomerAccountRequest{}, err
	}

	return responseData, nil
}
