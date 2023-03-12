package services

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/survivorbat/qq.maarten.dev/server/domain"
	"testing"
)

func TestDBCreatorService_GetOrCreate_CreatesNewUser(t *testing.T) {
	t.Parallel()
	// Arrange
	database := getDb(t)
	autoMigrate(t, database)

	service := &DBCreatorService{Database: database}

	authID := "23902349"

	// Act
	result, err := service.GetOrCreate(authID)

	// Assert
	assert.NoError(t, err)

	assert.Equal(t, authID, result.AuthID)

	assert.NotEmpty(t, result.Nickname)
	assert.NotEmpty(t, result.ID)
}

func TestDBCreatorService_GetOrCreate_ReturnsExistingUser(t *testing.T) {
	t.Parallel()
	// Arrange
	database := getDb(t)
	autoMigrate(t, database)

	service := &DBCreatorService{Database: database}

	authID := "23902349"

	creator := &domain.Creator{
		BaseObject: domain.BaseObject{ID: uuid.MustParse("3c97f06b-1078-46ef-a2c3-71fc4d9a3d3d")},
		AuthID:     authID,
		Nickname:   "existing name",
	}

	database.Create(creator)

	// Act
	result, err := service.GetOrCreate(authID)

	// Assert
	assert.NoError(t, err)

	assert.Equal(t, creator.Nickname, result.Nickname)
	assert.Equal(t, creator.ID, result.ID)
	assert.Equal(t, creator.AuthID, result.AuthID)
}

func TestDBCreatorService_GetOrCreate_ReturnsDatabaseError(t *testing.T) {
	t.Parallel()
	// Arrange
	database := getDb(t)

	// By not running this, we're sure it will return an error
	//autoMigrate(t, database

	service := &DBCreatorService{Database: database}

	authID := "23902349"

	// Act
	result, err := service.GetOrCreate(authID)

	// Assert
	assert.Empty(t, result)
	assert.ErrorContains(t, err, "no such table")
}

func TestDBCreatorService_GetByID_ReturnsUser(t *testing.T) {
	t.Parallel()
	// Arrange
	database := getDb(t)
	autoMigrate(t, database)

	service := &DBCreatorService{Database: database}

	creator := &domain.Creator{BaseObject: domain.BaseObject{ID: uuid.MustParse("6aacfb41-e478-46ec-857e-11221f2a97fc")}}

	database.Create(creator)

	// Act
	result, err := service.GetByID(creator.ID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, creator.ID, result.ID)
}

func TestDBCreatorService_GetByID_ReturnsDatabaseError(t *testing.T) {
	t.Parallel()
	// Arrange
	database := getDb(t)

	// By not running this, we're sure it will return an error
	//autoMigrate(t, database

	service := &DBCreatorService{Database: database}

	creator := &domain.Creator{BaseObject: domain.BaseObject{ID: uuid.MustParse("6aacfb41-e478-46ec-857e-11221f2a97fc")}}

	// Act
	result, err := service.GetByID(creator.ID)

	// Assert
	assert.Empty(t, result)
	assert.ErrorContains(t, err, "no such table")
}
