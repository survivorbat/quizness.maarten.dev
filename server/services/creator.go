package services

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/survivorbat/qq.maarten.dev/server/domain"
	"gorm.io/gorm"
)

// Compile-time interface checks
var _ CreatorService = new(DBCreatorService)

type CreatorService interface {
	GetOrCreate(authID string) (*domain.Creator, error)
	GetByID(id uuid.UUID) (*domain.Creator, error)
}

type DBCreatorService struct {
	Database *gorm.DB
}

func (c *DBCreatorService) GetOrCreate(authID string) (*domain.Creator, error) {
	result := &domain.Creator{AuthID: authID}
	result.GenerateNickname()

	if err := c.Database.FirstOrCreate(&result, map[string]any{"auth_id": authID}).Error; err != nil {
		logrus.WithError(err).Error("Failed to get or create")
		return nil, err
	}

	return result, nil
}

func (c *DBCreatorService) GetByID(id uuid.UUID) (*domain.Creator, error) {
	var result *domain.Creator
	if err := c.Database.Find(&result, id).Error; err != nil {
		logrus.WithError(err).Error("Failed to get by ID")
		return nil, err
	}

	return result, nil
}
