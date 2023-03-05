package services

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/survivorbat/qq.maarten.dev/server/domain"
	"gorm.io/gorm"
)

// Compile-time interface checks
var _ IQuizService = new(QuizService)

type IQuizService interface {
	GetByCreator(id uuid.UUID) ([]*domain.Quiz, error)
}

type QuizService struct {
	Database *gorm.DB
}

func (c *QuizService) GetByCreator(id uuid.UUID) ([]*domain.Quiz, error) {
	var result []*domain.Quiz
	if err := c.Database.Where("creator_id = ?", id).Find(&result).Error; err != nil {
		logrus.WithError(err).Error("Failed to get by creator")
		return nil, err
	}

	return result, nil
}
