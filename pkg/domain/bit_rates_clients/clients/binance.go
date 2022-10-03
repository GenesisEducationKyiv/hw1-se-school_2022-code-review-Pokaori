package clients

import (
	"bitcoin-service/interfaces"
	"encoding/json"
	"net/http"
	"strconv"
)

type decodedBinanceResponse struct {
	Price string `json:"price"`
}

type BinanceBitcoinRateClient struct {
	Domain string
	next   *interfaces.BitcoinRateClientInterface
}

func (converter *BinanceBitcoinRateClient) SetNext(next *interfaces.BitcoinRateClientInterface) {
	converter.next = next
}

func (converter *BinanceBitcoinRateClient) ExchangeRate(currency string) (float64, error) {
	rate, err := converter.getExchangeRate(currency)
	next := converter.next
	if err != nil && next != nil {
		return (*next).ExchangeRate(currency)
	}
	return rate, err
}

func (converter *BinanceBitcoinRateClient) getExchangeRate(currency string) (float64, error) {
	resp, err := http.Get(converter.Domain + "api/v3/ticker/price?symbol=BTC" + currency)
	if err != nil {
		return 0, err
	}

	var cryptoRate decodedBinanceResponse
	if err := json.NewDecoder(resp.Body).Decode(&cryptoRate); err != nil {
		return 0, err
	}
	return strconv.ParseFloat(string(cryptoRate.Price), 64)
}
