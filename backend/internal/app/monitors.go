package app

import (
	"time"

	"github.com/chamanbravo/upstat/internal/dto"
	"github.com/chamanbravo/upstat/internal/models"
)

func (a *App) CreateMonitor(u *dto.AddMonitorIn) (*models.Monitor, error) {
	return a.db.CreateMonitor(u)
}

func (a *App) DeleteMonitorById(id int) error {
	return a.db.DeleteMonitorById(id)
}

func (a *App) RetrieveUptime(id int, timestamp time.Time) (float64, error) {
	return a.db.RetrieveUptime(id, timestamp)
}

func (a *App) RetrieveAverageLatency(id int, timestamp time.Time) (float64, error) {
	return a.db.RetrieveAverageLatency(id, timestamp)
}

func (a *App) UpdateMonitorStatus(id int, status string) error {
	return a.db.UpdateMonitorStatus(id, status)
}

func (a *App) RetrieveMonitors() ([]*models.Monitor, error) {
	return a.db.RetrieveMonitors()
}

func (a *App) UpdateMonitorById(id int, monitor *dto.AddMonitorIn) error {
	return a.db.UpdateMonitorById(id, monitor)
}

func (a *App) FindMonitorById(id int) (*models.Monitor, error) {
	return a.db.FindMonitorById(id)
}

func (a *App) RetrieveHeartbeats(id, limit int) ([]*models.Heartbeat, error) {
	return a.db.RetrieveHeartbeats(id, limit)
}

func (a *App) StartMonitoringProcess(m *models.Monitor) {
	a.monitor.StartGoroutine(m)
}

func (a *App) StopMonitoringProcess(id int) {
	a.monitor.StopGoroutine(id)
}
