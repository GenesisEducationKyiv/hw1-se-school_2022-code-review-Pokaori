package decorators

import (
	"bitcoin-service/interfaces"
	"log"

	"github.com/patrickmn/go-cache"
)

type CachingRateClient struct {
	Wrapee      *interfaces.BitcoinRateClientInterface
	SystemCache cache.Cache
	next        *interfaces.BitcoinRateClientInterface
}

func (converter *CachingRateClient) SetNext(next *interfaces.BitcoinRateClientInterface) {
	(*converter.Wrapee).SetNext(next)
}

func (converter *CachingRateClient) ExchangeRate(currency string) (float64, error) {
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
