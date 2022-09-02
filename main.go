package main

import (
	"bitcoin-service/pkg/config"
	"bitcoin-service/pkg/routes"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	config.LoadEnv()
	r := mux.NewRouter()
	routes.RegisterBitcoinRoutes(r)
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe("0.0.0.0:8000", r))
}
