package controllers

import (
	"github.com/chamanbravo/upstat/internal/dto"
	"github.com/chamanbravo/upstat/pkg"
	"github.com/gofiber/fiber/v2"
)

// @Accept json
// @Produce json
// @Success 200 {object} dto.NeedSetup
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/users/setup [get]
func (h *Handler) Setup(c *fiber.Ctx) error {
	usersCount, err := h.app.UsersCount()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"needSetup": usersCount <= 0,
	})
}

// @Accept json
// @Produce json
// @Param body body dto.UpdatePasswordIn true "Body"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/users/update-password [post]
func (h *Handler) UpdatePassword(c *fiber.Ctx) error {
	updatePasswordBody := new(dto.UpdatePasswordIn)
	if err := c.BodyParser(updatePasswordBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	errors := pkg.BodyValidator.Validate(updatePasswordBody)
	if len(errors) > 0 {
		return c.Status(400).JSON(errors)
	}

	username := c.Locals("username").(string)

	user, err := h.app.FindUserByUsername(username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Internal server error",
			"message": err.Error(),
		})
	}
	if user == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Not found",
			"message": "User does not exist",
		})
	}

	if err = pkg.CheckHash(user.Password, updatePasswordBody.CurrentPassword); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid current password",
		})
	}

	hashedNewPassword, err := pkg.HashAndSalt(updatePasswordBody.NewPassword)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	err = h.app.UpdatePassword(username, hashedNewPassword)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   fiber.ErrInternalServerError,
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "success",
	})
}

// @Accept json
// @Produce json
// @Param username path string true "Username"
// @Param body body dto.UpdateAccountIn true "Body"
// @Success 200 {object} dto.SuccessResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/users/me [patch]
func (h *Handler) UpdateAccount(c *fiber.Ctx) error {
	username := c.Locals("username").(string)
	if username == "" {
		return c.Status(400).JSON(fiber.Map{
			"error":   "Bad Request",
			"message": "Username parameter is missing",
		})
	}

	updateAccountBody := new(dto.UpdateAccountIn)
	if err := c.BodyParser(updateAccountBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	errors := pkg.BodyValidator.Validate(updateAccountBody)
	if len(errors) > 0 {
		return c.Status(400).JSON(errors)
	}

	err := h.app.UpdateAccount(username, updateAccountBody)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   fiber.ErrInternalServerError,
			"message": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Account updated successfully.",
	})
}
