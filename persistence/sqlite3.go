package persistence

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func New(name string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(name), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}
