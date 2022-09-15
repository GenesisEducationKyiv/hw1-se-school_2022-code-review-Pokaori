package creators

import (
	"bitcoin-service/interfaces"
	"bitcoin-service/pkg/domain/bit_rates_clients/clients"
)

type CoinbaseBitcoinRateCreator struct {
	Domain string
}

func (creator *CoinbaseBitcoinRateCreator) CreateClient() interfaces.BitcoinRateClientInterface {
	return &clients.CoinbaseBitcoinRateClient{Domain: creator.Domain}
}
