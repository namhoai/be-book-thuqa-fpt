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

func (ds *DataStore) GetBooksbyStatus(status string) (*[]models.BookHistoryAll, error) {
	var history []models.BookHistoryAll
	query := `select * from book_history inner join book on book.id=book_history.book_id where book_history.status = ?`
	err := ds.Db.Raw(query, status).Scan(&history).Error
	return &history, err
}

func (ds *DataStore) CheckAvailability(id uint) (bool, error) {
	book := &models.Book{}
	err := ds.Db.Where("id = ?", id).First(book).Error
	if err != nil {
		return false, err
	}
	if book.Stock == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (ds *DataStore) ReserveBook(bookID, userID uint, reservedDate, returnDate *time.Time) error {
	book := &models.Book{}
	err := ds.Db.Where("id = ?", bookID).First(book).Error
	if err != nil {
		return err
	}
	if book.Stock == 0 {
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
	var bookHistory []models.BookHistory
	query := `select * from book_history where user_id = ?`
	err = ds.Db.Raw(query, userID).Scan(&bookHistory).Error
	if err != nil {
		return err
	}
	if len(bookHistory) >= 10 {
		return errors.New("maximum reserved books reached")
	}
	err = ds.Db.Model(book).Where("id = ?", bookID).Updates(map[string]interface{}{
		"stock": book.Stock - 1,
	}).Error
	if err != nil {
		return err
	}
	history := &models.BookHistory{}
	for _, v := range bookHistory {
		if v.BookID == bookID {
			return ds.Db.Model(history).Where("book_id = ? and user_id = ?", bookID, userID).Updates(map[string]interface{}{
				"status": "borrowed",
			}).Error
		}
	}
	history = &models.BookHistory{
		UserID:       userID,
		BookID:       bookID,
		ReservedDate: reservedDate,
		ReturnDate:   returnDate,
		Status:       "borrowed",
	}
	return ds.Db.Create(history).Error
}

func (ds *DataStore) AdminConfirmReturnBook(bookID, studentID uint) error {
	book := &models.Book{}
	err := ds.Db.Model(book).Where("id = ?", bookID).Updates(map[string]interface{}{
		"stock": book.Stock + 1,
	}).Error
	if err != nil {
		return err
	}
	history := &models.BookHistory{}
	return ds.Db.Model(history).Where("book_id = ? and user_id = ?", bookID, studentID).Updates(map[string]interface{}{
		"returnDate": time.Now(),
		"status":     "returned",
	}).Error
}

func (ds *DataStore) StudentReturnBook(bookID, userID uint, reservedDate, returnDate *time.Time) error {
	studentReturnBook := &models.StudentReturnBook{
		UserID:       userID,
		BookID:       bookID,
		ReservedDate: reservedDate,
		ReturnDate:   returnDate,
	}
	// err := ds.Db.Model(studentReturnBook).Where("book_id = ? and user_id = ?", bookID, userID).Updates(map[string]interface{}{
	// 	"returnDate": returnDate,
	// }).Error
	err := ds.Db.Create(studentReturnBook).Error
	if err != nil {
		return err
	}
	return nil
}

func (ds *DataStore) UpdateBookOverdue(currentTime *time.Time) error {
	var books []models.BookHistory
	query := `select * from book_history where DATE(return_date) < ?`
	err := ds.Db.Raw(query, currentTime).Scan(&books).Error
	if err != nil {
		return nil
	}
	return ds.Db.Model(books).Updates(map[string]interface{}{
		"status": "overdue",
	}).Error
}

func (ds *DataStore) GetBooksStudentOverdue(userID uint) (*[]models.BookHistoryAll, error) {
	var history []models.BookHistoryAll
	query := `select * from book_history inner join book on book.id=book_history.book_id where book_history.user_id = ? and book_history.status = 'overdue'`
	err := ds.Db.Raw(query, userID).Scan(&history).Error
	return &history, err
}

func (ds *DataStore) GetBooksStudentReserved(userID uint) (*[]models.BookHistoryAll, error) {
	var history []models.BookHistoryAll
	query := `select * from book_history inner join book on book.id=book_history.book_id where book_history.user_id = ? and book_history.status = 'borrowed'`
	err := ds.Db.Raw(query, userID).Scan(&history).Error
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
