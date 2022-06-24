package migrations

import migrate "github.com/rubenv/sql-migrate"

func init() {
	instance.add(&migrate.Migration{
		Id: "1568898228",
		Up: []string{
			`
			CREATE TABLE book_history (
			  book_id bigint(20) NOT NULL,
			  user_id bigint(20) NOT NULL,
			  reserved_date timestamp,
			  return_date timestamp,
			  status varchar(20) NOT NULL,
			);
			`,
			`
			CREATE TABLE book_queue (
			  book_id bigint(20) NOT NULL,
			  user_id bigint(20) NOT NULL,
			  reserved_date timestamp,
			  return_date timestamp,
			);
			`,
			`
			CREATE TABLE student_return_book (
			  book_id bigint(20) NOT NULL,
			  user_id bigint(20) NOT NULL,
			  reserved_date timestamp,
			  return_date timestamp,
			);
			`,
		},
		//language=SQL
		Down: []string{
			`DROP TABLE book_history;`,
			`DROP TABLE book_queue;`,
			`DROP TABLE book_queue;`,
		},
	})
}
