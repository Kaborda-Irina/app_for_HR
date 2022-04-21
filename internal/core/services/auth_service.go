package services

import (
	"github.com/inkoba/app_for_HR/internal/core/ports"
	"github.com/sirupsen/logrus"
)

type AuthService struct {
	logger         *logrus.Logger
	userRepository ports.IUserRepository
	cryptoService  ports.ICryptoService
}

var _ ports.IAuthService = (*AuthService)(nil)

func NewAuthService(userIRepository ports.IUserRepository, logger *logrus.Logger, cryptoService ports.ICryptoService) *AuthService {
	return &AuthService{
		logger,
		userIRepository,
		cryptoService,
	}
}

func (a AuthService) IsValidUser(username string, password string) error {
	user, err := a.userRepository.GetUserByUsername(username)
	if err != nil {
		a.logger.Error("Error when getting user by username ", err)
		return err
	}

	err = a.cryptoService.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return err
	}
	return nil
}
