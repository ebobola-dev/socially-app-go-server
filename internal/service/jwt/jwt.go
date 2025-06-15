package jwt_service

import (
	"fmt"
	"time"

	"github.com/ebobola-dev/socially-app-go-server/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type IJwtService interface {
	GenerateRegistration(email string) (string, error)
	ValidateRegistration(token string) (*RegistrationClaims, error)
}

type JwtService struct {
	cfg *config.JWTConfig
}

func NewJwtService(cfg *config.JWTConfig) IJwtService {
	return &JwtService{
		cfg: cfg,
	}
}

type RegistrationClaims struct {
	Email string `json:"email_address"`
	jwt.RegisteredClaims
}

func (s *JwtService) GenerateRegistration(email string) (string, error) {
	now := time.Now()
	claims := RegistrationClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * time.Duration(s.cfg.ACCESS_DURABILITY_HOURS))),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.cfg.ACCESS_SERCER_KEY)
}

func (s *JwtService) ValidateRegistration(tokenString string) (*RegistrationClaims, error) {
	claims := &RegistrationClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.cfg.ACCESS_SERCER_KEY, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
