package repository

import (
	"fmt"
	"time"

	"github.com/chamanbravo/upstat/internal/models"
)

func (r *Repository) RetrieveHeartbeats(id, limit int) ([]*models.Heartbeat, error) {
	stmt, err := r.db.Prepare("SELECT * FROM heartbeats WHERE monitor_id = $1 ORDER BY timestamp DESC limit $2")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(id, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query to get heartbeats: %w", err)
	}
	defer rows.Close()

	var heartbeats []*models.Heartbeat

	for rows.Next() {
		heartbeat := new(models.Heartbeat)
		err := rows.Scan(&heartbeat.ID, &heartbeat.MonitorId, &heartbeat.Timestamp, &heartbeat.StatusCode, &heartbeat.Status, &heartbeat.Latency, &heartbeat.Message)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows to get heartbeats: %w", err)
		}

		heartbeats = append(heartbeats, heartbeat)
	}

	return heartbeats, nil
}

func (r *Repository) SaveHeartbeat(heartbeat *models.Heartbeat) error {
	stmt, err := r.db.Prepare("INSERT INTO heartbeats(monitor_id, timestamp, status_code, status, latency, message) VALUES($1, $2, $3, $4, $5, $6)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(heartbeat.MonitorId, heartbeat.Timestamp, heartbeat.StatusCode, heartbeat.Status, heartbeat.Latency, heartbeat.Message)
	if err != nil {
		return fmt.Errorf("failed to save heartbeat: %w", err)
	}

	return nil
}

func (r *Repository) RetrieveHeartbeatsByTime(id int, startTime time.Time) ([]*models.Heartbeat, error) {
	stmt, err := r.db.Prepare("SELECT * FROM heartbeats WHERE monitor_id = $1 AND timestamp >= $2 ORDER BY timestamp ASC")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(id, startTime)
	if err != nil {
		return nil, fmt.Errorf("failed to get heartbeats: %w", err)
	}
	defer rows.Close()

	var heartbeats []*models.Heartbeat

	for rows.Next() {
		heartbeat := new(models.Heartbeat)
		err := rows.Scan(&heartbeat.ID, &heartbeat.MonitorId, &heartbeat.Timestamp, &heartbeat.StatusCode, &heartbeat.Status, &heartbeat.Latency, &heartbeat.Message)
		if err != nil {
			return nil, fmt.Errorf("faield to scan rows for heartbeats: %w", err)
		}
		heartbeats = append(heartbeats, heartbeat)
	}

	return heartbeats, nil
}
