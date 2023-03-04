package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type BaseObject struct {
	ID        uuid.UUID `gorm:"type:uuid" json:"id" example:"00000000-0000-0000-0000-000000000000"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (b *BaseObject) BeforeCreate(*gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}

	return nil
}
