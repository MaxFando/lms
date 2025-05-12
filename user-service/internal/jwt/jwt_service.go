package jwt

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
	"github.com/MaxFando/lms/user-service/internal/model"
)

type JWTService interface {
	GenerateTokens(user *model.User) (string, string, error)
	ParseToken(token string) (*UserClaims, error)
}

type jwtService struct {
	secret     string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

type UserClaims struct {
	UserID int64  `json:"user_id"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func NewJWTService(secret string, accessTTL, refreshTTL time.Duration) JWTService {
	return &jwtService{
		secret:     secret,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (j *jwtService) GenerateTokens(user *model.User) (string, string, error) {
	accessClaims := UserClaims{
		UserID: user.ID,
		Name:   user.Name,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.accessTTL)),
		},
	}
	refreshClaims := UserClaims{
		UserID: user.ID,
		Name:   user.Name,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.refreshTTL)),
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(j.secret))
	if err != nil {
		return "", "", err
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(j.secret))
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (j *jwtService) ParseToken(tokenStr string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	return claims, nil
}
