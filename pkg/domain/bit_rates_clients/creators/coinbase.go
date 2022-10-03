package creators

import (
	"bitcoin-service/interfaces"
	"bitcoin-service/pkg/domain/bit_rates_clients/clients"
)

type CoingateBitcoinRateCreator struct {
	Domain string
}

func (creator *CoingateBitcoinRateCreator) CreateClient() interfaces.BitcoinRateClientInterface {
	return &clients.CoingateBitcoinRateClient{Domain: creator.Domain}
}
