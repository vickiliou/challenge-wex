package gateway

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vickiliou/challenge-wex/internal/httpresponse"
)

type mockHttpClient struct {
	do func(req *http.Request) (*http.Response, error)
}

func (m *mockHttpClient) Do(req *http.Request) (*http.Response, error) {
	return m.do(req)
}

func TestGetExchangeRate(t *testing.T) {
	testCases := map[string]struct {
		body io.ReadCloser
		want *CurrencyExchangeRate
	}{
		"return only one exchange rate": {
			body: io.NopCloser(bytes.NewBufferString(`{
				"data": [
					{
						"country_currency_desc": "Canada-Dollar",
						"exchange_rate": "1.234"
					}
				]
			}`)),
			want: &CurrencyExchangeRate{
				CountryCurrencyDesc: "Canada-Dollar",
				ExchangeRate:        "1.234",
			},
		},
		"return more than one exchange rate, should return the first one": {
			body: io.NopCloser(bytes.NewBufferString(`{
				"data": [
					{
						"country_currency_desc": "Canada-Dollar",
						"exchange_rate": "1.234"
					},
					{
						"country_currency_desc": "Canada-Dollar",
						"exchange_rate": "1.456"
					}
				]
			}`)),
			want: &CurrencyExchangeRate{
				CountryCurrencyDesc: "Canada-Dollar",
				ExchangeRate:        "1.234",
			},
		},
	}

	for title, tc := range testCases {
		t.Run(title, func(t *testing.T) {
			mockClient := &mockHttpClient{
				do: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       tc.body,
					}, nil
				},
			}

			gw := NewGateway(mockClient)
			input := CurrencyExchangeRateRequest{
				TransactionDate: time.Date(2023, time.September, 21, 0, 0, 0, 0, time.UTC),
				Country:         "Canada",
				Currency:        "Dollar",
			}

			got, gotErr := gw.GetExchangeRate(input)
			assert.NoError(t, gotErr)
			assert.Equal(t, tc.want, got)
		})
	}

}

func TestGetExchangeRate_Error(t *testing.T) {
	someError := errors.New("some error")

	testCases := map[string]struct {
		mockClient     *mockHttpClient
		wantStatusCode int
		wantErr        string
	}{
		"failed to fetch exchange rates": {
			mockClient: &mockHttpClient{
				do: func(req *http.Request) (*http.Response, error) {
					return &http.Response{}, someError
				},
			},
			wantErr: someError.Error(),
		},
		"API request failed": {
			mockClient: &mockHttpClient{
				do: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						Body:       http.NoBody,
						StatusCode: http.StatusBadRequest,
					}, nil
				},
			},
			wantStatusCode: http.StatusBadRequest,
		},
		"failed to decode API response": {
			mockClient: &mockHttpClient{
				do: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						Body:       io.NopCloser(bytes.NewBufferString(`{"]}`)),
						StatusCode: http.StatusOK,
					}, nil
				},
			},
			wantErr: "failed to decode",
		},
		"no return from API": {
			mockClient: &mockHttpClient{
				do: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewBufferString(`{"data":[]}`)),
					}, nil
				},
			},
			wantErr: httpresponse.ErrNoCurrencyConversion.Error(),
		},
	}

	for title, tc := range testCases {
		t.Run(title, func(t *testing.T) {
			gw := NewGateway(tc.mockClient)

			input := CurrencyExchangeRateRequest{
				TransactionDate: time.Date(2023, time.September, 21, 0, 0, 0, 0, time.UTC),
				Country:         "Canada",
				Currency:        "Dollar",
			}

			got, gotErr := gw.GetExchangeRate(input)
			assert.Nil(t, got)
			assert.ErrorContains(t, gotErr, tc.wantErr)
		})
	}
}
