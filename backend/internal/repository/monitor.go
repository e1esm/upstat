package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/chamanbravo/upstat/internal/dto"
	"github.com/chamanbravo/upstat/internal/models"
	"github.com/chamanbravo/upstat/svcerr"
)

func (r *Repository) CreateMonitor(u *dto.AddMonitorIn) (*models.Monitor, error) {
	stmt, err := r.db.Prepare("INSERT INTO monitors(name, url, type, frequency, method, status) VALUES($1, $2, $3, $4, $5, $6) RETURNING id, frequency, url")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	monitor := new(models.Monitor)
	result := stmt.QueryRow(u.Name, u.URL, u.Type, u.Frequency, u.Method, "green").Scan(&monitor.ID, &monitor.Frequency, &monitor.Url)
	if result != nil {
		return nil, fmt.Errorf("failed to get inserted monitor: %w", result)
	}

	return monitor, nil
}

func (r *Repository) FindMonitorById(id int) (*models.Monitor, error) {
	stmt, err := r.db.Prepare("SELECT * FROM monitors WHERE id = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	monitor := new(models.Monitor)
	err = stmt.QueryRow(id).Scan(&monitor.ID, &monitor.Name, &monitor.Url, &monitor.Type, &monitor.Method, &monitor.Frequency, &monitor.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%w: failed to find monitor with id: %d", err, id)
		}

		return nil, err
	}

	return monitor, nil
}

func (r *Repository) UpdateMonitorById(id int, monitor *dto.AddMonitorIn) error {
	stmt, err := r.db.Prepare("UPDATE monitors SET name = $1, url = $2, type = $3, method = $4, frequency = $5  WHERE id = $6")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result := stmt.QueryRow(monitor.Name, monitor.URL, monitor.Type, monitor.Method, monitor.Frequency, id)
	if result.Err() != nil {
		return fmt.Errorf("failed to update monitor by id: %w", err)
	}

	return nil
}

func (r *Repository) RetrieveMonitors() ([]*models.Monitor, error) {
	stmt, err := r.db.Prepare("SELECT * FROM monitors")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("failed to get monitors: %w", err)
	}
	defer rows.Close()

	var monitors []*models.Monitor

	for rows.Next() {
		monitor := new(models.Monitor)
		err := rows.Scan(&monitor.ID, &monitor.Name, &monitor.Url, &monitor.Type, &monitor.Method, &monitor.Frequency, &monitor.Status)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row of monitors: %w", err)
		}
		monitors = append(monitors, monitor)
	}

	return monitors, nil
}

func (r *Repository) UpdateMonitorStatus(id int, status string) error {
	stmt, err := r.db.Prepare("UPDATE monitors SET status = $1 WHERE id = $2")
	if err != nil {
		log.Println("Error when trying to prepare statement")
		log.Println(err)
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(status, id)
	if err != nil {
		return fmt.Errorf("failed to update monitor status")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%w: id was not found", svcerr.ErrNoMonitorsFound, id)
	}

	return nil
}

func (r *Repository) RetrieveAverageLatency(id int, timestamp time.Time) (float64, error) {
	stmt, err := r.db.Prepare("SELECT AVG(latency) as average_latency FROM heartbeats WHERE monitor_id = $1 AND timestamp >= $2")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var averageLatency float64
	err = stmt.QueryRow(id, timestamp).Scan(&averageLatency)
	if err != nil {
		return 0, err
	}

	return averageLatency, nil
}

func (r *Repository) RetrieveUptime(id int, timestamp time.Time) (float64, error) {
	stmt, err := r.db.Prepare("SELECT (COUNT(CASE WHEN status = 'green' THEN 1 END) * 100.0) / COUNT(*) as green_percentage FROM heartbeats WHERE monitor_id = $1 AND timestamp >= $2")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var averageLatency float64
	err = stmt.QueryRow(id, timestamp).Scan(&averageLatency)
	if err != nil {
		return 0, err
	}

	return averageLatency, nil
}

func (r *Repository) DeleteMonitorById(id int) error {
	stmt, err := r.db.Prepare("DELETE FROM monitors WHERE id = $1")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%w: %d id not found", svcerr.ErrNoMonitorsFound, id)
	}

	return nil
}
