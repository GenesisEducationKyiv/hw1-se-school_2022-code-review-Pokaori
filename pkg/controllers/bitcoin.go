package controllers

import (
	"bitcoin-service/pkg/config"
	"bitcoin-service/pkg/models"
	"bitcoin-service/pkg/utils"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/mail"
)

var Storage models.EmailHandler
var Converter utils.BitcoinReader

type Response struct {
	status int
	body   []byte
}

func NewResponse(status int) *Response {
	return &Response{status: status}
}

func NewResponseWithBody(status int, body []byte) *Response {
	return &Response{status: status, body: body}
}

func Subscribe(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		writeIncorrectEmailResponse(&w)
		return
	}
	addr, err := mail.ParseAddress(r.Form.Get("email"))
	if err != nil {
		writeIncorrectEmailResponse(&w)
		return
	}

	res, err := Storage.AddEmail(addr.Address)

	var response Response
	if errors.Is(err, models.ErrDuplicateEmail) {
		response = *NewResponseWithBody(http.StatusConflict, []byte(err.Error()))
	} else if err != nil {
		response = *NewResponse(http.StatusInternalServerError)
		log.Println(err)
	} else {
		response = *NewResponseWithBody(http.StatusOK, []byte(res))
	}
	writeResponse(&w, response)
}

func GetRate(w http.ResponseWriter, r *http.Request) {
	rate, err := Converter.ExchangeRate(config.ExchangeRateUAH)
	var response Response
	if err != nil {
		response = *NewResponse(http.StatusInternalServerError)
		log.Println(err)
	} else {
		response = *NewResponseWithBody(http.StatusOK, []byte(fmt.Sprint(rate)))
	}
	writeResponse(&w, response)
}

func SendEmails(w http.ResponseWriter, r *http.Request) {
	rate, err := Converter.ExchangeRate(config.ExchangeRateUAH)
	if err != nil {
		log.Println(err)
		response := *NewResponse(http.StatusBadRequest)
		writeResponse(&w, response)
		return
	}

	var notifier utils.EmailNotifier = &utils.EmailBTCtoUAHNotifier{
		Host:     config.Settings.EmailHost,
		Port:     config.Settings.EmailPort,
		From:     config.Settings.EmailName,
		Password: config.Settings.EmailPass,
		Rate:     rate,
	}
	emails, err := Storage.GetAllEmails()

	var response Response
	if err != nil {
		log.Println(err)
		response = *NewResponse(http.StatusBadRequest)

	} else {
		go notifier.SendEmails(emails)
		response = *NewResponse(http.StatusOK)
	}
	writeResponse(&w, response)
}

func writeIncorrectEmailResponse(w *http.ResponseWriter) {
	response := *NewResponseWithBody(http.StatusBadRequest, []byte("Incorrect email"))
	writeResponse(w, response)
}

func writeResponse(w *http.ResponseWriter, res Response) {
	(*w).WriteHeader(res.status)
	if len(res.body) > 0 {
		_, err := (*w).Write(res.body)
		if err != nil {
			log.Println(err)
			(*w).WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
