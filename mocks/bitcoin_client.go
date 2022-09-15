package mocks

import "bitcoin-service/pkg/utils/bitcoin_rates/clients"

type BitcoinRateClientMock struct {
	Rate float64
}

func (converter *BitcoinRateClientMock) SetNext(next *clients.BitcoinRateClientInterface) {
}

func (converter *BitcoinRateClientMock) ExchangeRate(currency string) (float64, error) {
	return converter.Rate, nil
}
