package main

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/library/models"
	password_hash "github.com/library/password-hash"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestUser(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "User-Svc Handler Tests")
}

func setupUserData(db *gorm.DB) error {
	adminAccount, err := addAdmin()
	if err != nil {
		return err
	}
	return db.Create(adminAccount).Error
}

func cleanTestData(db *gorm.DB, adminEmail, userEmail string) error {
	if err := db.Exec(`delete from account where email = ?`, adminEmail).Error; err != nil {
		return err
	}
	if err := db.Exec(`delete from account where email = ?`, userEmail).Error; err != nil {
		return err
	}
	return nil
}

func addAdmin() (*models.Account, error) {
	password := "password"
	hashedPwd, err := password_hash.HashPassword(password)
	if err != nil {
		return nil, err
	}
	return &models.Account{
		Email:        "unit@admin.com",
		AccountRole:  models.AdminAccount,
		Password:     password,
		PasswordHash: hashedPwd,
	}, nil
}
