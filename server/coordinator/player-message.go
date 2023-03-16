package coordinator

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/survivorbat/qq.maarten.dev/server/inputs"
)

var validate = validator.New()

func init() {
	validate.SetTagName("binding")
}

type PlayerMessage struct {
	Action  PlayerAction    `json:"action"`
	Content json.RawMessage `json:"content"`

	// Optional
	Answer *inputs.Answer `json:"-"`
}

func (p *PlayerMessage) IsValid() bool {
	if !p.Action.IsValid() {
		return false
	}

	switch p.Action {
	case AnswerAction:
		if err := validate.Struct(p.Answer); err != nil {
			logrus.WithError(err).Error("Failed to validate")
			return false
		}
	}

	return true
}

func (p *PlayerMessage) Parse() error {
	switch p.Action {
	case AnswerAction:
		if err := json.Unmarshal(p.Content, &p.Answer); err != nil {
			logrus.WithError(err).Error("Failed to parse")
			return err
		}

		return nil

	default:
		logrus.Errorf("Unknown action %s", p.Action)
		return errors.New("unknown type")
	}
}
