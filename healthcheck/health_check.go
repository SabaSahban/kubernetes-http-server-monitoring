package healthcheck

import (
	"cloud-final/storage"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func StartHealthCheck(storageClient *storage.Storage, interval int) {
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	for {
		select {
		case <-ticker.C:
			checkServersHealth(storageClient)
		}
	}
}

func checkServersHealth(storageClient *storage.Storage) {
	servers, err := storageClient.GetAllServers()
	if err != nil {
		logrus.Error("Failed to retrieve servers: ", err)
		return
	}

	for _, server := range servers {
		checkServer(&server)
		err := storageClient.UpdateServer(server)
		if err != nil {
			logrus.Errorf("Failed to update server %s: %v", server.ID, err)
		}
	}
}

func checkServer(server *storage.ServerData) {
	resp, err := http.Get(server.Address)
	if err != nil {
		server.FailureCount++
		server.LastFailure = time.Now()
		logrus.Errorf("Failed to reach server %s: %v", server.Address, err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode == http.StatusOK {
		server.SuccessCount++
	} else {
		server.FailureCount++
		server.LastFailure = time.Now()
	}
}
