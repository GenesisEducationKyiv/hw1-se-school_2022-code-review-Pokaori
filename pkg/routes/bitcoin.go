package routes

import (
	"bitcoin-service/pkg/controllers"

	"github.com/gorilla/mux"
)

var RegisterBitcoinRoutes = func(router *mux.Router) {
	router.HandleFunc("/rate/", controllers.GetRate).Methods("GET")
	router.HandleFunc("/subscribe/", controllers.Subscribe).Methods("POST")
	router.HandleFunc("/sendEmails/", controllers.SendEmails).Methods("POST")
}
