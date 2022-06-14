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
			  issue_date timestamp,
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
		},
	})
}
