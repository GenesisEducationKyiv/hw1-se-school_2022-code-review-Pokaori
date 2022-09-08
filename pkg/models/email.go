package models

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"sort"
)

var ErrDuplicateEmail = errors.New("this email already exist")

type EmailJsonStorage struct {
	PathFile string
}

func (storage *EmailJsonStorage) GetAllEmails() ([]string, error) {
	jsonFile, err := storage.openOrCreateStorage()
	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result []string
	if len(byteValue) != 0 {
		err = json.Unmarshal([]byte(byteValue), &result)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (storage *EmailJsonStorage) AddEmail(email string) (string, error) {
	emails, err := storage.GetAllEmails()
	if err != nil {
		return "", err
	}

	emails, err = storage.addToSortedSlice(emails, email)
	if err != nil {
		return "", err
	}

	err = storage.saveToJsonFile(emails)
	return email, err
}

func (storage *EmailJsonStorage) addToSortedSlice(emails []string, email string) ([]string, error) {
	index := sort.SearchStrings(emails, email)

	//check emptiness and if index is not bigger than len to undestand if we can acess email by index; Then compare emails
	if len(emails) > 0 && len(emails) > index && emails[index] == email {
		return nil, ErrDuplicateEmail
	}

	// create empty email, so we can add one more
	emails = append(emails, "")

	//copy emails in a way, that empty email will be on right place
	if len(emails) > 1 {
		copy(emails[index+1:], emails[index:])
	}
	emails[index] = email
	return emails, nil
}

func (storage *EmailJsonStorage) saveToJsonFile(data []string) error {
	content, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(storage.PathFile, content, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (storage *EmailJsonStorage) openOrCreateStorage() (*os.File, error) {
	jsonFile, err := os.Open(storage.PathFile)
	if err != nil {
		log.Println(err)
		_, err = os.Create(storage.PathFile)
		if err != nil {
			log.Fatal(err)
		}
		jsonFile, err = os.Open(storage.PathFile)
		if err != nil {
			return nil, err
		}
		log.Println("Created new file storage for emails")
	}
	return jsonFile, err
}
