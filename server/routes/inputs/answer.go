package inputs

import "github.com/google/uuid"

type Answer struct {
	OptionID uuid.UUID `json:"optionID" example:"00000000-0000-0000-0000-000000000000"`
}
