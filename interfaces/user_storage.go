package interfaces

import "bitcoin-service/pkg/models"

type UsersStorageInterface interface {
	IsExist(user *models.User) (bool, error)
	Save(user *models.User) error
	GetAllUsers() ([]models.User, error)
}
