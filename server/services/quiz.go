package services

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/survivorbat/qq.maarten.dev/server/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Compile-time interface checks
var _ QuizService = new(DBQuizService)

type QuizService interface {
	GetByID(id uuid.UUID) (*domain.Quiz, error)
	GetByCreator(id uuid.UUID) ([]*domain.Quiz, error)
	CreateOrUpdate(quiz *domain.Quiz) error
	Delete(id uuid.UUID) error
}

type DBQuizService struct {
	Database *gorm.DB
}

func (c *DBQuizService) GetByID(id uuid.UUID) (*domain.Quiz, error) {
	var result *domain.Quiz
	if err := c.Database.Preload("Games").Where("id = ?", id).First(&result).Error; err != nil {
		logrus.WithError(err).Error("Failed to get by id")
		return nil, err
	}

	return result, nil
}
func (c *DBQuizService) GetByCreator(id uuid.UUID) ([]*domain.Quiz, error) {
	var result []*domain.Quiz
	if err := c.Database.Preload("MultipleChoiceQuestions.Options").Preload("Games").Where("creator_id = ?", id).Find(&result).Error; err != nil {
		logrus.WithError(err).Error("Failed to get by creator")
		return nil, err
	}

	return result, nil
}

func (c *DBQuizService) CreateOrUpdate(quiz *domain.Quiz) error {
	return c.Database.Transaction(func(tx *gorm.DB) error {
		if err := c.Database.Clauses(clause.OnConflict{UpdateAll: true}).Create(quiz).Error; err != nil {
			logrus.WithError(err).Error("Failed to create quiz")
			return err
		}

		if err := c.Database.Model(quiz).Association("MultipleChoiceQuestions").Replace(quiz.MultipleChoiceQuestions); err != nil {
			return err
		}

		return nil
	})
}

func (c *DBQuizService) Delete(id uuid.UUID) error {
	if err := c.Database.Delete(new(domain.Quiz), id).Error; err != nil {
		logrus.WithError(err).Error("Failed to delete quiz")
		return err
	}

	return nil
}
