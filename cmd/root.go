package cmd

import (
	"cloud-final/config"
	"cloud-final/handler"
	"cloud-final/healthcheck"
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

	serverHandler := handler.ServerHandler{
		RedisClient: rc,
	}

	go healthcheck.StartHealthCheck(rc, cfg.CheckInterval)

	e := echo.New()
	e.POST("/api/server", serverHandler.AddServer)
	e.GET("/api/server", serverHandler.GetServerStatus)
	e.GET("/api/server/all", serverHandler.GetAllServersStatus)

	port := cfg.Server.Port
	serverAddress := fmt.Sprintf(":%d", port)
	if err := e.Start(serverAddress); !errors.Is(err, http.ErrServerClosed) {
		logrus.Fatalf("Server error: %v", err)
	}
}
