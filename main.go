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
	controllers.Storage = &models.EmailJsonStorage{PathFile: config.Settings.EmailsStoragePath}
	controllers.Converter = &utils.BitcoinConverterCoingate{Domain: config.BitcoinCoingateDomain}
	r := mux.NewRouter()
	routes.RegisterBitcoinRoutes(r)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe("0.0.0.0:8000", r))
}
