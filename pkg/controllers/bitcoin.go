package controllers

import (
	"bitcoin-service/interfaces"
	"bitcoin-service/pkg/config"
	"bitcoin-service/pkg/models"
	"fmt"
	"log"
	"net/http"
	"net/mail"
)

type BitcoinController struct {
	storage   interfaces.UsersStorageInterface
	converter interfaces.BitcoinRateClientInterface
	notifier  interfaces.EmailNotifier
}

func NewBitcoinController(storage interfaces.UsersStorageInterface, converter interfaces.BitcoinRateClientInterface, notifier interfaces.EmailNotifier) *BitcoinController {
	return &BitcoinController{storage: storage, converter: converter, notifier: notifier}
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

	res, err := controller.storage.IsExist(models.NewUser(addr.Address))
	if err != nil {
		controller.writeResponse(&w, *NewResponse(http.StatusInternalServerError))
		log.Println(err)
		return
	}
	if res {
		controller.writeResponse(&w, *NewResponseWithBody(http.StatusConflict, []byte("this email already exists")))
		return
	}

	err = controller.storage.Save(models.NewUser(addr.Address))

	var response Response
	if err != nil {
		response = *NewResponse(http.StatusInternalServerError)
		log.Println(err)
	} else {
		response = *NewResponseWithBody(http.StatusOK, []byte(addr.Address))
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

	users, err := controller.storage.GetAllUsers()
	emails := controller.getEmailListFromUsers(users)

	var response Response
	if err != nil {
		log.Println(err)
		response = *NewResponse(http.StatusBadRequest)

	} else {
		go controller.notifier.SendEmails(emails, rate)
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

func (controller *BitcoinController) getEmailListFromUsers(users []models.User) []string {
	var emails []string
	for _, user := range users {
		emails = append(emails, user.Email)
	}
	return emails
}
