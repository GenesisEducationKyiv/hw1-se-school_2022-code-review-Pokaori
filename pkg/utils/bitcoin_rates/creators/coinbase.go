package creators

import (
	"bitcoin-service/pkg/utils/bitcoin_rates/clients"
)

type CoingateBitcoinRateCreator struct {
	Domain string
}

func (creator *CoingateBitcoinRateCreator) CreateClient() clients.BitcoinRateClientInterface {
	return &clients.CoingateBitcoinRateClient{Domain: creator.Domain}
}
