package data_store

func (ds *DataStore) DeleteBook(id uint) error {
	err := ds.Db.Exec(`delete from book where id = ?`, id).Error
	return err
}

func (ds *DataStore) DeleteRecordStudentReturnBook(id uint) error {
	err := ds.Db.Exec(`delete from student_return_book where book_id = ?`, id).Error
	return err
}
