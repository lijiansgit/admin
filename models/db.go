package models

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/lijiansgit/admin/config"
)

var (
	DB *gorm.DB
)

func Init() (err error) {
	DB, err = gorm.Open("sqlite3", config.Conf.DB.File)
	if err != nil {
		return err
	}

	// defer DB.Close()
	DB.LogMode(true)
	DB.DB().Ping()
	DB.DB().SetMaxIdleConns(10)
	DB.DB().SetMaxOpenConns(100)
	DB.DB().SetConnMaxLifetime(time.Hour)

	// Migrate the schema aaaa
	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&Routes{})

	return nil
}
