package routes

import (
	"github.com/google/uuid"
	"github.com/survivorbat/qq.maarten.dev/server/domain"
	"github.com/survivorbat/qq.maarten.dev/server/services"
)

type MockCreatorService struct {
	services.ICreatorService

	getByIDCalledWith   uuid.UUID
	getByIDReturns      *domain.Creator
	getByIDReturnsError error
}

func (m *MockCreatorService) GetByID(id uuid.UUID) (*domain.Creator, error) {
	m.getByIDCalledWith = id
	return m.getByIDReturns, m.getByIDReturnsError
}
