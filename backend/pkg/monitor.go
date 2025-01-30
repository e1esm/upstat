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
	"github.com/chamanbravo/upstat/pkg/alerts"
)

type DB interface {
	SaveIncident(incident *dto.SaveIncident) error
	LatestIncidentByMonitorId(id int) (*models.Incident, error)

	RetrieveMonitors() ([]*models.Monitor, error)
	FindMonitorById(id int) (*models.Monitor, error)
	UpdateMonitorStatus(id int, status string) error

	SaveHeartbeat(heartbeat *models.Heartbeat) error

	FindNotificationChannelsByMonitorId(id int) ([]models.Notification, error)
}

type Monitor struct {
	stopChannel   chan struct{}
	stopWaitGroup sync.WaitGroup
	goroutines    map[int]chan struct{}
	mutex         sync.Mutex
	db            DB
}

func New(db DB) *Monitor {
	return &Monitor{
		goroutines: make(map[int]chan struct{}, 0),
		mutex:      sync.Mutex{},
		db:         db,
	}
}

func (m *Monitor) StartGoroutine(monitor *models.Monitor) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	id := monitor.ID

	goroutineStopChannel := make(chan struct{})
	m.goroutines[id] = goroutineStopChannel

	m.stopWaitGroup.Add(1)

	go func() {
		defer func() {
			m.stopWaitGroup.Done()
			m.mutex.Lock()
			delete(m.goroutines, id)
			m.mutex.Unlock()
		}()

		for {
			select {
			case <-m.stopChannel:
				fmt.Printf("Goroutine with ID %d stopped\n", id)
				return
			case <-goroutineStopChannel:
				fmt.Printf("Goroutine with ID %d stopped by request\n", id)
				return
			default:
				monitor, err := m.db.FindMonitorById(id)
				if err != nil {
					log.Printf("Error retrieving updated monitor data: %v", err)
					continue
				}
				if monitor.Status != "yellow" {
					heartbeat := m.Ping(monitor)

					incidents, err := m.db.LatestIncidentByMonitorId(id)
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

						err = m.db.SaveIncident(newIncident)
						if err != nil {
							log.Printf("Error when trying to save incident: %v", err.Error())
						}

						notificationChannels, err := m.db.FindNotificationChannelsByMonitorId(id)
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

func (m *Monitor) StopGoroutine(id int) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if stopCh, exists := m.goroutines[id]; exists {
		close(stopCh)
	} else {
		fmt.Printf("Goroutine with ID %d not found\n", id)
	}
}

func (m *Monitor) StartGoroutineSetup() {
	monitors, err := m.db.RetrieveMonitors()
	if err != nil {
		log.Println("Error when trying to retrieve monitors")
		log.Println(err.Error())
	}

	for _, v := range monitors {
		if v.Status != "yellow" {
			m.StartGoroutine(v)
		}
	}
}

func (m *Monitor) Ping(monitor *models.Monitor) *models.Heartbeat {
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
			err := m.db.UpdateMonitorStatus(monitor.ID, "red")
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
			err := m.db.UpdateMonitorStatus(monitor.ID, "green")
			if err != nil {
				log.Printf("Error when trying to update monitor status: %v", err.Error())
			}
		}
	}
	heartbeat.Timestamp = time.Now().UTC()
	err = m.db.SaveHeartbeat(heartbeat)
	if err != nil {
		log.Printf("Error when trying to save heartbeat: %v", err.Error())
	}

	return heartbeat
}
