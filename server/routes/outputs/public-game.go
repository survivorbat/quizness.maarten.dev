package outputs

import (
	"github.com/google/uuid"
	"github.com/survivorbat/qq.maarten.dev/server/domain"
)

func NewPublicGame(game *domain.Game) *OutputGame {
	return &OutputGame{ID: game.ID}
}

type OutputGame struct {
	ID uuid.UUID `json:"id"`
}
