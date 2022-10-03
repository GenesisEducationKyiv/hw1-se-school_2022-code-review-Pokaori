package config

import (
	"bitcoin-service/interfaces"
	"bitcoin-service/pkg/domain/bit_rates_clients/creators"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type BitcoinRateCreatorInterface interface {
	CreateClient() interfaces.BitcoinRateClientInterface
}
type settings struct {
	EmailName              string
	EmailPass              string
	EmailHost              string
	EmailPort              int
	EmailsStoragePath      string
	CryptoCurrencyProvider string
	BitcoinRateCreators    map[string]BitcoinRateCreatorInterface
}

var Settings settings

func LoadEnv(path string) {
	err := godotenv.Load(path)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	Settings.EmailName = os.Getenv("EMAIL_NAME")
	Settings.EmailPass = os.Getenv("EMAIL_PASSWORD")
	Settings.EmailHost = os.Getenv("EMAIL_HOST")
	Settings.EmailPort, err = strconv.Atoi(os.Getenv("EMAIL_PORT"))
	if err != nil {
		log.Fatal("Incorrect value for EmailPort. Should be int.")
	}
	Settings.EmailsStoragePath = os.Getenv("EMAIL_STORAGE_PATH")
	Settings.CryptoCurrencyProvider = os.Getenv("CRYPTO_CURRENCY_PROVIDER")
	Settings.BitcoinRateCreators = map[string]BitcoinRateCreatorInterface{
		"coingate": &creators.CoingateBitcoinRateCreator{Domain: CoingateDomain},
		"coinbase": &creators.CoinbaseBitcoinRateCreator{Domain: CoinbaseDomain},
		"binance":  &creators.BinanceBitcoinRateCreator{Domain: BinanceDomain},
	}

}
