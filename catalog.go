package tango

import (
	"encoding/json"

	"github.com/go-resty/resty/v2"
)

type Catalog struct {
	CatalogName string  `json:"catalogName"`
	Brands      []Brand `json:"brands"`
}

type Brand struct {
	BrandKey          string            `json:"brandKey"`
	BrandName         string            `json:"brandName"`
	Disclaimer        string            `json:"disclaimer"`
	Description       string            `json:"description"`
	ShortDescription  string            `json:"shortDescription"`
	Terms             string            `json:"terms"`
	CreatedDate       string            `json:"createdDate"`
	LastUpdateDate    string            `json:"lastUpdateDate"`
	BrandRequirements BrandRequirements `json:"brandRequirements"`
	ImageUrls         ImageUrls         `json:"imageUrls"`
	Status            string            `json:"status"`
	Items             []Item            `json:"items"`
}

type BrandRequirements struct {
	DisplayInstructions            string `json:"displayInstructions"`
	TermsAndConditionsInstructions string `json:"termsAndConditionsInstructions"`
	DisclaimerInstructions         string `json:"disclaimerInstructions"`
	AlwaysShowDisclaimer           bool   `json:"alwaysShowDisclaimer"`
}

type ImageUrls struct {
	AdditionalProp string `json:"additionalProp"`
}

type Item struct {
	Utid                       string           `json:"utid"`
	RewardName                 string           `json:"rewardName"`
	CurrencyCode               string           `json:"currencyCode"`
	Status                     string           `json:"status"`
	ValueType                  string           `json:"valueType"`
	RewardType                 string           `json:"rewardType"`
	IsWholeAmountValueRequired bool             `json:"isWholeAmountValueRequired"`
	ExchangeRateRule           string           `json:"exchangeRateRule"`
	MinValue                   int              `json:"minValue"`
	MaxValue                   int              `json:"maxValue"`
	FaceValue                  int              `json:"faceValue"`
	Fee                        Fee              `json:"fee"`
	CreatedDate                string           `json:"createdDate"`
	LastUpdateDate             string           `json:"lastUpdateDate"`
	Countries                  []string         `json:"countries"`
	CredentialTypes            []string         `json:"credentialTypes"`
	RedemptionInstructions     string           `json:"redemptionInstructions"`
	ItemAvailability           ItemAvailability `json:"itemAvailability"`
	FulfillmentType            string           `json:"fulfillmentType"`
}

type Fee struct {
	Type  string `json:"type"`
	Value int    `json:"value"`
}

type ItemAvailability struct {
	ItemAvailabilityStatus string `json:"itemAvailabilityStatus"`
	Note                   string `json:"note"`
	ResolutionDate         string `json:"resolutionDate"`
	StatusPageUrl          string `json:"statusPageUrl"`
	LastModifiedDate       string `json:"lastModifiedDate"`
}

/*
Get all items in your accounts catalog.
https://developers.tangocard.com/reference/getcatalog-1
*/
func (c *TangoClient) GetCatalogItems() (Catalog, error) {
	url := ApiURL + "/catalog"

	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+c.Token).
		Get(url)

	if err != nil {
		return Catalog{}, err
	}

	var responseData Catalog
	err = json.Unmarshal(resp.Body(), &responseData)
	if err != nil {
		return Catalog{}, err
	}

	return responseData, nil
}
