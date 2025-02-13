package controllers

import (
	"strconv"
	"time"

	"github.com/chamanbravo/upstat/internal/dto"
	"github.com/chamanbravo/upstat/internal/models"
	"github.com/chamanbravo/upstat/pkg"
	"github.com/gofiber/fiber/v2"
)

// @Tags StatusPages
// @Accept json
// @Produce json
// @Param body body dto.CreateStatusPageIn true "Body"
// @Success 200 {object} dto.SuccessResponse
// @Success 400 {object} dto.ErrorResponse
// @Router /api/status-pages [post]
func (h *Handler) CreateStatusPage(c *fiber.Ctx) error {
	statusPage := new(dto.CreateStatusPageIn)
	if err := c.BodyParser(statusPage); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	errors := pkg.BodyValidator.Validate(statusPage)
	if len(errors) > 0 {
		return c.Status(400).JSON(errors)
	}

	err := h.app.CreateStatusPage(statusPage)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Success",
	})
}

// @Tags StatusPages
// @Accept json
// @Produce json
// @Success 200 {object} dto.ListStatusPagesOut
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/status-pages [get]
func (h *Handler) ListStatusPages(c *fiber.Ctx) error {
	statusPages, err := h.app.ListStatusPages()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"errors": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message":     "Status pages list",
		"statusPages": statusPages,
	})
}

// @Accept json
// @Produce json
// @Param id path string true "Status Page ID"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/status-pages/{id} [delete]
func (h *Handler) DeleteStatusPage(c *fiber.Ctx) error {
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

	err = h.app.DeleteStatusPageById(id)
	if err != nil {
		return c.JSON(fiber.Map{
			"error":   err.Error(),
			"message": "Unable to delete a status page",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Status Page channels deleted",
	})
}

// @Tags StatusPages
// @Accept json
// @Produce json
// @Param id path string true "Status Page ID"
// @Param body body dto.CreateStatusPageIn true "Body"
// @Success 200 {object} dto.SuccessResponse
// @Success 400 {object} dto.ErrorResponse
// @Router /api/status-pages/{id} [patch]
func (h *Handler) UpdateStatusPage(c *fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "ID parameter is missing",
		})
	}

	statusPage := new(dto.CreateStatusPageIn)
	if err := c.BodyParser(statusPage); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	errors := pkg.BodyValidator.Validate(statusPage)
	if len(errors) > 0 {
		return c.Status(400).JSON(errors)
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid ID parameter",
		})
	}

	err = h.app.UpdateStatusPage(id, statusPage)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Successfully updated.",
	})

}

// @Tags StatusPages
// @Accept json
// @Produce json
// @Param id path string true "Status Page Id"
// @Success 200 {object} dto.StatusPageInfo
// @Success 400 {object} dto.ErrorResponse
// @Router /api/status-pages/{id} [get]
func (h *Handler) StatusPageInfo(c *fiber.Ctx) error {
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

	statusPage, err := h.app.FindStatusPageById(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message":    "success",
		"statusPage": statusPage,
	})
}

// @Tags Notifications
// @Accept json
// @Produce json
// @Param slug path string true "Status Page Slug"
// @Success 200 {object} dto.StatusPageSummary
// @Success 400 {object} dto.ErrorResponse
// @Router /api/status-pages/{slug}/summary [get]
func (h *Handler) StatusSummary(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return c.Status(400).JSON(fiber.Map{
			"message": "ID parameter is missing",
		})
	}

	statusPageInfo, err := h.app.FindStatusPageBySlug(slug)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	if statusPageInfo == nil {
		return c.Status(400).JSON(fiber.Map{
			"message":        "Status page not found",
			"statusPageInfo": nil,
		})
	}

	monitors, err := h.app.RetrieveStatusPageMonitors(slug)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	var monitorsList []fiber.Map
	startTime := time.Now().Add(time.Duration(-45) * time.Hour * 24)
	heartbeatMap := make(map[string]dto.HeartbeatSummary)
	for _, v := range monitors {
		heartbeat, err := h.app.RetrieveHeartbeatsByTime(v.ID, startTime)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		for _, v := range heartbeat {
			dateKey := v.Timestamp.Format("2006-01-02")
			value := heartbeatMap[dateKey]

			value.Total++
			value.Timestamp = dateKey
			if v.Status == "green" {
				value.Up++
			} else {
				value.Down++
			}

			heartbeatMap[dateKey] = value
		}

		allHeartbeats := []dto.HeartbeatSummary{}
		for _, v := range heartbeatMap {
			allHeartbeats = append(allHeartbeats, v)
		}

		recentHeartbeats := []models.Heartbeat{}
		for _, hb := range heartbeat {
			if hb.Timestamp.After(time.Now().Add(-12 * time.Hour)) {
				recentHeartbeats = append(recentHeartbeats, *hb)
			}
		}

		monitorItem := fiber.Map{
			"id":     v.ID,
			"name":   v.Name,
			"recent": recentHeartbeats,
			"all":    allHeartbeats,
		}
		monitorsList = append(monitorsList, monitorItem)
	}

	return c.JSON(fiber.Map{
		"message":        "Status pages list",
		"statusPageInfo": statusPageInfo,
		"monitors":       monitorsList,
	})
}
