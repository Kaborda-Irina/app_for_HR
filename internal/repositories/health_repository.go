package repositories

import (
	"github.com/inkoba/app_for_HR/internal/core/ports"
	"github.com/sirupsen/logrus"
)

type HealthRepository struct {
	mc     *MongoConfig
	logger *logrus.Logger
}

//This line is for get feedback in case we are not implementing the interface correctly
var _ ports.IHealthRepository = (*HealthRepository)(nil)

func NewHealthRepository(mc *MongoConfig, logger *logrus.Logger) *HealthRepository {
	return &HealthRepository{
		mc,
		logger,
	}
}

func (h HealthRepository) Ping() error {
	err := h.mc.Ping()
	if err != nil {
		h.logger.Error(err)
		return err
	}
	return nil
}
