package handler

import (
	"cloud-final/storage"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type ServerHandler struct {
	RedisClient *storage.Storage
}

func (h *ServerHandler) AddServer(c echo.Context) error {
	serverAddress := c.FormValue("address")
	if serverAddress == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Server address is required"})
	}

	id := strconv.FormatInt(time.Now().UnixNano(), 10)

	serverData := storage.ServerData{
		ID:           id,
		Address:      serverAddress,
		SuccessCount: 0,
		FailureCount: 0,
		CreatedAt:    time.Now(),
	}

	err := h.RedisClient.AddServer(serverData)
	if err != nil {
		logrus.Errorf("Error adding server: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error adding server"})
	}

	return c.JSON(http.StatusOK, map[string]string{"id": id})
}

func (h *ServerHandler) GetServerStatus(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Server ID is required"})
	}

	serverData, err := h.RedisClient.GetServer(id)
	if err != nil {
		logrus.Errorf("Error retrieving server: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error retrieving server"})
	}

	return c.JSON(http.StatusOK, serverData)
}

func (h *ServerHandler) GetAllServersStatus(c echo.Context) error {
	servers, err := h.RedisClient.GetAllServers()
	if err != nil {
		logrus.Errorf("Error retrieving servers: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Error retrieving servers"})
	}

	return c.JSON(http.StatusOK, servers)
}
