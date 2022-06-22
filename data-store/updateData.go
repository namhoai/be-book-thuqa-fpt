package data_store

import (
	"github.com/library/models"
)

func (ds *DataStore) UpdateBook(bookID uint, newTitle, newISBN string,
	newStock uint, newAuthor, newYear string, newEdition uint,
	newCover, newAbstract, newCategory string, newRating uint) error {
	book := &models.Book{}
	err := ds.Db.Where("id = ?", bookID).First(book).Error
	if err != nil {
		return err
	}
	err = ds.Db.Model(book).Where("id = ?", bookID).Updates(map[string]interface{}{
		"name":     newTitle,
		"isbn":     newISBN,
		"stock":    newStock,
		"author":   newAuthor,
		"year":     newYear,
		"edition":  newEdition,
		"cover":    newCover,
		"abstract": newAbstract,
		"category": newCategory,
		"rating":   newRating,
	}).Error
	return err
}
