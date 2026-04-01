package handler

import (
	"database/sql"
	"smap-api/internal/model"
	"smap-api/internal/service"

	"github.com/gofiber/fiber/v2"
)

type RoleHandler struct {
	roleSvc *service.RoleService
}

func NewRoleHandler(roleSvc *service.RoleService) *RoleHandler {
	return &RoleHandler{roleSvc: roleSvc}
}

func (h *RoleHandler) GetAllRoles(c *fiber.Ctx) error {
	roles, err := h.roleSvc.GetAllRoles(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch roles",
		})
	}
	return c.JSON(roles)
}

func (h *RoleHandler) CreateRole(c *fiber.Ctx) error {
	var req model.Role
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}
	err := h.roleSvc.CreateRole(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create role",
		})
	}
	return c.JSON(fiber.Map{
		"message": "role created successfully",
	})
}

func (h *RoleHandler) UpdateRole(c *fiber.Ctx) error {
	var req model.Role
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}
	err := h.roleSvc.UpdateRole(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to update role",
		})
	}
	return c.JSON(fiber.Map{
		"message": "role updated successfully",
	})
}


func (h *RoleHandler) DeleteRole(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "id tidak valid",
		})
	}

	err = h.roleSvc.DeleteRole(c.Context(), uint(id))
	if err != nil {

		// 🔥 custom message
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "id tidak ditemukan",
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "gagal menghapus role",
		})
	}

	return c.JSON(fiber.Map{
		"message": "role berhasil dihapus",
	})
}