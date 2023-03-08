package inputs

type QuestionOption struct {
	Type       string `json:"type"`
	TextOption string `json:"textOption"`
}

type Question struct {
	Title             string `json:"title" binding:"required"`
	Description       string `json:"description" binding:"required"`
	DurationInSeconds uint   `json:"durationInSeconds" binding:"required"`
	Category          string `json:"category" binding:"required"`
	Order             uint   `json:"order" binding:"required"`
	Type              string `json:"type" binding:"required"`

	Options []*QuestionOption `json:"options"`
}

type Quiz struct {
	Name        string      `json:"name" binding:"required"`
	Description string      `json:"description" binding:"required"`
	Questions   []*Question `json:"questions"`
}
