package data_store

func (ds *DataStore) DeleteBook(id uint) error {
	err := ds.Db.Exec(`delete from book where id = ?`, id).Error
	return err
}
