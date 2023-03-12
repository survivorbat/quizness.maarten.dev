package inputs

type Game struct {
	PlayerLimit uint `json:"playerLimit" example:"25" binding:"required,min=2,max=25"` // desc: The max amount of players that may join this game
}
