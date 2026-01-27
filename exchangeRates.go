package tango

import (
	"encoding/json"

	"github.com/go-resty/resty/v2"
)

/*
Get a list of exchange rates.
https://developers.tangocard.com/reference/getexchangerates-1
*/
func (c *TangoClient) GetExchangeRates(baseCurrency, rewardCurrency string) (ExchangeRatesResponse, error) {
	url := ApiURL + "/exchangerates"

	if baseCurrency != "" && rewardCurrency != "" {
		url += "?baseCurrency=" + baseCurrency + "&rewardCurrency=" + rewardCurrency
	} else if baseCurrency != "" {
		url += "?baseCurrency=" + baseCurrency
	} else if rewardCurrency != "" {
		url += "?rewardCurrency=" + rewardCurrency
	}

	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+c.Token).
		Get(url)

	if err != nil {
		return ExchangeRatesResponse{}, err
	}

	var responseData ExchangeRatesResponse
	err = json.Unmarshal(resp.Body(), &responseData)
	if err != nil {
		return ExchangeRatesResponse{}, err
	}

	return responseData, nil

}

type ExchangeRatesRequest struct {
	RewardCurrency string `json:"rewardCurrency"`
	BaseCurrency   string `json:"baseCurrency"`
}

type ExchangeRatesResponse struct {
	Disclaimer    string          `json:"disclaimer"`
	ExchangeRates []ExchangeRates `json:"exchangeRates"`
}

type ExchangeRates struct {
	LastModifiedDate string  `json:"lastModifiedDate"`
	RewardCurrency   string  `json:"rewardCurrency"`
	BaseCurrency     string  `json:"baseCurrency"`
	BaseFx           float64 `json:"baseFx"`
}
