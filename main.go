package main

import (
	"bitcoin-service/pkg/config"
	"bitcoin-service/pkg/controllers"
	"bitcoin-service/pkg/models"
	"bitcoin-service/pkg/routes"
	"bitcoin-service/pkg/utils"
	"log"
	"net/http"

	"github.com/gorilla/mux"
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
	storage := &models.EmailJsonStorage{PathFile: config.Settings.EmailsStoragePath}
	converter := &utils.BitcoinConverterCoingate{Domain: config.BitcoinCoingateDomain}
	return controllers.NewBitcoinController(storage, converter)
}
