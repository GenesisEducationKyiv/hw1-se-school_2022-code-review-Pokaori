package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type decodedCoinbaseResponse struct {
	Data struct {
		Price string `json:"amount"`
	} `json:"data"`
}

type CoinbaseBitcoinRateClient struct {
	Domain string
	next   *BitcoinRateClientInterface
}

func (converter *CoinbaseBitcoinRateClient) SetNext(next *BitcoinRateClientInterface) {
	converter.next = next
}

func (converter *CoinbaseBitcoinRateClient) ExchangeRate(currency string) (float64, error) {
	rate, err := converter.getExchangeRate(currency)
	next := converter.next
	if err != nil && next != nil {
		return (*next).ExchangeRate(currency)
	}
	return rate, err
}

func (converter *CoinbaseBitcoinRateClient) getExchangeRate(currency string) (float64, error) {
	resp, err := http.Get(fmt.Sprintf(converter.Domain+"/v2/prices/BTC-%s/spot", currency))
	if err != nil {
		return 0, err
	}

	var cryptoRate decodedCoinbaseResponse
	if err := json.NewDecoder(resp.Body).Decode(&cryptoRate); err != nil {
		return 0, err
	}
	return strconv.ParseFloat(string(cryptoRate.Data.Price), 64)
}
