package clients

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBitcoinConverterReturnCorrectExchangeRate(t *testing.T) {
	var test_rate float64 = 1.4
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("%v", test_rate)))
	}))
	defer server.Close()
	converter := &CoingateBitcoinRateClient{Domain: server.URL}

	resp, err := converter.ExchangeRate("UAH")

	if err != nil {
		t.Fatalf("Got unexpected error:\n %v", err)
	} else if resp != test_rate {
		t.Fatalf("Got incorrect value. Expected %f, got %f", test_rate, resp)
	}

}

func TestBitcoinConverterExchangeRateUseCorrectCurrency(t *testing.T) {
	var test_rate float64 = 1.4
	currency := "UAH"
	path := "/v2/rates/merchant/BTC/" + currency
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != path {
			t.Fatalf("Incorrect path. Expected %s, got %s", path, r.URL.Path)
		}
		w.Write([]byte(fmt.Sprintf("%v", test_rate)))
	}))
	defer server.Close()
	converter := &CoingateBitcoinRateClient{Domain: server.URL}

	_, err := converter.ExchangeRate(currency)

	if err != nil {
		t.Fatalf("Got unexpected error:\n %v", err)
	}
}
