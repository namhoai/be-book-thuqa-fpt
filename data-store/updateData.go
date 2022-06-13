package data_store

import (
	"fmt"
	"strconv"

	"github.com/library/models"
)

func (ds *DataStore) UpdateNameOfBook(bookID uint, newName string) error {
	book := &models.Book{}
	err := ds.Db.Where("id = ?", bookID).First(book).Error
	if err != nil {
		return err
	}
	err = ds.Db.Model(book).Where("id = ?", bookID).Updates(map[string]interface{}{
		"name": newName,
	}).Error
	return err
}

func (ds *DataStore) UpdateSubjectOfBook(bookID uint, subjectName string) error {
	book := &models.Book{}
	err := ds.Db.Where("id = ?", bookID).First(book).Error
	if err != nil {
		return err
	}
	// fmt.Println("updating subject name of book")
	err = ds.Db.Model(book).Where("id = ?", bookID).Updates(map[string]interface{}{
		"subject": subjectName,
	}).Error
	if err != nil {
		return err
	}
	// fmt.Println("deleting subject_x_book old entry")
	err = ds.Db.Exec(`delete from subject_x_book where book_id = ?`, bookID).Error
	subject := &models.Subject{}
	// fmt.Println("finding subject id")
	if err := ds.Db.Where("name = ?", subjectName).First(subject).Error; err != nil {
		return err
	}
	subjectXBook := &models.SubjectXBook{
		SubjectID: subject.ID,
		BookID:    book.ID,
	}
	// fmt.Println("creating subject book")
	if err := ds.Db.Create(subjectXBook).Error; err != nil {
		fmt.Println("id of subject: ", subject.ID)
		return err
	}
	return err
}

func (ds *DataStore) UpdateAuthorOfBook(bookID uint, authorId string) error {
	book := &models.Book{}
	err := ds.Db.Where("id = ?", bookID).First(book).Error
	if err != nil {
		return err
	}
	// fmt.Println("updating authorId of book")
	err = ds.Db.Model(book).Where("id = ?", bookID).Updates(map[string]interface{}{
		"authorId": authorId,
	}).Error
	if err != nil {
		return err
	}
	// fmt.Println("deleting old book_x_author entry")
	err = ds.Db.Exec(`delete from book_x_author where book_id = ?`, bookID).Error

	authorID, err := strconv.Atoi(authorId)
	bookXAuthor := &models.BookXAuthor{
		BookID:   book.ID,
		AuthorID: uint(authorID),
	}
	// fmt.Println("creating author book entry")
	if err := ds.Db.Create(bookXAuthor).Error; err != nil {
		return err
	}
	return err
}

func (ds *DataStore) UpdateTitleOfBook(bookID uint, title string) error {
	book := &models.Book{}
	err := ds.Db.Where("id = ?", bookID).First(book).Error
	if err != nil {
		return err
	}
	err = ds.Db.Model(book).Where("id = ?", bookID).Updates(map[string]interface{}{
		"title": title,
	}).Error
	return err
}

func (ds *DataStore) UpdateISBNOfBook(bookID uint, isbn string) error {
	book := &models.Book{}
	err := ds.Db.Where("id = ?", bookID).First(book).Error
	if err != nil {
		return err
	}
	err = ds.Db.Model(book).Where("id = ?", bookID).Updates(map[string]interface{}{
		"isbn": isbn,
	}).Error
	return err
}
