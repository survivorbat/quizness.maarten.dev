package services

import (
	"github.com/survivorbat/qq.maarten.dev/server/domain"
	"gorm.io/gorm"
	"testing"
)

func autoMigrate(t *testing.T, db *gorm.DB) {
	err := db.AutoMigrate(&domain.Quiz{}, &domain.Creator{}, &domain.MultipleChoiceQuestion{}, &domain.QuestionOption{},
		&domain.Game{}, &domain.Player{}, &domain.GameAnswer{})
	if err != nil {
		t.Fatal(err.Error())
	}
}
