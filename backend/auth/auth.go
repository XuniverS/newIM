package auth

// User 用户信息
type User struct {
	ID       string
	Username string
	Email    string
	Password string
}

// Token 认证令牌
type Token struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    int64
}

// Authenticator 认证接口
type Authenticator interface {
	Register(username, email, password string) (*User, error)
	Login(username, password string) (*Token, error)
	ValidateToken(token string) (*User, error)
	RefreshToken(refreshToken string) (*Token, error)
}
