package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/chamanbravo/upstat/internal/dto"
	"github.com/chamanbravo/upstat/internal/models"
	"github.com/chamanbravo/upstat/svcerr"
)

func (r *Repository) CreateNotificationChannel(nc *dto.NotificationCreateIn) error {
	stmt, err := r.db.Prepare("INSERT INTO notifications(name, provider, data) VALUES($1, $2, $3)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	dataJson, err := json.Marshal(nc.Data)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(nc.Name, nc.Provider, dataJson)
	if err != nil {
		return fmt.Errorf("failed to create new notification channel: %w", err)
	}

	return nil
}

func (r *Repository) ListNotificationChannel() ([]dto.NotificationItem, error) {
	stmt, err := r.db.Prepare("SELECT id, name, provider FROM notifications")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("failed to get notification channels: %w", err)
	}

	var notifications []dto.NotificationItem
	for rows.Next() {
		var notification dto.NotificationItem
		err = rows.Scan(&notification.ID, &notification.Name, &notification.Provider)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row of notification channels: %w", err)
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func (r *Repository) DeleteNotificationChannel(id int) error {
	stmt, err := r.db.Prepare("DELETE FROM notifications WHERE id = $1")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return svcerr.ErrNotificationsNotFound
		}

		return fmt.Errorf("failed to delete notification channel: %w", err)
	}

	return nil
}

func (r *Repository) UpdateNotificationById(id int, nc *dto.NotificationCreateIn) error {
	stmt, err := r.db.Prepare("UPDATE notifications SET name = $1, provider = $2, data = $3 WHERE id = $4")
	if err != nil {
		return err
	}
	defer stmt.Close()

	dataJson, err := json.Marshal(nc.Data)
	if err != nil {
		return err
	}

	result := stmt.QueryRow(nc.Name, nc.Provider, dataJson, id)
	if result != nil {
		if err == sql.ErrNoRows {
			return svcerr.ErrNotificationsNotFound
		}
		return err
	}

	return nil
}

func (r *Repository) FindNotificationById(id int) (*models.Notification, error) {
	stmt, err := r.db.Prepare("SELECT id, name, provider, data::text FROM notifications WHERE id = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	notification := new(models.Notification)
	var dataStr string

	err = stmt.QueryRow(id).Scan(&notification.ID, &notification.Name, &notification.Provider, &dataStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, svcerr.ErrNotificationsNotFound
		}

		return nil, fmt.Errorf("failed to get notification channel: %w", err)
	}

	if err := json.Unmarshal([]byte(dataStr), &notification.Data); err != nil {
		return nil, err
	}

	return notification, nil
}

func (r *Repository) NotificationMonitor(monitorId int, notificationChannels []string) error {
	stmt, err := r.db.Prepare("INSERT INTO notifications_monitors(monitor_id, notification_id) VALUES($1, $2)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, notificationChannel := range notificationChannels {
		_, err = stmt.Exec(monitorId, notificationChannel)
		if err != nil {
			return fmt.Errorf("failed to create new notification channel: %w", err)
		}
	}

	return nil
}

func (r *Repository) FindNotificationChannelsByMonitorId(id int) ([]models.Notification, error) {
	stmt, err := r.db.Prepare(`
	SELECT
        n.id,
		n.name AS notification_name,
		n.provider,
		n.data::text
	FROM
		notifications_monitors nm
	JOIN
		monitors m ON nm.monitor_id = m.id
	JOIN
		notifications n ON nm.notification_id = n.id
	WHERE
		nm.monitor_id = $1;
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var notification models.Notification
		var dataStr string
		err = rows.Scan(&notification.ID, &notification.Name, &notification.Provider, &dataStr)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row of notifications: %w", err)
		}

		if err := json.Unmarshal([]byte(dataStr), &notification.Data); err != nil {
			return nil, err
		}

		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func (r *Repository) UpdateNotificationMonitorById(monitorId int, notificationChannels []string) error {
	stmt, err := r.db.Prepare("DELETE FROM notifications_monitors WHERE monitor_id = $1")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(monitorId)
	if err != nil {
		return fmt.Errorf("failed to delete notification monitor: %w", err)
	}

	err = r.NotificationMonitor(monitorId, notificationChannels)
	if err != nil {
		return fmt.Errorf("failed to create new notification monitor: %w", err)
	}

	return nil
}
