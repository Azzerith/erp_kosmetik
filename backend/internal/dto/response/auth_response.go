package response

type LoginResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	TokenType    string      `json:"token_type"`
	ExpiresIn    int         `json:"expires_in"`
	User         UserResponse `json:"user"`
}

type UserResponse struct {
	ID             uint64  `json:"id"`
	Email          string  `json:"email"`
	Name           string  `json:"name"`
	Phone          *string `json:"phone,omitempty"`
	AvatarURL      *string `json:"avatar_url,omitempty"`
	Role           string  `json:"role"`
	LoyaltyPoints  int     `json:"loyalty_points"`
	EmailVerified  bool    `json:"email_verified"`
	CreatedAt      string  `json:"created_at"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}