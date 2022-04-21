package handlers

import (
	"encoding/json"
	"github.com/inkoba/app_for_HR/internal/core/domain/request"
	"github.com/inkoba/app_for_HR/internal/core/ports"
	"github.com/sirupsen/logrus"
	"net/http"
)

type SalaryFilterHandler struct {
	salaryService ports.ISalaryService
	logger        *logrus.Logger
}

func NewSalaryFilterHandler(salaryService ports.ISalaryService, logger *logrus.Logger) *SalaryFilterHandler {
	return &SalaryFilterHandler{
		salaryService,
		logger,
	}
}

func (fh SalaryFilterHandler) Filter(w http.ResponseWriter, r *http.Request) {
	salaryFilteringCondition := request.ConditionForFilteringSalaries{}
	err := json.NewDecoder(r.Body).Decode(&salaryFilteringCondition)
	if err != nil {
		fh.logger.Error("Error decode in SalariesResponse struct", err)
		HandleError(w, err.Error(), fh.logger)
		return
	}

	filteredSalaries, err := fh.salaryService.GetSalariesByFilter(&salaryFilteringCondition)
	if err != nil {
		fh.logger.Error("Error getting filtered salaries", err)
		HandleError(w, err.Error(), fh.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(filteredSalaries)
	if err != nil {
		HandleError(w, err.Error(), fh.logger)
	}
}
