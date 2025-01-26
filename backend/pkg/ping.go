package pkg

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/chamanbravo/upstat/internal/models"
	"github.com/chamanbravo/upstat/internal/repository"
)

func Ping(monitor *models.Monitor) *models.Heartbeat {
	heartbeat := new(models.Heartbeat)
	heartbeat.MonitorId = monitor.ID
	fmt.Printf("Pinging %v at %v \n", monitor.Name, monitor.Url)
	startTime := time.Now()
	response, err := http.Get(monitor.Url)
	if err != nil {
		heartbeat.Status = "red"
		heartbeat.StatusCode = "error"
		heartbeat.Message = "unable to ping"
		heartbeat.Latency = 0
		if monitor.Status != "red" {
			err := repository.UpdateMonitorStatus(monitor.ID, "red")
			if err != nil {
				log.Printf("Error when trying to update monitor status: %v", err.Error())
			}
		}
	} else {
		heartbeat.Status = "green"
		heartbeat.StatusCode = strings.Split(response.Status, " ")[0]
		heartbeat.Message = strings.Split(response.Status, " ")[1]
		latency := time.Since(startTime).Milliseconds()
		heartbeat.Latency = int(latency)
		defer response.Body.Close()

		if monitor.Status != "green" {
			err := repository.UpdateMonitorStatus(monitor.ID, "green")
			if err != nil {
				log.Printf("Error when trying to update monitor status: %v", err.Error())
			}
		}
	}
	heartbeat.Timestamp = time.Now().UTC()
	err = repository.SaveHeartbeat(heartbeat)
	if err != nil {
		log.Printf("Error when trying to save heartbeat: %v", err.Error())
	}

	return heartbeat
}
