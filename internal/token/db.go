package token

import (
	"errors"
	"net"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Token struct {
	ExpiresAt time.Time `json:"expires_at"`
	ID        string    `json:"id"`
	Token     string    `json:"token"`
	Ip        net.IP    `json:"ip"`
}

var db *gorm.DB

func InitDB() {
	var err error
	db, err = gorm.Open(sqlite.Open("./internal/token/tokendb.sqlite3"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	if db.AutoMigrate(&Token{}) != nil {
		panic("failed to migrate tokens database")
	}
}

func AddToken(token Token) error {
	err := db.Create(&token).Error
	if err != nil {
		return errors.New("failed to create auth session")
	}

	return nil
}

func GetToken(id string) *Token {
	var token Token
	db.Where("token = ?", id).Or("ip = ?", []byte(id)).Find(&token).Scan(&token)

	if token.Token == "" {
		return nil
	}

	return &token
}

func DeleteToken(id string) error {
	return db.Delete(&Token{}, "id = ?", id).Error
}
