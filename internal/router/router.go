package router

import (
	"smap-api/internal/config"
	"smap-api/internal/handler"
	"smap-api/internal/middleware"
	"smap-api/internal/ws"

	"github.com/gofiber/fiber/v2"
	fiberws "github.com/gofiber/websocket/v2"
)

type Deps struct {
	AuthHandler    *handler.AuthHandler
	RoleHandler    *handler.RoleHandler
	AssetHandler   *handler.AssetHandler

}

func Setup(app *fiber.App, d Deps) {
	// ── Static files ──────────────────────────────────────────────────────────
	app.Static("/uploads", "./uploads")

	// ── Health check ──────────────────────────────────────────────────────────
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	// ── WebSocket ─────────────────────────────────────────────────────────────
	// Upgrade check middleware
	app.Use("/ws", func(c *fiber.Ctx) error {
		if fiberws.IsWebSocketUpgrade(c) {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	app.Get("/ws", fiberws.New(ws.ChatHandler))

	// ── Auth (public) ─────────────────────────────────────────────────────────
	auth := app.Group("/auth")
	auth.Get("/", middleware.AuthRequired(),middleware.RequireRoles(config.DB, "superadmin") , d.AuthHandler.GetAllUsers)
	auth.Post("/login", d.AuthHandler.Login)

	// ── Role (public) ─────────────────────────────────────────────────────────
	role := app.Group("/role")
	role.Get("/" , d.RoleHandler.GetAllRoles)
	role.Post("/" , d.RoleHandler.CreateRole)
	role.Put("/" , d.RoleHandler.UpdateRole)
	role.Delete("/:id" , d.RoleHandler.DeleteRole)

	// ── Asset (public) ─────────────────────────────────────────────────────────
	asset := app.Group("/asset")
	asset.Get("/" , d.AssetHandler.GetAssets)

}
