package repositories

import (
	"bitcoin-service/pkg/models"
	"io/ioutil"
	"os"
	"sort"
	"testing"
)

func TestEmailHandlerGetEmailsFromEmptyFile(t *testing.T) {
	dir := setupTempDir(t)
	defer os.RemoveAll(dir)
	storage := UserJsonStorage{PathFile: dir + "/data.json"}

	users, err := storage.GetAllUsers()

	if err != nil {
		t.Fatalf("Got unexpected error:\n %v", err)
	}
	if len(users) != 0 {
		t.Fatalf("Got unexpected result from file. Expected [], got %v", users)
	}
}

func TestEmailHandlerAddEmail(t *testing.T) {
	dir := setupTempDir(t)
	test_email := "example@example.com"
	defer os.RemoveAll(dir)
	storage := UserJsonStorage{PathFile: dir + "/data.json"}

	err := storage.Save(models.NewUser(test_email))

	if err != nil {
		t.Fatalf("Got unexpected error:\n %v", err)
	}

}

func TestEmailHandlerSortingAddEmail(t *testing.T) {
	dir := setupTempDir(t)
	emails := []string{"b@example.com", "c@example.com", "a@example.com"}
	defer os.RemoveAll(dir)
	storage := UserJsonStorage{PathFile: dir + "/data.json"}

	for _, email := range emails {
		err := storage.Save(models.NewUser(email))
		if err != nil {
			t.Fatalf("Got unexpected error while adding:\n %v", err)
		}
	}
	result_users, err := storage.GetAllUsers()
	if err != nil {
		t.Fatalf("Got unexpected error:\n %v", err)
	}
	sort.Strings(emails)

	if len(emails) != len(result_users) {
		t.Fatalf("Got unexpected result length. Expected %d, got %d", len(emails), len(result_users))
	}
	for k, email := range emails {
		if email != result_users[k].Email {
			t.Fatalf("Got unexpected result email. Expected %v, got %v", email, result_users[k].Email)
		}
	}

}

func TestEmailHandlerIsExistFalse(t *testing.T) {
	dir := setupTempDir(t)
	email := "example@example.com"
	defer os.RemoveAll(dir)
	storage := UserJsonStorage{PathFile: dir + "/data.json"}

	res, err := storage.IsExist(models.NewUser(email))

	if err != nil {
		t.Fatalf("Got unexpected error while adding:\n %v", err)
	}
	if res {
		t.Fatalf("Email %s exists", email)
	}
}

func TestEmailHandlerIsExistTrue(t *testing.T) {
	dir := setupTempDir(t)
	email := "example@example.com"
	defer os.RemoveAll(dir)
	storage := UserJsonStorage{PathFile: dir + "/data.json"}

	err := storage.Save(models.NewUser(email))
	if err != nil {
		t.Fatalf("Got unexpected error while adding:\n %v", err)
	}
	res, err := storage.IsExist(models.NewUser(email))

	if err != nil {
		t.Fatalf("Got unexpected error while adding:\n %v", err)
	}
	if !res {
		t.Fatalf("Email %s doesn't exist", email)
	}
}

func setupTempDir(t *testing.T) string {
	dir, err := ioutil.TempDir("./", "test")
	if err != nil {
		t.Fatalf("Unable to create Temp Dir : %s", dir)
	}
	return dir
}
