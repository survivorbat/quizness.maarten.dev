package services

import (
	"github.com/google/uuid"
	"github.com/survivorbat/qq.maarten.dev/server/domain"
	"gorm.io/gorm"
)

// Compile-time interface checks
var _ ICreatorService = new(CreatorService)

type ICreatorService interface {
	GetOrCreate(authID string) (*domain.Creator, error)
	GetByID(id uuid.UUID) (*domain.Creator, error)
}

type CreatorService struct {
	Database *gorm.DB
}

func (c *CreatorService) GetOrCreate(authID string) (*domain.Creator, error) {
	var result *domain.Creator
	if err := c.Database.FirstOrCreate(&result, map[string]any{"auth_id": authID}).Error; err != nil {
		return nil, err
	}

	return result, nil
}

func (c *CreatorService) GetByID(id uuid.UUID) (*domain.Creator, error) {
	var result *domain.Creator
	if err := c.Database.FirstOrCreate(&result, id).Error; err != nil {
		return nil, err
	}

	return result, nil
}
