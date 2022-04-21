package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/inkoba/app_for_HR/internal/core/ports"
	"github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"net/http"
)

type SalaryHandler struct {
	logger        *logrus.Logger
	salaryService ports.ISalaryService
}

var _ ports.ISalaryHandler = (*SalaryHandler)(nil)

func NewSalaryHandler(salaryService ports.ISalaryService, logger *logrus.Logger) ports.ISalaryHandler {
	return &SalaryHandler{
		logger,
		salaryService,
	}
}

func (sh SalaryHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		sh.logger.Error(err)
		HandleError(w, err.Error(), sh.logger)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		sh.logger.Error(err)
		HandleError(w, err.Error(), sh.logger)
		return
	}

	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			sh.logger.Error(err)
			return
		}
	}(file)

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		sh.logger.Error(err)
		HandleError(w, err.Error(), sh.logger)
		return
	}

	report, err := sh.salaryService.Create(buf.Bytes())
	if err != nil {
		sh.logger.Error(err)
		HandleError(w, err.Error(), sh.logger)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&report)
	if err != nil {
		sh.logger.Error(err)
		HandleError(w, err.Error(), sh.logger)
	}
}
