package services

import (
	"github.com/inkoba/app_for_HR/internal/core/domain"
	"github.com/inkoba/app_for_HR/internal/core/ports"
	"github.com/sirupsen/logrus"
)

type UserService struct {
	userRepository ports.IUserRepository
	logger         *logrus.Logger
	appCrypto      ports.ICryptoService
}

var _ ports.IUserService = (*UserService)(nil)

func NewUserService(userRepository ports.IUserRepository, logger *logrus.Logger, appCrypto ports.ICryptoService) *UserService {
	return &UserService{
		userRepository,
		logger,
		appCrypto,
	}
}

func (us UserService) Get(id string) (*domain.User, error) {
	user, err := us.userRepository.Get(id)
	if err != nil {
		us.logger.Error(err)
		return nil, err
	}

	return user, nil
}

func (us UserService) GetAll() ([]*domain.User, error) {
	users, err := us.userRepository.GetAll()
	if err != nil {
		us.logger.Error(err)
		return nil, err
	}
	return users, err
}

func (us UserService) Create(user *domain.User) (string, error) {
	hashedPassword, err := us.appCrypto.GetHashedPassword([]byte(user.Password))
	if err != nil {
		us.logger.Error(err)
		return "", err
	}
	newUser := domain.User{
		Username: user.Username,
		Password: hashedPassword,
		IsAdmin:  false,
	}

	result, err := us.userRepository.Create(&newUser)
	if err != nil {
		us.logger.Error(err)
		return "", err
	}
	return result, err
}

func (us UserService) Delete(id string) error {
	err := us.userRepository.Delete(id)
	if err != nil {
		us.logger.Error(err)
	}
	return err
}

func (us UserService) GetUserByUsername(username string) (*domain.User, error) {
	user, err := us.userRepository.GetUserByUsername(username)
	if err != nil {
		us.logger.Error(err)
		return nil, err
	}
	return user, nil
}
