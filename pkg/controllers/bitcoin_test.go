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
	dir := setupStorage(t)
	defer os.RemoveAll(dir)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/Subscribe", nil)

	Subscribe(rr, req)

	if rr.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("Status code returned, %d, did not match expected code %d", rr.Result().StatusCode, http.StatusBadRequest)
	}
}

func TestControllerSubscribeSuccessful(t *testing.T) {
	dir := setupStorage(t)
	defer os.RemoveAll(dir)
	email := "example@example.com"
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/Subscribe", nil)
	req.Form = url.Values{"email": []string{email}}

	Subscribe(rr, req)

	if rr.Result().StatusCode != http.StatusOK {
		t.Errorf("Status code returned, %d, did not match expected code %d", rr.Result().StatusCode, http.StatusOK)
	}
	if string(rr.Body.Bytes()) != email {
		t.Errorf("Incorrect body returned. Expected %s, got %s", email, string(rr.Body.Bytes()))
	}
}

func TestControllerSubscribeDuplicate(t *testing.T) {
	dir := setupStorage(t)
	defer os.RemoveAll(dir)
	email := "example@example.com"
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/subscribe/", nil)
	req.Form = url.Values{"email": []string{email}}

	Subscribe(rr, req)
	rr = httptest.NewRecorder()
	Subscribe(rr, req)

	if rr.Result().StatusCode != http.StatusConflict {
		t.Errorf("Status code returned, %d, did not match expected code %d", rr.Result().StatusCode, http.StatusOK)
	}
}

func TestControllerGetRateSuccessful(t *testing.T) {
	test_rate := "0.4"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(test_rate))
	}))
	defer server.Close()
	Converter = &utils.BitcoinConverterCoingate{Domain: server.URL}
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/rate/", nil)

	GetRate(rr, req)

	if rr.Result().StatusCode != http.StatusOK {
		t.Errorf("Status code returned, %d, did not match expected code %d", rr.Result().StatusCode, http.StatusOK)
	}
	if string(rr.Body.Bytes()) != test_rate {
		t.Errorf("Incorrect rate returned. Expected %s, got %s", test_rate, string(rr.Body.Bytes()))
	}
}

func TestControllerSendEmailsSuccessful(t *testing.T) {
	Converter = &utils.BitcoinConverterCoingate{Domain: config.BitcoinCoingateDomain}
	dir := setupStorage(t)
	defer os.RemoveAll(dir)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/Subscribe", nil)
	Subscribe(rr, req)
	rr = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/sendEmails/", nil)

	SendEmails(rr, req)

	if rr.Result().StatusCode != http.StatusOK {
		t.Errorf("Status code returned, %d, did not match expected code %d", rr.Result().StatusCode, http.StatusOK)
	}
}

func setupStorage(t *testing.T) string {
	dir := config.SetupTempDir(t)
	Storage = &models.EmailJsonStorage{PathFile: dir + "/data.json"}
	return dir
}
