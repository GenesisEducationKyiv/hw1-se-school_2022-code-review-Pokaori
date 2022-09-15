package main

import (
	"bitcoin-service/pkg/config"
	"bitcoin-service/pkg/controllers"
	"bitcoin-service/pkg/repositories"
	"bitcoin-service/pkg/routes"
	"bitcoin-service/pkg/utils/bitcoin_rates/clients"
	ratedecorators "bitcoin-service/pkg/utils/bitcoin_rates/clients/decorators"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/patrickmn/go-cache"
)

func main() {
	config.LoadEnv(".env")
	controller := initController()
	r := mux.NewRouter()
	routes.RegisterBitcoinRoutes(r, controller)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe("0.0.0.0:8000", r))
}

func initController() *controllers.BitcoinController {
	storage := &repositories.UserJsonStorage{PathFile: config.Settings.EmailsStoragePath}
	converter := initConverter()
	return controllers.NewBitcoinController(storage, *converter)
}

func initConverter() *clients.BitcoinRateClientInterface {
	var converter clients.BitcoinRateClientInterface = config.Settings.BitcoinRateCreators[config.Settings.CryptoCurrencyProvider].CreateClient()
	mainConverter := &converter

	prevConv := mainConverter
	for k, v := range config.Settings.BitcoinRateCreators {
		if k != config.Settings.CryptoCurrencyProvider {
			client := v.CreateClient()
			(*prevConv).SetNext(&client)
			prevConv = &client
		}
	}

	var decorator clients.BitcoinRateClientInterface = &ratedecorators.DecoratornRateClient{
		Wrapee:      mainConverter,
		SystemCache: *cache.New(5*time.Minute, 5*time.Minute)}
	return &decorator
}
