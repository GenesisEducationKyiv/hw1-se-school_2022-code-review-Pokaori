package main

import (
	"bitcoin-service/interfaces"
	"bitcoin-service/pkg/application/email_notifier"
	"bitcoin-service/pkg/config"
	"bitcoin-service/pkg/controllers"
	ratedecorators "bitcoin-service/pkg/domain/bit_rates_clients/clients/decorators"
	"bitcoin-service/pkg/domain/repositories"
	"bitcoin-service/pkg/routes"
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
	notifier := initEmailNotifier()
	return controllers.NewBitcoinController(storage, *converter, *notifier)
}

func initConverter() *interfaces.BitcoinRateClientInterface {
	var converter interfaces.BitcoinRateClientInterface = config.Settings.BitcoinRateCreators[config.Settings.CryptoCurrencyProvider].CreateClient()
	mainConverter := &converter

	prevConv := mainConverter
	for k, v := range config.Settings.BitcoinRateCreators {
		if k != config.Settings.CryptoCurrencyProvider {
			client := v.CreateClient()
			(*prevConv).SetNext(&client)
			prevConv = &client
		}
	}

	var decorator interfaces.BitcoinRateClientInterface = &ratedecorators.CachingRateClient{
		Wrapee:      mainConverter,
		SystemCache: *cache.New(5*time.Minute, 5*time.Minute)}
	return &decorator
}

func initEmailNotifier() *interfaces.EmailNotifier {
	var notifier interfaces.EmailNotifier = &email_notifier.EmailBTCtoUAHNotifier{
		Host:     config.Settings.EmailHost,
		Port:     config.Settings.EmailPort,
		From:     config.Settings.EmailName,
		Password: config.Settings.EmailPass,
	}
	return &notifier
}
