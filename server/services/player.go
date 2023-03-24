package services

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/survivorbat/qq.maarten.dev/server/domain"
	"gorm.io/gorm"
)

type PlayerService interface {
	GetByID(playerID uuid.UUID) (*domain.Player, error)
	GetByGame(gameID uuid.UUID) ([]*domain.Player, error)
	Create(player *domain.Player) error
	Delete(player *domain.Player) error
}

type DBPlayerService struct {
	Database *gorm.DB
}

func (g *DBPlayerService) GetByID(playerID uuid.UUID) (*domain.Player, error) {
	var result *domain.Player

	if err := g.Database.First(&result, playerID).Error; err != nil {
		logrus.WithError(err).Error("Failed to fetch by id")
		return nil, err
	}

	return result, nil
}

func (g *DBPlayerService) GetByGame(gameID uuid.UUID) ([]*domain.Player, error) {
	var result []*domain.Player

	if err := g.Database.Preload("Game").Where("game_id = ?", gameID).Find(&result).Error; err != nil {
		logrus.WithError(err).Error("Failed to fetch by game")
		return nil, err
	}

	return result, nil
}

func (g *DBPlayerService) Create(player *domain.Player) error {
	player.GenerateNickname()
	player.GenerateColors()

	if err := g.Database.Create(&player).Error; err != nil {
		logrus.WithError(err).Error("Failed to create")
		return err
	}

	return nil
}

func (g *DBPlayerService) Delete(player *domain.Player) error {
	if err := g.Database.Delete(player).Error; err != nil {
		logrus.WithError(err).Error("Failed to delete")
		return err
	}

	return nil
}
