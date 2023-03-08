package inputs

// TokenInput is used to exchange the Google login code to an authenticated JWT token
type TokenInput struct {
	Code string `json:"code" example:"..."` // desc: The token obtained from your Google login
}
