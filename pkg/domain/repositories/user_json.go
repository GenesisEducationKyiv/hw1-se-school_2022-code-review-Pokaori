package repositories

import (
	"bitcoin-service/pkg/models"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sort"
)

type UserJsonStorage struct {
	PathFile string
}

func (storage *UserJsonStorage) GetAllUsers() ([]models.User, error) {
	jsonFile, err := storage.openOrCreateStorage()
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result []models.User
	if len(byteValue) != 0 {
		err = json.Unmarshal([]byte(byteValue), &result)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (storage *UserJsonStorage) Save(user *models.User) error {
	users, err := storage.GetAllUsers()
	if err != nil {
		return err
	}

	users, err = storage.addToSortedSlice(users, user)
	if err != nil {
		return err
	}

	err = storage.saveToJsonFile(users)
	return err
}

func (storage *UserJsonStorage) addToSortedSlice(users []models.User, user *models.User) ([]models.User, error) {
	index := storage.searchByEmail(users, user.Email)

	// create empty email, so we can add one more
	users = append(users, *models.NewUser(""))

	//copy emails in a way, that empty email will be on the right place
	if len(users) > 1 {
		copy(users[index+1:], users[index:])
	}
	users[index] = *user
	return users, nil
}

func (storage *UserJsonStorage) IsExist(user *models.User) (bool, error) {
	users, err := storage.GetAllUsers()
	if err != nil {
		return false, err
	}
	index := storage.searchByEmail(users, user.Email)

	//check emptiness and if index is not bigger than len to undestand if we can acess email by index; Then compare emails
	return len(users) > 0 && len(users) > index && users[index].Email == user.Email, nil
}

func (storage *UserJsonStorage) saveToJsonFile(data []models.User) error {
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

func (storage *UserJsonStorage) openOrCreateStorage() (*os.File, error) {
	jsonFile, err := os.Open(storage.PathFile)
	if err != nil {
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

func (storage *UserJsonStorage) searchByEmail(users []models.User, email string) int {
	return sort.Search(len(users), func(i int) bool {
		return users[i].Email >= email
	})
}
