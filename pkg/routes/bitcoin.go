package routes

import (
	"bitcoin-service/pkg/controllers"

	"github.com/gorilla/mux"
)

var RegisterBitcoinRoutes = func(router *mux.Router, controller *controllers.BitcoinController) {
	router.HandleFunc("/rate/", controller.GetRate).Methods("GET")
	router.HandleFunc("/subscribe/", controller.Subscribe).Methods("POST")
	router.HandleFunc("/sendEmails/", controller.SendEmails).Methods("POST")
}
