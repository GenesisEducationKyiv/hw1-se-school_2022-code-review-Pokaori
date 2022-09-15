package clients

import (
	"io/ioutil"
	"net/http"
	"strconv"
)

type CoingateBitcoinRateClient struct {
	Domain string
	next   *BitcoinRateClientInterface
}

func (converter *CoingateBitcoinRateClient) SetNext(next *BitcoinRateClientInterface) {
	converter.next = next
}

func (converter *CoingateBitcoinRateClient) ExchangeRate(currency string) (float64, error) {
	rate, err := converter.getExchangeRate(currency)
	next := converter.next
	if err != nil && next != nil {
		return (*next).ExchangeRate(currency)
	}
	return rate, err
}

func (converter *CoingateBitcoinRateClient) getExchangeRate(currency string) (float64, error) {
	resp, err := http.Get(converter.Domain + "/v2/rates/merchant/BTC/" + currency)
	if err != nil {
		return 0, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(string(body), 64)
}
