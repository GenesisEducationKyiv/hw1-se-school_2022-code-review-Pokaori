package interfaces

type BitcoinRateClientInterface interface {
	ExchangeRate(currency string) (float64, error)
	SetNext(*BitcoinRateClientInterface)
}
