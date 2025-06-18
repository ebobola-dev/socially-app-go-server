package hash

import "golang.org/x/crypto/bcrypt"

type IHashService interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(password, hash string) bool
}

type HashService struct{}

func NewHashService() IHashService {
	return &HashService{}
}

func (s *HashService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	return string(bytes), err
}

func (s *HashService) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
