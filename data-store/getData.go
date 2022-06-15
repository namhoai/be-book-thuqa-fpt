package data_store

import (
	"github.com/library/models"
)

func (ds *DataStore) GetBooks() (*[]models.Book, error) {
	var books []models.Book
	err := ds.Db.Find(&books).Error
	return &books, err
}

func (ds *DataStore) GetBookByID(id uint) (*models.Book, error) {
	books := &models.Book{}
	err := ds.Db.Where("id=?", id).Find(books).Error
	return books, err
}

func (ds *DataStore) GetBooksByTitle(title string) (*[]models.Book, error) {
	var books []models.Book
	err := ds.Db.Where("name=?", title).Find(&books).Error
	return &books, err
}

func (ds *DataStore) GetBooksByISBN(isbn string) (*[]models.Book, error) {
	var books []models.Book
	err := ds.Db.Where("isbn=?", isbn).Find(&books).Error
	return &books, err
}

// func (ds *DataStore) GetBooksByName(name string) (*[]models.Book, error) {
// 	var books []models.Book
// 	query := `select * from book where name like '%` + name + `%'`
// 	err := ds.Db.Raw(query).Scan(&books).Error
// 	return &books, err
// }

func (ds *DataStore) GetBooksByStock(stock uint) (*[]models.Book, error) {
	var books []models.Book
	err := ds.Db.Where("stock=?", stock).Find(&books).Error
	return &books, err
}

func (ds *DataStore) GetBooksByAuthor(author string) (*[]models.Book, error) {
	var books []models.Book
	err := ds.Db.Where("author=?", author).Find(&books).Error
	return &books, err
}

func (ds *DataStore) GetBooksByYear(year string) (*[]models.Book, error) {
	var books []models.Book
	err := ds.Db.Where("year=?", year).Find(&books).Error
	return &books, err
}

func (ds *DataStore) GetBooksByEdition(edition uint) (*[]models.Book, error) {
	var books []models.Book
	err := ds.Db.Where("edition=?", edition).Find(&books).Error
	return &books, err
}

func (ds *DataStore) GetBooksByAvailable(available bool) (*[]models.Book, error) {
	var books []models.Book
	err := ds.Db.Where("available=?", available).Find(&books).Error
	return &books, err
}

func (ds *DataStore) GetUserByEmail(email string) (*models.Account, error) {
	user := &models.Account{}
	err := ds.Db.Where("email=? and account_role='user'", email).Find(user).Error
	return user, err
}

func (ds *DataStore) GetUserByID(id uint) (*models.Account, error) {
	user := &models.Account{}
	err := ds.Db.Where("id = ?", id).Find(user).Error
	return user, err
}

func (ds *DataStore) GetUsers() (*[]models.Account, error) {
	var users []models.Account
	query := `select * from account where account_role="user"`
	err := ds.Db.Raw(query).Scan(&users).Error
	return &users, err
}

func (ds *DataStore) GetAllBooksReturnByUser() (*[]models.StudentReturnBook, error) {
	var books []models.StudentReturnBook
	err := ds.Db.Find(&books).Error
	return &books, err
}

func (ds *DataStore) GetBooksReturnByUser(bookId uint) (*models.StudentReturnBook, error) {
	book := &models.StudentReturnBook{}
	err := ds.Db.Where("id = ?", bookId).Find(book).Error
	return book, err
}
