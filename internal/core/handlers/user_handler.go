package handlers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/inkoba/app_for_HR/internal/core/domain"
	"github.com/inkoba/app_for_HR/internal/core/ports"
	"github.com/sirupsen/logrus"
	"net/http"
)

type UserHandler struct {
	userService ports.IUserService
	logger      *logrus.Logger
}

func NewUserHandler(service ports.IUserService, logger *logrus.Logger) ports.IUserHandler {
	return UserHandler{
		service,
		logger,
	}
}

func (ah UserHandler) GetAll(w http.ResponseWriter, _ *http.Request) {
	users, err := ah.userService.GetAll()
	if err != nil {
		ah.logger.Error(err)
		HandleError(w, err.Error(), ah.logger)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		ah.logger.Error(err)
		HandleError(w, err.Error(), ah.logger)
	}
}

func (ah UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	user, err := ah.userService.Get(id)
	if err != nil {
		ah.logger.Error(err)
		HandleError(w, err.Error(), ah.logger)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		ah.logger.Error(err)
		HandleError(w, err.Error(), ah.logger)
	}
}

func (ah UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user domain.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		ah.logger.Error("Unable to decode request body ", err)
		HandleError(w, err.Error(), ah.logger)
		return
	}

	_, err = ah.userService.GetUserByUsername(user.Username)
	if err == nil {
		ah.logger.Error("User is existed in database ", err)
		HandleError(w, errors.New("User is exist in database").Error(), ah.logger)
		return
	}

	id, err := ah.userService.Create(&user)
	if err != nil {
		ah.logger.Error(err)
		HandleError(w, err.Error(), ah.logger)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(id)
	if err != nil {
		ah.logger.Error(err)
		HandleError(w, err.Error(), ah.logger)
	}
}
func (ah UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := ah.userService.Delete(id)
	if err != nil {
		ah.logger.Error(err)
		HandleError(w, err.Error(), ah.logger)
	}
}
