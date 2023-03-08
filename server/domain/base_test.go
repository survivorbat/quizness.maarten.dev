package domain

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

func TestBaseObject_BeforeCreate_SetsUUID(t *testing.T) {
	t.Parallel()
	// Arrange
	object := BaseObject{}

	// Act
	err := object.BeforeCreate(&gorm.DB{})

	// Assert
	assert.Nil(t, err)
	assert.NotEqual(t, uuid.Nil, object.ID)
}

func TestBaseObject_BeforeCreate_DoesNotOverwriteUUID(t *testing.T) {
	t.Parallel()
	// Arrange
	object := BaseObject{
		ID: uuid.MustParse("1d5858c2-d57a-47d7-b800-5268bdf2dfe4"),
	}

	// Act
	err := object.BeforeCreate(&gorm.DB{})

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, uuid.MustParse("1d5858c2-d57a-47d7-b800-5268bdf2dfe4"), object.ID)
}
