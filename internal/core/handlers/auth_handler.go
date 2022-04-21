package handlers

import (
	"encoding/json"
	"github.com/golang-jwt/jwt"
	"github.com/inkoba/app_for_HR/internal/core/ports"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

type Credentials struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

type Claims struct {
	Username string `json:"username"`
	IsAdmin  bool   `json:"isAdmin"`
	jwt.StandardClaims
}

type AuthHandler struct {
	authService ports.IAuthService
	userService ports.IUserService
	logger      *logrus.Logger
}

func NewAuthHandler(service ports.IAuthService, userService ports.IUserService, logger *logrus.Logger) ports.IAuthHandler {
	return AuthHandler{
		service,
		userService,
		logger,
	}
}

func (ah AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ah.logger.Info("Start LogURL")
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		ah.logger.Error("Error decode in Credentials struct", err)
		HandleError(w, err.Error(), ah.logger)
		return
	}

	err = ah.authService.IsValidUser(creds.Username, creds.Password)
	if err != nil {
		ah.logger.Error("Error validation username or password", err)
		HandleError(w, err.Error(), ah.logger)
		return
	}
	ah.logger.Info("User is valid ", creds.Username)

	expirationTime := time.Now().Add(15 * time.Minute)

	user, err := ah.userService.GetUserByUsername(creds.Username)
	if err != nil {
		ah.logger.Error("Error get user", err)
		HandleError(w, err.Error(), ah.logger)
		return
	}

	claims := &Claims{
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtKey := os.Getenv("JWT_SECRET_KEY")

	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		ah.logger.Error("Error create token", err)
		HandleError(w, err.Error(), ah.logger)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true,
	})

}
