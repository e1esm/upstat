package app

import (
	"time"

	"github.com/chamanbravo/upstat/internal/dto"
	"github.com/chamanbravo/upstat/internal/models"
)

type DB interface {
	FindUserByUsernameAndEmail(user *dto.UserSignUp) (*models.User, error)
	FindUserByUsername(username string) (*models.User, error)
	UsersCount() (int, error)
	UpdateAccount(username string, u *dto.UpdateAccountIn) error
	UpdatePassword(username string, newPassword string) error
	SaveUser(u *dto.UserSignUp) error

	CreateNotificationChannel(nc *dto.NotificationCreateIn) error
	ListNotificationChannel() ([]dto.NotificationItem, error)
	FindNotificationById(id int) (*models.Notification, error)
	UpdateNotificationById(id int, nc *dto.NotificationCreateIn) error
	DeleteNotificationChannel(id int) error
	FindNotificationChannelsByMonitorId(id int) ([]models.Notification, error)

	CreateStatusPage(u *dto.CreateStatusPageIn) error
	ListStatusPages() ([]*models.StatusPage, error)
	DeleteStatusPageById(id int) error
	UpdateStatusPage(id int, statusPage *dto.CreateStatusPageIn) error
	FindStatusPageById(id int) (*models.StatusPage, error)
	FindStatusPageBySlug(slug string) (*models.StatusPage, error)
	RetrieveStatusPageMonitors(slug string) ([]*models.Monitor, error)
	RetrieveHeartbeatsByTime(id int, startTime time.Time) ([]*models.Heartbeat, error)
	StatusPageMonitor(monitorId int, statusPages []string) error
	UpdateStatusPageMonitorById(monitorId int, statusPages []string) error
	FindStatusPageByMonitorId(id int) ([]models.StatusPage, error)

	CreateMonitor(u *dto.AddMonitorIn) (*models.Monitor, error)
	DeleteMonitorById(id int) error
	RetrieveUptime(id int, timestamp time.Time) (float64, error)
	RetrieveMonitors() ([]*models.Monitor, error)
	UpdateMonitorById(id int, monitor *dto.AddMonitorIn) error
	RetrieveAverageLatency(id int, timestamp time.Time) (float64, error)
	UpdateMonitorStatus(id int, status string) error
	FindMonitorById(id int) (*models.Monitor, error)
	RetrieveHeartbeats(id, limit int) ([]*models.Heartbeat, error)

	UpdateNotificationMonitorById(monitorId int, notificationChannels []string) error
	NotificationMonitor(monitorId int, notificationChannels []string) error
}

type Monitor interface {
	StartGoroutine(monitor *models.Monitor)
	StopGoroutine(id int)
	StartGoroutineSetup()
}

type App struct {
	db      DB
	monitor Monitor
}

func New(db DB, m Monitor) *App {
	return &App{
		db:      db,
		monitor: m,
	}
}
