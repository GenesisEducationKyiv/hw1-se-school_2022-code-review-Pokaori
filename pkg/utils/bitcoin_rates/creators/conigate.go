package creators

import (
	"bitcoin-service/pkg/utils/bitcoin_rates/clients"
)

type CoinbaseBitcoinRateCreator struct {
	Domain string
}

func (creator *CoinbaseBitcoinRateCreator) CreateClient() clients.BitcoinRateClientInterface {
	return &clients.CoinbaseBitcoinRateClient{Domain: creator.Domain}
}
