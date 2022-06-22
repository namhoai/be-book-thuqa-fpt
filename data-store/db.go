package data_store

import (
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/library/envConfig"
	"github.com/library/migrations"
	"github.com/library/models"
	"github.com/sirupsen/logrus"
)

type DataStore struct {
	Db *gorm.DB
}

type DbUtil interface {
	InsertData
	GetData
	BookReserve
	DeleteData
	UpdateData
	VerifyUser(models.LoginDetails) (*models.Account, error)
}

type InsertData interface {
	CreateUserAccount(models.Account) error
	CreateBook(book models.Book) error
}

type GetData interface {
	GetBooks() (*[]models.Book, error)
	GetBooksByTitle(string) (*[]models.Book, error)
	GetBooksByISBN(string) (*[]models.Book, error)
	GetBookByID(uint) (*models.Book, error)
	GetBooksByAuthor(string) (*[]models.Book, error)
	GetBooksByStock(uint) (*[]models.Book, error)
	GetBooksByYear(string) (*[]models.Book, error)
	GetBooksByEdition(uint) (*[]models.Book, error)
	GetBooksByAvailable(bool) (*[]models.Book, error)
	GetUserByEmail(string) (*models.Account, error)
	GetUserByID(uint) (*models.Account, error)
	GetUsers() (*[]models.Account, error)
	GetAllBooksReturnByUser() (*[]models.StudentReturnBook, error)
	GetBooksReturnByUser(uint) (*models.StudentReturnBook, error)
}

type BookReserve interface {
	GetHistory(uint) (*[]models.BookHistory, error)
	GetBooksbyStatus(string) (*[]models.BookHistory, error)
	GetCompleteHistory() (*[]models.BookHistory, error)
	CheckAvailability(uint) (bool, error)
	ReserveBook(uint, uint, *time.Time, *time.Time) error
	AdminConfirmReturnBook(uint) error
	StudentReturnBook(uint, uint, *time.Time, *time.Time) error
	UpdateBookOverdue(*time.Time) error
	GetBooksStudentOverdue(uint, string) (*[]models.BookHistory, error)
	GetBooksStudentReserved(uint) (*[]models.BookHistoryAll, error)
	GetAllBooksStudentReturned() (*[]models.StudentReturnBook, error)
	GetBooksStudentReturned(uint) (*[]models.StudentReturnBook, error)
}

type DeleteData interface {
	DeleteBook(uint) error
	DeleteRecordStudentReturnBook(uint) error
}

type UpdateData interface {
	UpdateBook(uint, string, string, uint, string, string, uint, string, string, string, uint) error
}

var retryAttempts = 0

func DbConnect(dbConfig *envConfig.Env, testing bool) *DataStore {
	var sqlUrl string
	if testing {
		sqlUrl = dbConfig.TestSqlUrl
	} else {
		sqlUrl = dbConfig.SqlUrl
	}
	db, err := gorm.Open(dbConfig.SqlDialect, sqlUrl)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Info("DB connection not established, retrying ...")
		time.Sleep(time.Second * 5)
		retryAttempts++
		if retryAttempts > 5 {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Fatal("DB connection not established")
		}
		return DbConnect(dbConfig, testing)
	} else {
		err = migrations.InitMySQL(db)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Fatal("error running migrations")
		}
		return &DataStore{Db: db}
	}
}
