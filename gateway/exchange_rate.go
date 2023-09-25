package gateway

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/vickiliou/challenge-wex/internal/httpresponse"
)

const (
	baseURL    = "https://api.fiscaldata.treasury.gov/services/api/fiscal_service/"
	endpoint   = "v1/accounting/od/rates_of_exchange"
	fields     = "?fields=country_currency_desc,exchange_rate,record_date"
	sort       = "&sort=-record_date"
	dateFormat = "2006-01-02"
)

// CurrencyExchangeRateRequest represents the request structure for exchange rate.
type CurrencyExchangeRateRequest struct {
	TransactionDate time.Time
	CurrencyCountry string
	Currency        string
}

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
func (g *Gateway) GetExchangeRate(input CurrencyExchangeRateRequest) (*CurrencyExchangeRate, error) {
	url := constructExchangeRateURL(input)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	res, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch exchange rates: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", res.StatusCode)
	}

	var resp CurrencyExchangeRateResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("failed to decode exchange rates response: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, httpresponse.ErrConvertTargetCurrency
	}

	return &resp.Data[0], nil
}

// constructExchangeRateURL constructs the URL for fetching exchange rates based on
// the target country, target currency, and transaction date.
func constructExchangeRateURL(input CurrencyExchangeRateRequest) string {
	sixMonthsAgo := input.TransactionDate.AddDate(0, -6, 0).Format(dateFormat)
	fmtTxnDate := input.TransactionDate.Format(dateFormat)

	filterParam := fmt.Sprintf("&filter=country_currency_desc:eq:%s-%s,record_date:lte:%s,record_date:gte:%s",
		input.CurrencyCountry, input.Currency, fmtTxnDate, sixMonthsAgo)

	url := baseURL + endpoint + fields + filterParam + sort

	return url
}
