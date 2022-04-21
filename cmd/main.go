package main

import (
	"fmt"
	"github.com/inkoba/app_for_HR/internal/config"
	"github.com/inkoba/app_for_HR/internal/initialization"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	c, logger, err := config.LoadConfig()
	if err != nil {
		logger.Fatal("Error during loading from config file")
	}

	errs := make(chan error)

	initialization.Initialize(c, logger, errs)

	go func() {
		cs := make(chan os.Signal, 1)
		signal.Notify(cs, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-cs)
	}()

	logger.Println("exit", <-errs)
}
