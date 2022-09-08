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

type EmailNotifier interface {
	SendEmails(emails []string)
}

type EmailHandler interface {
	AddEmail(email string) (string, error)
	GetAllEmails() ([]string, error)
}

type BitcoinReader interface {
	ExchangeRate(currency string) (float64, error)
}

type BitcoinController struct {
	storage   EmailHandler
	converter BitcoinReader
}

func NewBitcoinController(storage EmailHandler, converter BitcoinReader) *BitcoinController {
	return &BitcoinController{storage: storage, converter: converter}
}

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

func (controller *BitcoinController) Subscribe(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		controller.writeIncorrectEmailResponse(&w)
		return
	}
	addr, err := mail.ParseAddress(r.Form.Get("email"))
	if err != nil {
		controller.writeIncorrectEmailResponse(&w)
		return
	}

	res, err := controller.storage.AddEmail(addr.Address)

	var response Response
	if errors.Is(err, models.ErrDuplicateEmail) {
		response = *NewResponseWithBody(http.StatusConflict, []byte(err.Error()))
	} else if err != nil {
		response = *NewResponse(http.StatusInternalServerError)
		log.Println(err)
	} else {
		response = *NewResponseWithBody(http.StatusOK, []byte(res))
	}
	controller.writeResponse(&w, response)
}

func (controller *BitcoinController) GetRate(w http.ResponseWriter, r *http.Request) {
	rate, err := controller.converter.ExchangeRate(config.ExchangeRateUAH)
	var response Response
	if err != nil {
		response = *NewResponse(http.StatusInternalServerError)
		log.Println(err)
	} else {
		response = *NewResponseWithBody(http.StatusOK, []byte(fmt.Sprint(rate)))
	}
	controller.writeResponse(&w, response)
}

func (controller *BitcoinController) SendEmails(w http.ResponseWriter, r *http.Request) {
	rate, err := controller.converter.ExchangeRate(config.ExchangeRateUAH)
	if err != nil {
		log.Println(err)
		response := *NewResponse(http.StatusBadRequest)
		controller.writeResponse(&w, response)
		return
	}

	var notifier EmailNotifier = &utils.EmailBTCtoUAHNotifier{
		Host:     config.Settings.EmailHost,
		Port:     config.Settings.EmailPort,
		From:     config.Settings.EmailName,
		Password: config.Settings.EmailPass,
		Rate:     rate,
	}
	emails, err := controller.storage.GetAllEmails()

	var response Response
	if err != nil {
		log.Println(err)
		response = *NewResponse(http.StatusBadRequest)

	} else {
		go notifier.SendEmails(emails)
		response = *NewResponse(http.StatusOK)
	}
	controller.writeResponse(&w, response)
}

func (controller *BitcoinController) writeIncorrectEmailResponse(w *http.ResponseWriter) {
	response := *NewResponseWithBody(http.StatusBadRequest, []byte("Incorrect email"))
	controller.writeResponse(w, response)
}

func (controller *BitcoinController) writeResponse(w *http.ResponseWriter, res Response) {
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
