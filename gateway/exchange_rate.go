package gateway

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/vickiliou/challenge-wex/internal/httpresponse"
)

const (
	baseURL  = "https://api.fiscaldata.treasury.gov/services/api/fiscal_service/"
	endpoint = "v1/accounting/od/rates_of_exchange"
	fields   = "?fields=country_currency_desc,exchange_rate,record_date"
	sort     = "&sort=-record_date"
)

// CurrencyExchangeRate represents currency exchange rate data.
type CurrencyExchangeRate struct {
	CountryCurrencyDesc string `json:"country_currency_desc"`
	ExchangeRate        string `json:"exchange_rate"`
}

// CurrencyExchangeRateResponse represents the response structure for exchange rate.
type CurrencyExchangeRateResponse struct {
	Data []CurrencyExchangeRate `json:"data"`
}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Gateway is responsible for fetching exchange rate data.
type Gateway struct {
	httpClient httpClient
}

// NewGateway creates and returns a new instance of the Gateway.
func NewGateway(httpClient httpClient) *Gateway {
	return &Gateway{
		httpClient: httpClient,
	}
}

// GetExchangeRate fetches the exchange rate for a specific date and returns the closest available rate.
func (g *Gateway) GetExchangeRate(txnDate time.Time) (*CurrencyExchangeRate, error) {
	url := constructExchangeRateURL(txnDate)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %v", err.Error())
	}

	res, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch exchange rates: %v", err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", res.StatusCode)
	}

	var resp CurrencyExchangeRateResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode exchange rates response: %v", err)
	}

	if len(resp.Data) == 0 {
		return nil, httpresponse.ErrConvertTargetCurrency
	}

	return &resp.Data[0], nil
}

// constructExchangeRateURL constructs the URL for fetching exchange rates based on
// the target country, target currency, and transaction date.
func constructExchangeRateURL(txnDate time.Time) string {
	targetCountry := "Canada"
	targetCurrency := "Dollar"

	sixMonthsAgo := txnDate.AddDate(0, -6, 0).Format("2006-01-02")
	fmtTxnDate := txnDate.Format("2006-01-02")

	filterParam := fmt.Sprintf("&filter=country_currency_desc:eq:%s-%s,record_date:lte:%s,record_date:gte:%s",
		targetCountry, targetCurrency, fmtTxnDate, sixMonthsAgo)

	url := baseURL + endpoint + fields + filterParam + sort

	return url
}
