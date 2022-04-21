package handlers

import (
	"fmt"
	"github.com/inkoba/app_for_HR/internal/core/ports"
	"github.com/sirupsen/logrus"
	"net/http"
)

type HealthHandler struct {
	healthService ports.IHealthService
	logger        *logrus.Logger
}

func NewHealthHandler(service ports.IHealthService, logger *logrus.Logger) *HealthHandler {
	return &HealthHandler{
		service,
		logger,
	}
}

func (hh HealthHandler) Ping(w http.ResponseWriter, _ *http.Request) {
	hh.logger.Info("entering health check end point")
	w.WriteHeader(http.StatusOK)
	err := hh.healthService.Ping()
	if err != nil {
		_, err := fmt.Fprintf(w, "Mongo is NOT up and running")
		if err != nil {
			hh.logger.Error(err)
			return
		}
	} else {
		_, err := fmt.Fprintf(w, "API is up and running")
		if err != nil {
			hh.logger.Error(err)
			return
		}
	}
}
