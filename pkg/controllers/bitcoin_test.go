package controllers

import (
	"bitcoin-service/mocks"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestControllerSubscribeIncorrectEmail(t *testing.T) {
	controller := setupController(t, 0.4)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/Subscribe", nil)

	controller.Subscribe(rr, req)

	if rr.Result().StatusCode != http.StatusBadRequest {
		t.Errorf("Status code returned, %d, did not match expected code %d", rr.Result().StatusCode, http.StatusBadRequest)
	}
}

func TestControllerSubscribeSuccessful(t *testing.T) {
	controller := setupController(t, 0.4)
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
	controller := setupController(t, 0.4)
	email := "example@example.com"
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/subscribe/", nil)
	req.Form = url.Values{"email": []string{email}}

	controller.Subscribe(rr, req)
	rr = httptest.NewRecorder()
	controller.Subscribe(rr, req)

	if rr.Result().StatusCode != http.StatusConflict {
		t.Errorf("Status code returned, %d, did not match expected code %d", rr.Result().StatusCode, http.StatusConflict)
	}
}

func TestControllerGetRateSuccessful(t *testing.T) {
	test_rate := 0.4
	controller := setupController(t, test_rate)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/rate/", nil)

	controller.GetRate(rr, req)

	if rr.Result().StatusCode != http.StatusOK {
		t.Errorf("Status code returned, %d, did not match expected code %d", rr.Result().StatusCode, http.StatusOK)
	}
	if string(rr.Body.Bytes()) != fmt.Sprintf("%.1f", test_rate) {
		t.Errorf("Incorrect rate returned. Expected %.1f, got %s", test_rate, string(rr.Body.Bytes()))
	}
}

func TestControllerSendEmailsSuccessful(t *testing.T) {
	controller := setupController(t, 0.4)
	rr := httptest.NewRecorder()

	req := httptest.NewRequest("POST", "/subscribe/", nil)
	controller.Subscribe(rr, req)
	rr = httptest.NewRecorder()
	req = httptest.NewRequest("GET", "/sendEmails/", nil)

	controller.SendEmails(rr, req)

	if rr.Result().StatusCode != http.StatusOK {
		t.Errorf("Status code returned, %d, did not match expected code %d", rr.Result().StatusCode, http.StatusOK)
	}
}

func setupController(t *testing.T, test_rate float64) *BitcoinController {
	storage := &mocks.UsersStorageMock{}
	converter := &mocks.BitcoinRateClientMock{Rate: test_rate}
	notifier := &mocks.EmailNotifierMock{}
	return NewBitcoinController(storage, converter, notifier)
}

func createTestServer(test_rate string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(test_rate))
	}))
}
