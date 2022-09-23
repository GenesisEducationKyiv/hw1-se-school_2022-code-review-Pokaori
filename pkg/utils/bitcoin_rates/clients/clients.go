package clients

type BitcoinRateClientInterface interface {
	ExchangeRate(currency string) (float64, error)
	SetNext(*BitcoinRateClientInterface)
}
