package dto

type LoginInput struct {
	Email    string `json:"email" binding:"required,email,email_advanced"`
	Password string `json:"password" binding:"required"`
}
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiredAt    string `json:"expired_at"`
}
type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
