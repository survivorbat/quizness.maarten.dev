package domain

type Creator struct {
	BaseObject

	NickName string `json:"nickname"`

	// Never expose this
	AuthID string `json:"-" gorm:"unique"`

	Quizes []*Quiz `json:"quizes" gorm:"foreignKey:CreatorID"`
}
