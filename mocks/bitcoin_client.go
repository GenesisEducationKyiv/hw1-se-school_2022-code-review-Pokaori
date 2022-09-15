package mocks

import "bitcoin-service/interfaces"

type BitcoinRateClientMock struct {
	Rate float64
}

func (converter *BitcoinRateClientMock) SetNext(next *interfaces.BitcoinRateClientInterface) {
}

func (converter *BitcoinRateClientMock) ExchangeRate(currency string) (float64, error) {
	return converter.Rate, nil
}
