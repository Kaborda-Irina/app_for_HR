package handlers

import (
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

type MiddlewareHandler struct {
	logger *logrus.Logger
}

func NewMiddlewareHandler(logger *logrus.Logger) *MiddlewareHandler {
	return &MiddlewareHandler{
		logger: logger,
	}
}

func (mw MiddlewareHandler) LogURL(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mw.logger.Info(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func (mw MiddlewareHandler) CheckJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mw.logger.Info("Start jwtMiddleware")
		tokenString := r.Header.Get("jwt")

		claims := jwt.MapClaims{}

		_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})
		if err != nil {
			mw.logger.Error("Error parsing jwt: ", err)
		}

		if claims["isAdmin"] == true {
			mw.logger.Info("Authenticated user")
			next.ServeHTTP(w, r)
		} else {
			mw.logger.Info("Unauthenticated user")
			w.WriteHeader(http.StatusForbidden)
			return
		}

	})
}
