package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/chamanbravo/upstat/internal/dto"
	"github.com/chamanbravo/upstat/pkg"
	"github.com/gofiber/fiber/v2"
)

// @Tags Monitors
// @Accept json
// @Produce json
// @Param body body dto.AddMonitorIn true "Body"
// @Success 200 {object} dto.SuccessResponse
// @Success 400 {object} dto.ErrorResponse
// @Router /api/monitors [post]
func (h *Handler) CreateMonitor(c *fiber.Ctx) error {
	newMonitor := new(dto.AddMonitorIn)
	if err := c.BodyParser(newMonitor); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	errors := pkg.BodyValidator.Validate(newMonitor)
	if len(errors) > 0 {
		return c.Status(400).JSON(errors)
	}

	monitor, err := h.app.CreateMonitor(newMonitor)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	err = h.app.NotificationMonitor(monitor.ID, newMonitor.NotificationChannels)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	err = h.app.StatusPageMonitor(monitor.ID, newMonitor.StatusPages)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	h.app.StartMonitoringProcess(monitor)

	return c.Status(200).JSON(fiber.Map{
		"message": "success",
	})
}

// @Tags Monitors
// @Accept json
// @Produce json
// @Param id path string true "Monitor ID"
// @Success 200 {object} dto.MonitorInfoOut
// @Success 400 {object} dto.ErrorResponse
// @Router /api/monitors/{id} [get]
func (h *Handler) MonitorInfo(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "ID parameter is missing",
		})
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid ID parameter",
		})
	}

	monitor, err := h.app.FindMonitorById(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success",
		"monitor": monitor,
	})
}

// @Tags Monitors
// @Accept json
// @Produce json
// @Param id path string true "Monitor ID"
// @Param body body dto.AddMonitorIn true "Body"
// @Success 200 {object} dto.SuccessResponse
// @Success 400 {object} dto.ErrorResponse
// @Router /api/monitors/{id} [patch]
func (h *Handler) UpdateMonitor(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "ID parameter is missing",
		})
	}

	monitor := new(dto.AddMonitorIn)
	if err := c.BodyParser(monitor); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	errors := pkg.BodyValidator.Validate(monitor)
	if len(errors) > 0 {
		return c.Status(400).JSON(errors)
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid ID parameter",
		})
	}

	err = h.app.UpdateMonitorById(id, monitor)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	err = h.app.UpdateNotificationMonitorById(id, monitor.NotificationChannels)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	err = h.app.UpdateStatusPageMonitorById(id, monitor.StatusPages)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success",
	})

}

// @Tags Monitors
// @Accept json
// @Produce json
// @Param id path string true "Monitor ID"
// @Success 200 {object} dto.SuccessResponse
// @Success 400 {object} dto.ErrorResponse
// @Router /api/monitors/{id}/pause [patch]
func (h *Handler) PauseMonitor(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "ID parameter is missing",
		})
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid ID parameter",
		})
	}

	h.app.StopMonitoringProcess(id)
	err = h.app.UpdateMonitorStatus(id, "yellow")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success",
	})
}

// @Tags Monitors
// @Accept json
// @Produce json
// @Param id path string true "Monitor ID"
// @Success 200 {object} dto.SuccessResponse
// @Success 400 {object} dto.ErrorResponse
// @Router /api/monitors/{id}/resume [patch]
func (h *Handler) ResumeMonitor(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "ID parameter is missing",
		})
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid ID parameter",
		})
	}

	monitor, err := h.app.FindMonitorById(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	h.app.StartMonitoringProcess(monitor)
	err = h.app.UpdateMonitorStatus(id, "green")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success",
	})
}

// @Accept json
// @Produce json
// @Success 200 {object} dto.MonitorsListOut
// @Success 400 {object} dto.ErrorResponse
// @Router /api/monitors [get]
func (h *Handler) MonitorsList(c *fiber.Ctx) error {
	monitors, err := h.app.RetrieveMonitors()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	var monitorsList []fiber.Map
	for _, v := range monitors {
		heartbeat, err := h.app.RetrieveHeartbeats(v.ID, 10)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		monitorItem := fiber.Map{
			"id":        v.ID,
			"name":      v.Name,
			"url":       v.Url,
			"frequency": v.Frequency,
			"status":    v.Status,
			"heartbeat": heartbeat,
		}
		monitorsList = append(monitorsList, monitorItem)
	}

	return c.Status(200).JSON(fiber.Map{
		"message":  "success",
		"monitors": monitorsList,
	})
}

// @Tags Monitors
// @Accept json
// @Produce json
// @Param id path string true "Monitor ID"
// @Param startTime query time.Time true "Start Time" format(json)
// @Success 200 {object} dto.HeartbeatsOut
// @Success 400 {object} dto.ErrorResponse
// @Router /api/monitors/{id}/heartbeat [get]
func (h *Handler) RetrieveHeartbeat(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "ID parameter is missing",
		})
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid ID parameter",
		})
	}

	query := new(dto.RetrieveHeartbeatIn)
	if err := c.QueryParser(query); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	errors := pkg.BodyValidator.Validate(query)
	if len(errors) > 0 {
		return c.Status(400).JSON(errors)
	}

	heartbeat, err := h.app.RetrieveHeartbeatsByTime(id, query.StartTime)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{"message": "success", "heartbeat": heartbeat})
}

// @Tags Monitors
// @Accept json
// @Produce json
// @Param id path string true "Monitor ID"
// @Success 200 {object} dto.MonitorSummaryOut
// @Success 400 {object} dto.ErrorResponse
// @Router /api/monitors/{id}/summary [get]
func (h *Handler) MonitorSummary(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "ID parameter is missing",
		})
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid ID parameter",
		})
	}

	averageLatency, err := h.app.RetrieveAverageLatency(id, time.Now().Add(-time.Hour*24))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	dayUptime, err := h.app.RetrieveUptime(id, time.Now().Add(-time.Hour*24))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	monthUptime, err := h.app.RetrieveUptime(id, time.Now().Add(-time.Hour*30*24))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success",
		"summary": fiber.Map{
			"averageLatency": averageLatency,
			"dayUptime":      dayUptime,
			"monthUptime":    monthUptime,
		},
	})
}

// @Tags Monitors
// @Accept json
// @Produce json
// @Param id path string true "Monitor ID"
// @Success 200 {object} dto.SuccessResponse
// @Success 400 {object} dto.ErrorResponse
// @Router /api/monitors/{id} [delete]
func (h *Handler) DeleteMonitor(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "ID parameter is missing",
		})
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid ID parameter",
		})
	}

	h.app.StopMonitoringProcess(id)
	err = h.app.DeleteMonitorById(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{"message": "success"})
}

// @Tags Monitors
// @Accept json
// @Produce json
// @Param id path string true "Monitor ID"
// @Success 200 {object} dto.CertificateExpiryCountDown
// @Success 400 {object} dto.ErrorResponse
// @Router /api/monitors/{id}/cert-exp-countdown [get]
func (h *Handler) CertificateExpiryCountDown(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "ID parameter is missing",
		})
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid ID parameter",
		})
	}

	monitor, err := h.app.FindMonitorById(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	response, err := http.Get(monitor.Url)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	tlsInfo := response.TLS
	cert := tlsInfo.PeerCertificates[0]
	expirationDate := cert.NotAfter
	daysUnitlExp := int(expirationDate.Sub(time.Now()).Hours() / 24)

	return c.Status(200).JSON(fiber.Map{
		"message":             "success",
		"daysUntilExpiration": daysUnitlExp,
	})
}

// @Tags Monitors
// @Accept json
// @Produce json
// @Param id path string true "Monitor ID"
// @Success 200 {object} dto.NotificationListOut
// @Success 400 {object} dto.ErrorResponse
// @Router /api/monitors/{id}/notifications [get]
func (h *Handler) NotificationChannelListOfMonitor(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "ID parameter is missing",
		})
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid ID parameter",
		})
	}

	notification, err := h.app.FindNotificationChannelsByMonitorId(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message":       "success",
		"notifications": notification,
	})
}

// @Tags Monitors
// @Accept json
// @Produce json
// @Param id path string true "Monitor ID"
// @Success 200 {object} dto.ListStatusPagesOut
// @Success 400 {object} dto.ErrorResponse
// @Router /api/monitors/{id}/status-pages [get]
func (h *Handler) StatusPagesListOfMonitor(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "ID parameter is missing",
		})
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid ID parameter",
		})
	}

	statusPages, err := h.app.FindStatusPageByMonitorId(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message":     "success",
		"statusPages": statusPages,
	})
}
