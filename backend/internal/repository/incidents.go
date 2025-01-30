package repository

import (
	"database/sql"
	"fmt"

	"github.com/chamanbravo/upstat/internal/dto"
	"github.com/chamanbravo/upstat/internal/models"
	"github.com/chamanbravo/upstat/svcerr"
)

func (r *Repository) SaveIncident(incident *dto.SaveIncident) error {
	stmt, err := r.db.Prepare("INSERT INTO incidents(type, description, is_positive, monitor_id) VALUES($1, $2, $3, $4)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(incident.Type, incident.Description, incident.IsPositive, incident.MonitorId)
	if err != nil {
		return fmt.Errorf("failed to save incident: %w", err)
	}

	return nil
}

func (r *Repository) LatestIncidentByMonitorId(id int) (*models.Incident, error) {
	stmt, err := r.db.Prepare(`
    SELECT id, type, description, is_positive, monitor_id
    FROM incidents
    WHERE monitor_id = $1
    ORDER BY id DESC
    LIMIT 1;
    `)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(id)

	var incident models.Incident
	err = row.Scan(&incident.ID, &incident.Type, &incident.Description, &incident.IsPositive, &incident.MonitorId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, svcerr.ErrNoIncidentsFound
		}
		return nil, fmt.Errorf("faield to scan rows of incidents: %w", err)
	}

	return &incident, nil
}
