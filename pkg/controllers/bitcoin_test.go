package controllers

import (
	"bitcoin-service/pkg/config"
	"bitcoin-service/pkg/models"
	"bitcoin-service/pkg/utils"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

func TestControllerSubscribeIncorrectEmail(t *testing.T) {
	controller, dir := setupController(t)
	defer os.RemoveAll(dir)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/Subscribe", nil)

	controller.Subscribe(rr, req)

	if rr.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("Status code returned, %d, did not match expected code %d", rr.Result().StatusCode, http.StatusBadRequest)
	}
}

func TestControllerSubscribeSuccessful(t *testing.T) {
	controller, dir := setupController(t)
	defer os.RemoveAll(dir)
	email := "example@example.com"
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/Subscribe", nil)
	req.Form = url.Values{"email": []string{email}}

	controller.Subscribe(rr, req)

	if rr.Result().StatusCode != http.StatusOK {
		t.Errorf("Status code returned, %d, did not match expected code %d", rr.Result().StatusCode, http.StatusOK)
	}
	if string(rr.Body.Bytes()) != email {
		t.Errorf("Incorrect body returned. Expected %s, got %s", email, string(rr.Body.Bytes()))
	}
}

func TestControllerSubscribeDuplicate(t *testing.T) {
	controller, dir := setupController(t)
	defer os.RemoveAll(dir)
	email := "example@example.com"
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/subscribe/", nil)
	req.Form = url.Values{"email": []string{email}}

	controller.Subscribe(rr, req)
	rr = httptest.NewRecorder()
	controller.Subscribe(rr, req)

	if rr.Result().StatusCode != http.StatusConflict {
		t.Errorf("Status code returned, %d, did not match expected code %d", rr.Result().StatusCode, http.StatusOK)
	}
}

func TestControllerGetRateSuccessful(t *testing.T) {
	controller, dir := setupController(t)
	defer os.RemoveAll(dir)
	test_rate := "0.4"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(test_rate))
	}))
	defer server.Close()
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/rate/", nil)

	controller.GetRate(rr, req)

	if rr.Result().StatusCode != http.StatusOK {
		t.Errorf("Status code returned, %d, did not match expected code %d", rr.Result().StatusCode, http.StatusOK)
	}
	if string(rr.Body.Bytes()) != test_rate {
		t.Errorf("Incorrect rate returned. Expected %s, got %s", test_rate, string(rr.Body.Bytes()))
	}
}

func TestControllerSendEmailsSuccessful(t *testing.T) {
	controller, dir := setupController(t)
	defer os.RemoveAll(dir)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/Subscribe", nil)
	controller.Subscribe(rr, req)
	rr = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/sendEmails/", nil)

	controller.SendEmails(rr, req)

	if rr.Result().StatusCode != http.StatusOK {
		t.Errorf("Status code returned, %d, did not match expected code %d", rr.Result().StatusCode, http.StatusOK)
	}
}

func setupController(t *testing.T) (*BitcoinController, string) {
	dir := config.SetupTempDir(t)
	storage := &models.EmailJsonStorage{PathFile: dir + "/data.json"}
	converter := &utils.BitcoinConverterCoingate{Domain: config.BitcoinCoingateDomain}
	return NewBitcoinController(storage, converter), dir
}
