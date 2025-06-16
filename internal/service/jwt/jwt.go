package jwt_service

import (
	"fmt"
	"time"

	"github.com/ebobola-dev/socially-app-go-server/internal/config"
	"github.com/ebobola-dev/socially-app-go-server/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type IJwtService interface {
	GenerateRegistration(email string) (string, error)
	ValidateRegistration(token string) (*RegistrationClaims, error)
	GenerateUserPair(userId uuid.UUID, deviceId string) (string, *model.RefreshToken, error)
	ValidateUserAccess(accessTokenString string) (*UserClaims, error)
	ValidateUserRefresh(refreshTokenString string) (*UserClaims, error)
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

type UserClaims struct {
	ID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

func (s *JwtService) GenerateRegistration(email string) (string, error) {
	now := time.Now()
	claims := RegistrationClaims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute * time.Duration(s.cfg.ACCESS_DURABILITY_MIN))),
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

func (s *JwtService) GenerateUserPair(userId uuid.UUID, deviceId string) (string, *model.RefreshToken, error) {
	now := time.Now()
	access_registered_claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute * time.Duration(s.cfg.ACCESS_DURABILITY_MIN))),
		IssuedAt:  jwt.NewNumericDate(now),
	}
	refresh_registered_claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour * time.Duration(s.cfg.REFRESH_DURABILITY_DAYS))),
		IssuedAt:  jwt.NewNumericDate(now),
	}
	access_claims := UserClaims{
		ID:               userId,
		RegisteredClaims: access_registered_claims,
	}
	refresh_claims := UserClaims{
		ID:               userId,
		RegisteredClaims: refresh_registered_claims,
	}
	access_token := jwt.NewWithClaims(jwt.SigningMethodHS256, access_claims)
	refresh_token := jwt.NewWithClaims(jwt.SigningMethodHS256, refresh_claims)
	access_string_token, a_err := access_token.SignedString(s.cfg.ACCESS_SERCER_KEY)
	if a_err != nil {
		return "", nil, a_err
	}
	refresh_string_token, r_err := refresh_token.SignedString(s.cfg.REFRESH_SERCER_KEY)
	if r_err != nil {
		return "", nil, r_err
	}
	rt_obj := &model.RefreshToken{
		UserID:    userId,
		DeviceID:  deviceId,
		Value:     refresh_string_token,
		ExpiresAt: refresh_registered_claims.ExpiresAt.Time,
	}
	return access_string_token, rt_obj, nil
}

func (s *JwtService) ValidateUserAccess(accessTokenString string) (*UserClaims, error) {
	claims := &UserClaims{}
	token, err := jwt.ParseWithClaims(accessTokenString, claims, func(token *jwt.Token) (interface{}, error) {
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

func (s *JwtService) ValidateUserRefresh(refreshTokenString string) (*UserClaims, error) {
	claims := &UserClaims{}
	token, err := jwt.ParseWithClaims(refreshTokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.cfg.REFRESH_SERCER_KEY, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
