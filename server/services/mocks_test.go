package services

import (
	"fmt"
	"github.com/survivorbat/qq.maarten.dev/server/domain"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func autoMigrate(t *testing.T, db *gorm.DB) {
	err := db.AutoMigrate(&domain.Quiz{}, &domain.Creator{}, &domain.MultipleChoiceQuestion{}, &domain.QuestionOption{})
	if err != nil {
		t.Fatal(err.Error())
	}
}

func getDb(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())))
	if err != nil {
		t.Fatal(err.Error())
	}

	// Use foreign keys
	if err := db.Exec("PRAGMA foreign_keys = ON;").Error; err != nil {
		t.Error(err.Error())
		return nil
	}

	return db
}
