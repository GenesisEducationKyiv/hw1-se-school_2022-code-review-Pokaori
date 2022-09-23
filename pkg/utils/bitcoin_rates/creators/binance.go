package creators

import (
	"bitcoin-service/pkg/utils/bitcoin_rates/clients"
)

type BinanceBitcoinRateCreator struct {
	Domain string
}

func (creator *BinanceBitcoinRateCreator) CreateClient() clients.BitcoinRateClientInterface {
	return &clients.BinanceBitcoinRateClient{Domain: creator.Domain}
}
