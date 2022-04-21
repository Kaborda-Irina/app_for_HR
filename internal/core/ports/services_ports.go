package ports

import (
	"github.com/inkoba/app_for_HR/internal/core/domain"
	"github.com/inkoba/app_for_HR/internal/core/domain/request"
	"github.com/inkoba/app_for_HR/internal/core/domain/response"
)

//go:generate mockgen -source=services_ports.go -destination=mocks/mock.go

type IHealthService interface {
	Ping() error
}
type IUserService interface {
	Get(id string) (*domain.User, error)
	GetAll() ([]*domain.User, error)
	Create(user *domain.User) (string, error)
	Delete(id string) error
	GetUserByUsername(username string) (*domain.User, error)
}
type IAuthService interface {
	IsValidUser(username string, password string) error
}
type ISalaryService interface {
	Create(file []byte) (*response.SalaryUploadReport, error)
	GetSalariesByFilter(filterSalary *request.ConditionForFilteringSalaries) ([]*response.SalariesResponse, error)
}

type ICryptoService interface {
	GetHashedPassword([]byte) (string, error)
	CompareHashAndPassword(hashedPassword []byte, password []byte) error
}
