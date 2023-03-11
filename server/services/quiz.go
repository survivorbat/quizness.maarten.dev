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
	GetByID(id uuid.UUID) (*domain.Quiz, error)
	GetByCreator(id uuid.UUID) ([]*domain.Quiz, error)
	Create(quiz *domain.Quiz) error
	Delete(id uuid.UUID) error
}

type QuizService struct {
	Database *gorm.DB
}

func (c *QuizService) GetByID(id uuid.UUID) (*domain.Quiz, error) {
	var result *domain.Quiz
	if err := c.Database.Find(&result, id).Error; err != nil {
		logrus.WithError(err).Error("Failed to get by id")
		return nil, err
	}

	return result, nil
}
func (c *QuizService) GetByCreator(id uuid.UUID) ([]*domain.Quiz, error) {
	var result []*domain.Quiz
	if err := c.Database.Preload("MultipleChoiceQuestions.Options").Where("creator_id = ?", id).Find(&result).Error; err != nil {
		logrus.WithError(err).Error("Failed to get by creator")
		return nil, err
	}

	return result, nil
}

func (c *QuizService) Create(quiz *domain.Quiz) error {
	if err := c.Database.Create(quiz).Error; err != nil {
		logrus.WithError(err).Error("Failed to create quiz")
		return err
	}

	return nil
}

func (c *QuizService) Delete(id uuid.UUID) error {
	if err := c.Database.Delete(new(domain.Quiz), id).Error; err != nil {
		logrus.WithError(err).Error("Failed to delete quiz")
		return err
	}

	return nil
}
