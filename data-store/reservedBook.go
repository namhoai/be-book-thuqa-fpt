package data_store

import (
	"errors"
	"time"

	"github.com/library/models"
)

func (ds *DataStore) GetCompleteHistory() (*[]models.BookHistory, error) {
	var history []models.BookHistory
	err := ds.Db.Find(&history).Error
	return &history, err
}

func (ds *DataStore) GetHistory(id uint) (*[]models.BookHistory, error) {
	var history []models.BookHistory
	query := `select * from book_history where book_id = ?`
	err := ds.Db.Raw(query, id).Scan(&history).Error
	return &history, err
}

func (ds *DataStore) GetBooksbyStatus(status string) (*[]models.BookHistory, error) {
	var history []models.BookHistory
	query := `select * from book_history where status = ?`
	err := ds.Db.Raw(query, status).Scan(&history).Error
	return &history, err
}

func (ds *DataStore) CheckAvailability(id uint) (bool, error) {
	book := &models.Book{}
	err := ds.Db.Where("id = ?", id).First(book).Error
	if err != nil {
		return false, err
	}
	return book.Available, nil
}

func (ds *DataStore) ReserveBook(bookID, userID uint, reservedDate, returnDate *time.Time) error {
	book := &models.Book{}
	err := ds.Db.Where("id = ?", bookID).First(book).Error
	if err != nil {
		return err
	}
	if book.Available == false {
		queue := &models.BookQueue{
			UserID:       userID,
			BookID:       bookID,
			ReservedDate: reservedDate,
			ReturnDate:   returnDate,
		}
		err = ds.Db.Create(queue).Error
		if err != nil {
			return err
		}
		return errors.New("book unavailable")
	}
	err = ds.Db.Model(book).Where("id = ?", bookID).Updates(map[string]interface{}{
		"available":      false,
		"available_date": returnDate,
	}).Error
	if err != nil {
		return err
	}
	history := &models.BookHistory{
		UserID:       userID,
		BookID:       bookID,
		ReservedDate: reservedDate,
		ReturnDate:   reservedDate,
		Status:       "borrowed",
	}
	return ds.Db.Create(history).Error
}

func (ds *DataStore) AdminConfirmReturnBook(bookID uint) error {
	book := &models.Book{}
	err := ds.Db.Model(book).Where("id = ?", bookID).Updates(map[string]interface{}{
		"available":      true,
		"available_date": time.Now(),
	}).Error
	if err != nil {
		return err
	}
	history := &models.BookHistory{}
	return ds.Db.Model(history).Where("book_id = ?", bookID).Updates(map[string]interface{}{
		"returnDate": time.Now(),
		"status":     "returned",
	}).Error
}

func (ds *DataStore) StudentReturnBook(bookID, userID uint, returnDate *time.Time) error {
	studentReturnBook := &models.StudentReturnBook{}
	err := ds.Db.Model(studentReturnBook).Where("bookId = ? and userId", bookID, userID).Updates(map[string]interface{}{
		"returnDate": returnDate,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func (ds *DataStore) UpdateBookOverdue(currentTime *time.Time) error {
	var books []models.BookHistory
	query := `select * from book where DATE(returnDate) < ?`
	err := ds.Db.Raw(query, currentTime).Scan(&books).Error
	if err != nil {
		return nil
	}
	return ds.Db.Model(books).Updates(map[string]interface{}{
		"status": "overdue",
	}).Error
}

func (ds *DataStore) GetBooksStudentReserved(userID uint, status string) (*[]models.BookHistory, error) {
	var history []models.BookHistory
	query := `select * from book_history where id = ? and status = ?`
	err := ds.Db.Raw(query, userID, status).Scan(&history).Error
	return &history, err
}

func (ds *DataStore) GetAllBooksStudentReturned() (*[]models.StudentReturnBook, error) {
	var books []models.StudentReturnBook
	err := ds.Db.Find(&books).Error
	return &books, err
}

func (ds *DataStore) GetBooksStudentReturned(bookID uint) (*[]models.StudentReturnBook, error) {
	var returnBook []models.StudentReturnBook
	query := `select * from student_return_book where book_id = ?`
	err := ds.Db.Raw(query, bookID).Scan(&returnBook).Error
	return &returnBook, err
}
