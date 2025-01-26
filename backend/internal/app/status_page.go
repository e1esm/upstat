package app

import (
	"time"

	"github.com/chamanbravo/upstat/internal/dto"
	"github.com/chamanbravo/upstat/internal/models"
)

func (a *App) CreateStatusPage(u *dto.CreateStatusPageIn) error {
	return a.db.CreateStatusPage(u)
}

func (a *App) ListStatusPages() ([]*models.StatusPage, error) {
	return a.db.ListStatusPages()
}

func (a *App) DeleteStatusPageById(id int) error {
	return a.db.DeleteStatusPageById(id)
}

func (a *App) UpdateStatusPage(id int, statusPage *dto.CreateStatusPageIn) error {
	return a.db.UpdateStatusPage(id, statusPage)
}

func (a *App) FindStatusPageById(id int) (*models.StatusPage, error) {
	return a.db.FindStatusPageById(id)
}

func (a *App) FindStatusPageBySlug(slug string) (*models.StatusPage, error) {
	return a.db.FindStatusPageBySlug(slug)
}

func (a *App) RetrieveStatusPageMonitors(slug string) ([]*models.Monitor, error) {
	return a.db.RetrieveStatusPageMonitors(slug)
}

func (a *App) RetrieveHeartbeatsByTime(id int, startTime time.Time) ([]*models.Heartbeat, error) {
	return a.db.RetrieveHeartbeatsByTime(id, startTime)
}

func (a *App) StatusPageMonitor(monitorId int, statusPages []string) error {
	return a.db.StatusPageMonitor(monitorId, statusPages)
}

func (a *App) UpdateStatusPageMonitorById(monitorId int, statusPages []string) error {
	return a.db.UpdateStatusPageMonitorById(monitorId, statusPages)
}

func (a *App) FindStatusPageByMonitorId(id int) ([]models.StatusPage, error) {
	return a.db.FindStatusPageByMonitorId(id)
}
