package handlers

import (
	"encoding/json"
	"github.com/inkoba/app_for_HR/internal/core/domain/response"
	"github.com/sirupsen/logrus"

	"net/http"
)

func HandleError(w http.ResponseWriter, message string, logger *logrus.Logger) {
	errorMessage := response.ErrorMessage{}
	errorMessage.Errors = append(errorMessage.Errors, message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	err := json.NewEncoder(w).Encode(&errorMessage)
	if err != nil {
		logger.Error(err)
	}
	return
}
