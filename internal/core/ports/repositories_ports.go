package ports

import (
	"github.com/inkoba/app_for_HR/internal/core/domain"
	"github.com/inkoba/app_for_HR/internal/core/domain/request"
)

//go:generate mockgen -source=repositories_ports.go -destination=mocks/mock_repository.go

type IUserRepository interface {
	Get(id string) (*domain.User, error)
	GetAll() ([]*domain.User, error)
	Create(user *domain.User) (string, error)
	Delete(id string) error
	GetUserByUsername(username string) (*domain.User, error)
}
type ISalaryRepository interface {
	Create(salaries []*domain.Salary) error
	DeleteAll() error
	GetFilteredSalaries(filterSalary *request.ConditionForFilteringSalaries) ([]*domain.Salary, error)
}

type IHealthRepository interface {
	Ping() error
}
