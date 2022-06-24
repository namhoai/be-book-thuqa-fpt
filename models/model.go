package models

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	AdminAccount = "admin"
	UserAccount  = "user"
)

type BaseModel struct {
	ID        uint `gorm:"primary_key" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

type Account struct {
	BaseModel
	Email        string `json:"email"`
	AccountRole  string `json:"accountRole"`
	Password     string `gorm:"-" json:"password"`
	Status       string `json:"status"`
	PasswordHash string `json:"-"`
}

func (Account) TableName() string {
	return "account"
}

type Book struct {
	BaseModel
	Name     string `json:"name"`
	ISBN     string `json:"isbn"`
	Stock    uint   `json:"stock"`
	Author   string `json:"author"`
	Year     string `json:"year"`
	Edition  uint   `json:"edition"`
	Cover    string `json:"cover"`
	Abstract string `json:"abstract"`
	Category string `json:"category"`
	Rating   uint   `json:"rating"`
}

func (Book) TableName() string {
	return "book"
}

type BookHistory struct {
	UserID       uint       `json:"userId"`
	BookID       uint       `json:"bookId"`
	ReservedDate *time.Time `json:"reservedDate"`
	ReturnDate   *time.Time `json:"returnDate"`
	Status       string     `json:"status"`
}

type BookHistoryAll struct {
	BookHistory
	Name  string `json:"name"`
	Cover string `json:"cover"`
}

func (BookHistory) TableName() string {
	return "book_history"
}

type BookQueue struct {
	UserID       uint       `json:"userId"`
	BookID       uint       `json:"bookId"`
	ReservedDate *time.Time `json:"reservedDate"`
	ReturnDate   *time.Time `json:"returnDate"`
}

func (BookQueue) TableName() string {
	return "book_queue"
}

type StudentReturnBook struct {
	UserID       uint       `json:"userId"`
	BookID       uint       `json:"bookId"`
	ReservedDate *time.Time `json:"reservedDate"`
	ReturnDate   *time.Time `json:"returnDate"`
}

func (StudentReturnBook) TableName() string {
	return "student_return_book"
}

type Response struct {
	AccountRole string `json:"accountRole"`
	Token       string `json:"token"`
	UserId      uint   `json:"userId"`
}

type LoginDetails struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	AccountRole string `json:"accountRole"`
}

type AuthInfo struct {
	Role string
	jwt.StandardClaims
}
