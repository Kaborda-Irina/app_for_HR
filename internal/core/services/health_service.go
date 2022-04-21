package services

import (
	"github.com/inkoba/app_for_HR/internal/core/ports"
	"github.com/sirupsen/logrus"
)

type HealthService struct {
	healthRepository ports.IHealthService
	logger           *logrus.Logger
}

//This line is for get feedback in case we are not implementing the interface correctly
var _ ports.IHealthService = (*HealthService)(nil)

func NewHealthService(healthRepository ports.IHealthService, logger *logrus.Logger) *HealthService {
	return &HealthService{
		healthRepository,
		logger,
	}
}

func (h HealthService) Ping() error {
	err := h.healthRepository.Ping()
	if err != nil {
		h.logger.Error(err)
		return err
	}
	return nil
}
