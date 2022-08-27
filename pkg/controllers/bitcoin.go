package controllers

import (
	"bitcoin-service/pkg/config"
	"bitcoin-service/pkg/models"
	"bitcoin-service/pkg/utils"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/mail"
)

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

	var storage models.EmailHandler = &models.EmailJsonStorage{PathFile: config.Settings.EmailsStoragePath}
	res, err := storage.AddEmail(addr.Address)

	var response Response
	if errors.Is(err, models.ErrDuplicateEmail) {
		response = *NewResponseWithBody(http.StatusConflict, []byte(err.Error()))
	} else if err != nil {
		response = *NewResponse(http.StatusInternalServerError)
		log.Println(err)
	} else {
		response = formJsonResponsewithStatus(res, http.StatusOK)
	}
	writeResponse(&w, response)
}

func GetRate(w http.ResponseWriter, r *http.Request) {
	var converter utils.BitcoinReader = &utils.BitcoinConverterCoingate{Domain: config.BitcoinCoingateDomain}
	rate, err := converter.ExchangeRate(config.ExchangeRateUAH)
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
	var converter utils.BitcoinReader = &utils.BitcoinConverterCoingate{Domain: config.BitcoinCoingateDomain}
	rate, err := converter.ExchangeRate(config.ExchangeRateUAH)
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
	var storage models.EmailHandler = &models.EmailJsonStorage{PathFile: config.Settings.EmailsStoragePath}
	emails, err := storage.GetAllEmails()
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

func formJsonResponsewithStatus(data string, status int) Response {
	jsonBody, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return *NewResponse(http.StatusInternalServerError)
	}
	return *NewResponseWithBody(http.StatusOK, jsonBody)
}

func writeResponse(w *http.ResponseWriter, res Response) {
	if len(res.body) > 0 {
		_, err := (*w).Write(res.body)
		if err != nil {
			log.Println(err)
			(*w).WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	(*w).WriteHeader(res.status)
}
