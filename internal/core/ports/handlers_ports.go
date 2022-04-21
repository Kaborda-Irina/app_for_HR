package ports

import (
	"net/http"
)

type IHealthHandler interface {
	Ping(w http.ResponseWriter, r *http.Request)
}

type ISalaryHandler interface {
	UploadFile(w http.ResponseWriter, r *http.Request)
}
type IUserHandler interface {
	GetAll(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type IAuthHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
}
type IMiddlewareHandler interface {
	LogURL(next http.Handler) http.Handler
	CheckJWT(next http.Handler) http.Handler
}

type IFilterHandler interface {
	Filter(w http.ResponseWriter, r *http.Request)
}
