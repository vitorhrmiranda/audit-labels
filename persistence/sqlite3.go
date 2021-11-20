package persistence

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func New() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("audit.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}
