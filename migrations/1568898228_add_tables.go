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
			  PRIMARY KEY (book_id,user_id),
			  FOREIGN KEY (user_id) REFERENCES account (id) ON DELETE CASCADE,
			  FOREIGN KEY (book_id) REFERENCES book (id) ON DELETE CASCADE
			);
			`,
			`
			CREATE TABLE book_queue (
			  book_id bigint(20) NOT NULL,
			  user_id bigint(20) NOT NULL,
			  reserved_date timestamp,
			  return_date timestamp,
			  PRIMARY KEY (book_id,user_id),
			  FOREIGN KEY (user_id) REFERENCES account (id) ON DELETE CASCADE,
			  FOREIGN KEY (book_id) REFERENCES book (id) ON DELETE CASCADE
			);
			`,
			`
			CREATE TABLE student_return_book (
			  book_id bigint(20) NOT NULL,
			  user_id bigint(20) NOT NULL,
			  reserved_date timestamp,
			  return_date timestamp,
			  PRIMARY KEY (book_id,user_id),
			  FOREIGN KEY (user_id) REFERENCES account (id) ON DELETE CASCADE,
			  FOREIGN KEY (book_id) REFERENCES book (id) ON DELETE CASCADE
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
