package data_store

import (
	"github.com/library/models"
)

func (ds *DataStore) CreateUserAccount(acc models.Account) error {
	return ds.Db.Create(&acc).Error
}

func (ds *DataStore) VerifyUser(details models.LoginDetails) (*models.Account, error) {
	account := &models.Account{}
	err := ds.Db.Where("email=? AND account_role=?", details.Email, details.AccountRole).First(account).Error
	return account, err
}

func (ds *DataStore) CreateBook(book models.Book) error {
	err := ds.Db.Create(&book).Error
	if err != nil {
		return err
	}
	return err
}
