package inputs

type QuestionOptionInput struct {
	Type       string `json:"type"`
	TextOption string `json:"textOption"`
}

type QuestionInput struct {
	Title             string `json:"title" binding:"required"`
	Description       string `json:"description" binding:"required"`
	DurationInSeconds uint   `json:"durationInSeconds" binding:"required"`
	Category          string `json:"category" binding:"required"`
	Order             uint   `json:"order" binding:"required"`
	Type              string `json:"type" binding:"required"`

	Options []*QuestionOptionInput `json:"options"`
}

type QuizInput struct {
	Name        string           `json:"name" binding:"required"`
	Description string           `json:"description" binding:"required"`
	Questions   []*QuestionInput `json:"questions"`
}
