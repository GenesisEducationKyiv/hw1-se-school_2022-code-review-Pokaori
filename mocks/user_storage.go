package mocks

import "bitcoin-service/pkg/models"

type UsersStorageMock struct {
	users []models.User
}

func (storage *UsersStorageMock) GetAllUsers() ([]models.User, error) {
	return storage.users, nil
}

func (storage *UsersStorageMock) IsExist(user *models.User) (bool, error) {
	for _, u := range storage.users {
		if u.Email == user.Email {
			return true, nil
		}
	}
	return false, nil
}

func (storage *UsersStorageMock) Save(user *models.User) error {
	storage.users = append(storage.users, *user)
	return nil
}
