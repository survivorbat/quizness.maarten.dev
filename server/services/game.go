package services

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/survivorbat/qq.maarten.dev/server/domain"
	"gorm.io/gorm"
)

// Compile-time interface checks
var _ GameService = new(DBGameService)

type GameService interface {
	GetByQuiz(quizID uuid.UUID) ([]*domain.Game, error)
	Create(game *domain.Game) error
}

type DBGameService struct {
	Database *gorm.DB
}

func (g *DBGameService) GetByQuiz(quizId uuid.UUID) ([]*domain.Game, error) {
	var result []*domain.Game

	if err := g.Database.Where("quiz_id = ?", quizId).Find(&result).Error; err != nil {
		logrus.WithError(err).Error("Failed to fetch by quiz")
		return nil, err
	}

	return result, nil
}

func (g *DBGameService) Create(game *domain.Game) error {
	if err := g.Database.Create(game).Error; err != nil {
		logrus.WithError(err).Error("Failed to create")
		return err
	}

	return nil
}
