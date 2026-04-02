package handler

import (
	"smap-api/internal/model"
	"smap-api/internal/service"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	userSvc *service.UserService
}

func NewAuthHandler(userSvc *service.UserService) *AuthHandler {
	return &AuthHandler{userSvc: userSvc}
}

func (h *AuthHandler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.userSvc.GetAllUsers(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(users)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req model.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	resp, err := h.userSvc.Login(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(resp)
}

// GetMe godoc
// GET /api/v1/auth/me
func (h *AuthHandler) GetMe(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}
	user, err := h.userSvc.GetMe(c.Context(), userID.(uint))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(user)
}

// CompleteOnboarding godoc
// POST /api/v1/auth/onboarding
// Body: OnboardingRequest (first_name, last_name, no_handphone, password, photo_cropped)
func (h *AuthHandler) CompleteOnboarding(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	var req model.OnboardingRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	if err := h.userSvc.CompleteOnboarding(c.Context(), userID.(uint), &req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Onboarding berhasil. Selamat datang!"})
}