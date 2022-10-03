package creators

import (
	"bitcoin-service/interfaces"
	"bitcoin-service/pkg/domain/bit_rates_clients/clients"
)

type BinanceBitcoinRateCreator struct {
	Domain string
}

func (creator *BinanceBitcoinRateCreator) CreateClient() interfaces.BitcoinRateClientInterface {
	return &clients.BinanceBitcoinRateClient{Domain: creator.Domain}
}
