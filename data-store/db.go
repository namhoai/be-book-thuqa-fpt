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
	BookIssue
	DeleteData
	UpdateData
	VerifyUser(models.LoginDetails) (*models.Account, error)
}

type InsertData interface {
	CreateUserAccount(models.Account) error
	CreateAuthor(author models.Author) error
	CreateBook(book models.Book) error
	CreateSubject(book models.Subject) error
}

type GetData interface {
	GetSubjects() (*[]models.Subject, error)
	GetAuthors() (*[]models.Author, error)
	GetBooks() (*[]models.Book, error)
	GetBooksByName(string) (*[]models.Book, error)
	GetBookByBookTitle(string) (*[]models.Book, error)
	GetBookByBookISBN(string) (*[]models.Book, error)
	GetBookByID(uint) (*models.Book, error)
	GetBooksByAuthor(uint) (*[]models.Book, error)
	GetBooksBySubject(uint) (*[]models.Book, error)
	GetAuthorsByName(string) (*[]models.Author, error)
	GetAuthorByID(uint) (*models.Author, error)
	GetUserByName(string) (*[]models.Account, error)
	GetUserByEmail(string) (*[]models.Account, error)
	GetUserByID(uint) (*models.Account, error)
	GetUsers() (*[]models.Account, error)
}

type BookIssue interface {
	GetHistory(uint) (*[]models.BookHistory, error)
	GetCompleteHistory() (*[]models.BookHistory, error)
	CheckAvailability(uint) (bool, error)
	IssueBook(uint, uint) error
	ReturnBook(uint) error
}

type DeleteData interface {
	DeleteBook(uint) error
}

type UpdateData interface {
	UpdateNameOfBook(uint, string) error
	UpdateSubjectOfBook(uint, string) error
	UpdateAuthorOfBook(uint, string) error
	UpdateTitleOfBook(uint, string) error
	UpdateISBNOfBook(uint, string) error
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
