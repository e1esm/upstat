package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/chamanbravo/upstat/internal/dto"
	"github.com/chamanbravo/upstat/internal/models"
	"github.com/chamanbravo/upstat/svcerr"
)

func (r *Repository) CreateStatusPage(u *dto.CreateStatusPageIn) error {
	stmt, err := r.db.Prepare("INSERT INTO status_pages (name, slug) VALUES($1, $2)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(u.Name, u.Slug)
	if err != nil {
		return fmt.Errorf("failed to create status page: %w", err)
	}

	return nil
}

func (r *Repository) ListStatusPages() ([]*models.StatusPage, error) {
	stmt, err := r.db.Prepare("SELECT id, name, slug FROM status_pages")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("failed to get status pages: %w", err)
	}
	defer rows.Close()

	var statuspages []*models.StatusPage

	for rows.Next() {
		statuspage := new(models.StatusPage)
		err := rows.Scan(&statuspage.ID, &statuspage.Name, &statuspage.Slug)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows of status pages: %w", err)
		}
		statuspages = append(statuspages, statuspage)
	}

	return statuspages, nil
}

func (r *Repository) UpdateStatusPage(id int, statusPage *dto.CreateStatusPageIn) error {
	stmt, err := r.db.Prepare("UPDATE status_pages SET name = $1, slug = $2 WHERE id = $3")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result := stmt.QueryRow(statusPage.Name, statusPage.Slug, id)
	if result != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("%w with id: %d", svcerr.ErrStatusPageNotFound, id)
		}
		return err
	}

	return nil
}

func (r *Repository) DeleteStatusPageById(id int) error {
	stmt, err := r.db.Prepare("DELETE FROM status_pages WHERE id = $1")
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
		return fmt.Errorf("%w with id: %d", svcerr.ErrStatusPageNotFound, id)
	}

	return nil
}

func (r *Repository) FindStatusPageById(id int) (*models.StatusPage, error) {
	stmt, err := r.db.Prepare("SELECT * FROM status_pages WHERE id = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	statusPage := new(models.StatusPage)
	err = stmt.QueryRow(id).Scan(&statusPage.ID, &statusPage.Name, &statusPage.Slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, svcerr.ErrNoMonitorsFound
		}

		return nil, fmt.Errorf("%w: failed to find status page with id: %d", err, id)
	}

	return statusPage, nil
}

func (r *Repository) StatusPageMonitor(monitorId int, statusPages []string) error {
	stmt, err := r.db.Prepare("INSERT INTO status_pages_monitors(monitor_id, status_pages_id) VALUES($1, $2)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, statusPage := range statusPages {
		_, err = stmt.Exec(monitorId, statusPage)
		if err != nil {
			return fmt.Errorf("failed to create new status page for the monitor")
		}
	}

	return nil
}

func (r *Repository) UpdateStatusPageMonitorById(monitorId int, statusPages []string) error {
	stmt, err := r.db.Prepare("DELETE FROM status_pages_monitors WHERE monitor_id = $1")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(monitorId)
	if err != nil {
		log.Println("Error when trying to delete status page")
		log.Println(err)
		return fmt.Errorf("failed to delete status page: %w", err)
	}

	err = r.StatusPageMonitor(monitorId, statusPages)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) FindStatusPageByMonitorId(id int) ([]models.StatusPage, error) {
	stmt, err := r.db.Prepare(`
	SELECT
        n.id,
		n.name,
        n.slug
	FROM
		status_pages_monitors nm
	JOIN
		monitors m ON nm.monitor_id = m.id
	JOIN
		status_pages n ON nm.status_pages_id= n.id
	WHERE
		nm.monitor_id = $1;
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statusPages []models.StatusPage
	for rows.Next() {
		var statusPage models.StatusPage
		err = rows.Scan(&statusPage.ID, &statusPage.Name, &statusPage.Slug)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows for status page: %w", err)
		}

		statusPages = append(statusPages, statusPage)
	}

	return statusPages, nil
}

func (r *Repository) FindStatusPageBySlug(slug string) (*models.StatusPage, error) {
	stmt, err := r.db.Prepare("SELECT * FROM status_pages WHERE slug = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	statusPage := new(models.StatusPage)
	err = stmt.QueryRow(slug).Scan(&statusPage.ID, &statusPage.Name, &statusPage.Slug)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, svcerr.ErrStatusPageNotFound
		}

		return nil, err
	}

	return statusPage, nil
}

func (r *Repository) RetrieveStatusPageMonitors(slug string) ([]*models.Monitor, error) {
	stmt, err := r.db.Prepare(`
	SELECT DISTINCT
		spm.monitor_id,
		m.name	
	FROM
		status_pages_monitors spm
	JOIN
		monitors m ON spm.monitor_id = m.id
	LEFT JOIN
		heartbeats hb ON spm.monitor_id = hb.monitor_id
	WHERE
		spm.status_pages_id = (SELECT id FROM status_pages WHERE slug = $1);
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(slug)
	if err != nil {
		return nil, err
	}

	var monitors []*models.Monitor
	for rows.Next() {
		monitor := new(models.Monitor)
		err = rows.Scan(&monitor.ID, &monitor.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows for monitors: %w", err)
		}

		monitors = append(monitors, monitor)
	}

	return monitors, nil
}
