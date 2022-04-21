package services

import (
	"github.com/inkoba/app_for_HR/internal/core/ports"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type CryptoService struct {
	logger *logrus.Logger
}

var _ ports.ICryptoService = (*CryptoService)(nil)

func NewHashPassword(logger *logrus.Logger) *CryptoService {
	return &CryptoService{
		logger,
	}
}
func (hp CryptoService) GetHashedPassword(password []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(password, 12)
	if err != nil {
		hp.logger.Error(err)
		return "", err
	}
	return string(hash), nil
}

func (hp CryptoService) CompareHashAndPassword(hashedPassword []byte, password []byte) error {
	err := bcrypt.CompareHashAndPassword(hashedPassword, password)
	if err != nil {
		hp.logger.Error(err)
		return err
	}
	return nil
}
