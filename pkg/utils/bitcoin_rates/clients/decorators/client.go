package decorators

import (
	"bitcoin-service/pkg/utils/bitcoin_rates/clients"
	"log"

	"github.com/patrickmn/go-cache"
)

type DecoratornRateClient struct {
	Wrapee      *clients.BitcoinRateClientInterface
	SystemCache cache.Cache
	next        *clients.BitcoinRateClientInterface
}

func (converter *DecoratornRateClient) SetNext(next *clients.BitcoinRateClientInterface) {
	(*converter.Wrapee).SetNext(next)
}

func (converter *DecoratornRateClient) ExchangeRate(currency string) (float64, error) {
	cahce_title := "bitcoin_rate_res"
	if rate, found := converter.SystemCache.Get(cahce_title); found {
		return rate.(float64), nil
	}

	result, err := (*converter.Wrapee).ExchangeRate(currency)
	if err == nil {
		log.Printf("Got bitcoin exchange rate: %f", result)
		converter.SystemCache.Set(cahce_title, result, cache.DefaultExpiration)
	}
	return result, err

}
