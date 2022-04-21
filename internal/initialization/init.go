package initialization

import (
	"github.com/gorilla/mux"
	"github.com/inkoba/app_for_HR/internal/config"
	"github.com/inkoba/app_for_HR/internal/core/handlers"
	"github.com/inkoba/app_for_HR/internal/core/services"
	"github.com/inkoba/app_for_HR/internal/repositories"
	"github.com/sirupsen/logrus"
	"net/http"
)

func Initialize(c config.Config, logger *logrus.Logger, errs chan error) {
	logger.Println("Starting mongo connection")
	mongoConfig := repositories.NewMongoConfig(c, logger)
	logger.Println("Mongo connection is successful")

	healthRepository := repositories.NewHealthRepository(mongoConfig, logger)
	userRepository := repositories.NewUserRepository(mongoConfig, logger)
	salaryRepository := repositories.NewSalaryRepository(mongoConfig, logger)

	appCrypto := services.NewHashPassword(logger)
	userService := services.NewUserService(userRepository, logger, appCrypto)
	authService := services.NewAuthService(userRepository, logger, appCrypto)
	healthService := services.NewHealthService(healthRepository, logger)
	salaryService := services.NewSalaryService(c.CurrencyConfig, salaryRepository, logger)

	userHandler := handlers.NewUserHandler(userService, logger)
	authHandler := handlers.NewAuthHandler(authService, userService, logger)
	healthHandler := handlers.NewHealthHandler(healthService, logger)
	salaryHandler := handlers.NewSalaryHandler(salaryService, logger)
	middlewareHandler := handlers.NewMiddlewareHandler(logger)
	filterHandler := handlers.NewSalaryFilterHandler(salaryService, logger)

	router := mux.NewRouter()
	router.Use(middlewareHandler.LogURL)
	subRouter := router.PathPrefix("/api/").Subrouter()
	subRouter.Use(middlewareHandler.CheckJWT)
	logger.Println("Сreating routes")

	router.HandleFunc("/api/health", healthHandler.Ping).Methods("GET")
	subRouter.HandleFunc("/users", userHandler.GetAll).Methods("GET")
	subRouter.HandleFunc("/users/{id:[a-zA-Z0-9]*}", userHandler.Get).Methods("GET")
	subRouter.HandleFunc("/users", userHandler.Create).Methods("POST")
	subRouter.HandleFunc("/users/{id:[a-zA-Z0-9]*}", userHandler.Delete).Methods("DELETE")

	router.HandleFunc("/api/salaries", salaryHandler.UploadFile).Methods("POST")

	router.HandleFunc("/api/login", authHandler.Login).Methods("POST")
	router.HandleFunc("/api/filter", filterHandler.Filter).Methods("POST")
	http.Handle("/", router)

	go func() {
		logger.Println("transport", "HTTP", "addr", ":"+c.Port)
		errs <- http.ListenAndServe(":"+c.Port, router)
	}()
}
