package inputs

// Token is used to exchange the Google login code to an authenticated JWT token
type Token struct {
	Code string `json:"code" example:"..."` // desc: The token obtained from your Google login
}
