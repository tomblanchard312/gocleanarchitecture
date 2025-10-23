package auth

// TokenGeneratorAdapter adapts JWTManager to the TokenGenerator interface
type TokenGeneratorAdapter struct {
	jwtManager *JWTManager
}

func NewTokenGeneratorAdapter(jwtManager *JWTManager) *TokenGeneratorAdapter {
	return &TokenGeneratorAdapter{
		jwtManager: jwtManager,
	}
}

func (a *TokenGeneratorAdapter) GenerateToken(userID, username, email string) (string, error) {
	return a.jwtManager.GenerateToken(userID, username, email)
}

func (a *TokenGeneratorAdapter) ValidateToken(token string) (userID string, username string, email string, err error) {
	claims, err := a.jwtManager.ValidateToken(token)
	if err != nil {
		return "", "", "", err
	}
	return claims.UserID, claims.Username, claims.Email, nil
}

