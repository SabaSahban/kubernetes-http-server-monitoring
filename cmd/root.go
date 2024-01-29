package cmd

import (
	"cloud-final/config"
	"cloud-final/handler"
	"cloud-final/storage"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// Execute application.
func Execute() {
	// load configs
	cfg, err := config.Load()
	if err != nil {
		logrus.Fatalf(err.Error())
	}

	// create redis client
	rc := storage.New(cfg)

	weatherHandler := handler.WeatherHandler{
		RedisClient:  rc,
		NinjasConfig: cfg.Ninjas,
	}

	e := echo.New()
	e.GET("/weather/:city", weatherHandler.WeatherInfo)
	port := cfg.Server.Port
	serverAddress := fmt.Sprintf(":%d", port)
	if err := e.Start(serverAddress); !errors.Is(err, http.ErrServerClosed) {
		logrus.Fatalf("Server error: %v", err)
	}
}
