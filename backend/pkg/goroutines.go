package pkg

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/chamanbravo/upstat/internal/dto"
	"github.com/chamanbravo/upstat/internal/models"
	"github.com/chamanbravo/upstat/internal/repository"
	"github.com/chamanbravo/upstat/pkg/alerts"
)

var (
	stopChannel   chan struct{}
	stopWaitGroup sync.WaitGroup
	goroutines    = make(map[int]chan struct{})
	mutex         sync.Mutex
)

func StartGoroutine(monitor *models.Monitor) {
	mutex.Lock()
	defer mutex.Unlock()

	id := monitor.ID

	goroutineStopChannel := make(chan struct{})
	goroutines[id] = goroutineStopChannel

	stopWaitGroup.Add(1)

	go func() {
		defer func() {
			stopWaitGroup.Done()
			mutex.Lock()
			delete(goroutines, id)
			mutex.Unlock()
		}()

		for {
			select {
			case <-stopChannel:
				fmt.Printf("Goroutine with ID %d stopped\n", id)
				return
			case <-goroutineStopChannel:
				fmt.Printf("Goroutine with ID %d stopped by request\n", id)
				return
			default:
				monitor, err := repository.FindMonitorById(id)
				if err != nil {
					log.Printf("Error retrieving updated monitor data: %v", err)
					continue
				}
				if monitor.Status != "yellow" {
					heartbeat := Ping(monitor)

					incidents, err := repository.LatestIncidentByMonitorId(id)
					if err != nil {
						log.Printf("Error when trying to retrieve incident: %v", err.Error())
					}

					if incidents == nil || (incidents.IsPositive != (heartbeat.Status == "green")) {
						var incidentType string
						if heartbeat.Status == "green" {
							incidentType = "UP"
						} else {
							incidentType = "DOWN"
						}
						newIncident := &dto.SaveIncident{
							Type: incidentType, Description: heartbeat.Message, IsPositive: heartbeat.Status == "green", MonitorId: id,
						}

						err = repository.SaveIncident(newIncident)
						if err != nil {
							log.Printf("Error when trying to save incident: %v", err.Error())
						}

						notificationChannels, err := repository.FindNotificationChannelsByMonitorId(id)
						if err != nil {
							log.Printf("Error when trying to retrieve notificationChannels: %v", err.Error())
						}

						discordMessage := alerts.DiscordAlertMessage(heartbeat, monitor)
						if err == nil {
							for _, v := range notificationChannels {
								jsonData, err := json.Marshal(discordMessage)
								if err == nil {
									_, err := http.Post(v.Data.WebhookUrl, "application/json", strings.NewReader(string(jsonData)))
									if err != nil {
										log.Printf("Error when trying to send heartbeat to webhook: %v", err.Error())
									}
								} else {
									log.Printf("Error when trying to convert heartbeat to JSON: %v", err.Error())
								}
							}
						} else {
							log.Printf("Error retrieving notification channels: %v", err)
						}
					}

				}
				time.Sleep(time.Duration(monitor.Frequency) * time.Second)
			}
		}
	}()
}

func StopGoroutine(id int) {
	mutex.Lock()
	defer mutex.Unlock()

	if stopCh, exists := goroutines[id]; exists {
		close(stopCh)
	} else {
		fmt.Printf("Goroutine with ID %d not found\n", id)
	}
}

func StartGoroutineSetup() {
	monitors, err := repository.RetrieveMonitors()
	if err != nil {
		log.Println("Error when trying to retrieve monitors")
		log.Println(err.Error())
	}

	for _, v := range monitors {
		if v.Status != "yellow" {
			StartGoroutine(v)
		}
	}
}
