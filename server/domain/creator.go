package domain

type Creator struct {
	BaseObject

	NickName string `json:"nickname" gorm:"unique"`

	// Never expose this
	AuthID string `json:"-" gorm:"unique"`

	Quizzes []*Quiz `json:"quizzes" gorm:"foreignKey:CreatorID"`
}
