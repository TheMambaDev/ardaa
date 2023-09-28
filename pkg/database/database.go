package database

import (
	"ardaa/domain"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type database struct {
	dsn string
}

func NewDatabase(dsn string) *database {
	return &database{
		dsn: dsn,
	}
}

func (database *database) ConnectMysql() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(database.dsn), &gorm.Config{})
	if err != nil {
		slog.Error("Database: could not connect to mysql, err: ", err)
		return nil, err
	}

	return db, nil
}

func (database *database) ConnectSqlite() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(database.dsn), &gorm.Config{})
	if err != nil {
		slog.Error("Database: could not connect to sqlite, err: ", err)
		return nil, err
	}

	return db, nil
}

func Automigrate(db *gorm.DB) error {
	// manually enter the models you want to migrate
	return db.AutoMigrate(&domain.User{})
}
