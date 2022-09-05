package models

import (
	"bitcoin-service/pkg/config"
	"os"
	"reflect"
	"sort"
	"testing"
)

func TestEmailHandlerGetEmailsFromEmptyFile(t *testing.T) {
	dir := config.SetupTempDir(t)
	defer os.RemoveAll(dir)
	storage := EmailJsonStorage{PathFile: dir + "/data.json"}

	emails, err := storage.GetAllEmails()

	if err != nil {
		t.Fatalf("Got unexpected error:\n %v", err)
	}
	if len(emails) != 0 {
		t.Fatalf("Got unexpected result from file. Expected [], got %v", emails)
	}
}

func TestEmailHandlerAddEmail(t *testing.T) {
	dir := config.SetupTempDir(t)
	test_email := "example@example.com"
	defer os.RemoveAll(dir)
	storage := EmailJsonStorage{PathFile: dir + "/data.json"}

	email, err := storage.AddEmail(test_email)

	if err != nil {
		t.Fatalf("Got unexpected error:\n %v", err)
	}
	if email != test_email {
		t.Fatalf("Got unexpected result email. Expected %v, got %v", test_email, email)
	}

}

func TestEmailHandlerSortingAddEmail(t *testing.T) {
	dir := config.SetupTempDir(t)
	emails := []string{"b@example.com", "c@example.com", "a@example.com"}
	defer os.RemoveAll(dir)
	storage := EmailJsonStorage{PathFile: dir + "/data.json"}

	for _, email := range emails {
		_, err := storage.AddEmail(email)
		if err != nil {
			t.Fatalf("Got unexpected error while adding:\n %v", err)
		}
	}
	result_emails, err := storage.GetAllEmails()
	if err != nil {
		t.Fatalf("Got unexpected error:\n %v", err)
	}
	sort.Strings(emails)

	if !reflect.DeepEqual(emails, result_emails) {
		t.Fatalf("Got unexpected result email. Expected %v, got %v", emails, result_emails)
	}

}

func TestEmailHandlerAddEmailDuplicate(t *testing.T) {
	dir := config.SetupTempDir(t)
	email := "example@example.com"
	defer os.RemoveAll(dir)
	storage := EmailJsonStorage{PathFile: dir + "/data.json"}

	_, err := storage.AddEmail(email)
	if err != nil {
		t.Fatalf("Got unexpected error while adding:\n %v", err)
	}
	_, err = storage.AddEmail(email)

	if err != ErrDuplicateEmail {
		t.Fatalf("Haven't recieved Duplicate Email Error. Got %v", err)
	}
}
