package utils

import (
	"io/ioutil"
	"net/http"
	"strconv"
)

type BitcoinReader interface {
	ExchangeRate(currency string) (float64, error)
}

type BitcoinConverterCoingate struct {
	Domain string
}

func (converter *BitcoinConverterCoingate) ExchangeRate(currency string) (float64, error) {
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
