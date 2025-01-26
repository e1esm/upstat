package app

import (
	"github.com/chamanbravo/upstat/internal/dto"
	"github.com/chamanbravo/upstat/internal/models"
)

func (a *App) CreateNotificationChannel(nc *dto.NotificationCreateIn) error {
	return a.db.CreateNotificationChannel(nc)
}

func (a *App) ListNotificationChannel() ([]dto.NotificationItem, error) {
	return a.db.ListNotificationChannel()
}

func (a *App) FindNotificationById(id int) (*models.Notification, error) {
	return a.db.FindNotificationById(id)
}

func (a *App) UpdateNotificationById(id int, nc *dto.NotificationCreateIn) error {
	return a.db.UpdateNotificationById(id, nc)
}

func (a *App) DeleteNotificationChannel(id int) error {
	return a.db.DeleteNotificationChannel(id)
}

func (a *App) UpdateNotificationMonitorById(monitorId int, notificationChannels []string) error {
	return a.db.UpdateNotificationMonitorById(monitorId, notificationChannels)
}

func (a *App) NotificationMonitor(monitorId int, notificationChannels []string) error {
	return a.db.NotificationMonitor(monitorId, notificationChannels)
}

func (a *App) FindNotificationChannelsByMonitorId(id int) ([]models.Notification, error) {
	return a.db.FindNotificationChannelsByMonitorId(id)
}
