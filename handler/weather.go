package handler

import (
	"awesomeProject/config"
	"awesomeProject/storage"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

const (
	redisKeyPrefix = "weather:"
	apiURLFormat   = "https://api.api-ninjas.com/v1/weather?city=%s"
)

type WeatherHandler struct {
	RedisClient  *storage.Storage
	NinjasConfig config.Ninjas
}

type WeatherInfo struct {
	WindSpeed       float64 `json:"wind_speed"`
	WindDegrees     int     `json:"wind_degrees"`
	Temperature     int     `json:"temp"`
	Humidity        int     `json:"humidity"`
	Sunset          int64   `json:"sunset"`
	MinTemperature  int     `json:"min_temp"`
	CloudPercentage int     `json:"cloud_pct"`
	FeelsLike       int     `json:"feels_like"`
	Sunrise         int64   `json:"sunrise"`
	MaxTemperature  int     `json:"max_temp"`
}

func (h *WeatherHandler) getWeatherInfo(city string) (*WeatherInfo, error) {
	podName := os.Getenv("HOSTNAME")
	podIP := os.Getenv("POD_IP")

	logrus.Infof("Handling request in pod %s (%s)", podName, podIP)

	redisKey := fmt.Sprintf("%s%s", redisKeyPrefix, city)
	cachedData, err := h.RedisClient.Read(redisKey)

	if err != nil && err != redis.Nil {
		logrus.WithError(err).Errorf("Failed to read from Redis for city %s", city)
		return nil, err
	}

	if err == nil {
		var weatherInfo WeatherInfo
		if err := json.Unmarshal([]byte(cachedData), &weatherInfo); err == nil {
			logrus.Infof("Data found in Redis for city %s: %+v", city, weatherInfo)
			return &weatherInfo, nil
		}
	}

	apiURL := fmt.Sprintf(apiURLFormat, city)
	apiKey := h.NinjasConfig.ApiKey

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		logrus.WithError(err).Errorf("Failed to create HTTP request for city %s", city)
		return nil, err
	}

	req.Header.Set("X-Api-Key", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logrus.WithError(err).Errorf("Failed to make API request for city %s", city)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("API request failed with status code %d", resp.StatusCode)
		logrus.Errorf("API request failed for city %s with status code %d", city, resp.StatusCode)
		return nil, err
	}

	var weatherInfo WeatherInfo
	err = json.NewDecoder(resp.Body).Decode(&weatherInfo)
	if err != nil {
		logrus.WithError(err).Errorf("Failed to decode JSON for city %s", city)
		return nil, err
	}

	err = h.RedisClient.Write(redisKey, &weatherInfo)
	if err != nil {
		logrus.WithError(err).Errorf("Failed to write to Redis for city %s", city)
		return nil, err
	}

	logrus.Infof("Weather info retrieved successfully for city %s: %+v", city, weatherInfo)
	return &weatherInfo, nil
}

func (h *WeatherHandler) WeatherInfo(c echo.Context) error {
	city := c.Param("city")

	weatherInfo, err := h.getWeatherInfo(city)
	if err != nil {
		logrus.WithError(err).Errorf("Failed to retrieve weather info for city %s", city)
		return c.JSON(http.StatusInternalServerError,
			map[string]string{"error": "Failed to retrieve weather information"})
	}

	return c.JSON(http.StatusOK, weatherInfo)
}
