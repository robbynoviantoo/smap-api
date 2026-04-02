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
	AuthHandler             *handler.AuthHandler
	RoleHandler             *handler.RoleHandler
	AssetHandler            *handler.AssetHandler
	AssetPendingHandler     *handler.AssetPendingHandler
	ChatHandler             *handler.ChatHandler
	MaintenanceHandler      *handler.AssetMaintenanceHandler
	BorrowHandler           *handler.AssetBorrowHandler
	EventHandler            *handler.EventHandler
	PengadaanPendingHandler *handler.PengadaanAssetPendingHandler
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

	api := app.Group("/api/v1")

	// ── Auth (public) ─────────────────────────────────────────────────────────
	auth := api.Group("/auth")
	auth.Get("/", middleware.AuthRequired(), middleware.RequireRoles(config.DB, "superadmin"), d.AuthHandler.GetAllUsers)
	auth.Post("/login", d.AuthHandler.Login)

	// ── Role (public) ─────────────────────────────────────────────────────────
	role := api.Group("/role")
	role.Get("/", d.RoleHandler.GetAllRoles)
	role.Post("/", d.RoleHandler.CreateRole)
	role.Put("/", d.RoleHandler.UpdateRole)
	role.Delete("/:id", d.RoleHandler.DeleteRole)

	// ── Asset (public) ─────────────────────────────────────────────────────────
	asset := api.Group("/asset")
	asset.Get("/", d.AssetHandler.GetAssets)
	asset.Get("/export", d.AssetHandler.ExportAssets)
	asset.Post("/import", d.AssetHandler.ImportAssets)
	// Create goes to pending table
	asset.Post("/", middleware.AuthRequired(), d.AssetPendingHandler.CreatePending)
	asset.Put("/", d.AssetHandler.UpdateAsset)
	asset.Delete("/:id", d.AssetHandler.DeleteAsset)

	// ── Asset Pending (Review) ─────────────────────────────────────────────────
	assetPending := api.Group("/asset_pending")
	// Require Auth + Admin/Reviewer roles for these
	assetPending.Get("/", middleware.AuthRequired(), d.AssetPendingHandler.GetAllPendings)
	assetPending.Post("/:id/review", middleware.AuthRequired(), d.AssetPendingHandler.ReviewPending)

	// ── Chat (REST API) ────────────────────────────────────────────────────────
	chat := api.Group("/chat")
	chat.Get("/history", middleware.AuthRequired(), d.ChatHandler.GetHistory)

	// ── Asset Maintenance ──────────────────────────────────────────────────────
	mnt := api.Group("/asset_maintenance")
	mnt.Post("/:id/schedule", middleware.AuthRequired(), d.MaintenanceHandler.Schedule)
	mnt.Post("/:id/start", middleware.AuthRequired(), d.MaintenanceHandler.Start)
	mnt.Post("/:id/finish", middleware.AuthRequired(), d.MaintenanceHandler.Finish)

	// Admin/Reviewer routes for maintenance pendings
	mntPending := api.Group("/asset_maintenance_pending")
	mntPending.Get("/", middleware.AuthRequired(), d.MaintenanceHandler.GetAllPendings)
	mntPending.Post("/:id/review", middleware.AuthRequired(), d.MaintenanceHandler.ReviewPending)

	// ── Asset Borrow ───────────────────────────────────────────────────────────
	borrow := api.Group("/asset_borrow")
	borrow.Post("/:id/borrow", middleware.AuthRequired(), d.BorrowHandler.Borrow)
	borrow.Post("/:id/return", middleware.AuthRequired(), d.BorrowHandler.Return)

	// Admin/Reviewer routes for borrow pendings
	borrowPending := api.Group("/asset_borrow_pending")
	borrowPending.Get("/", middleware.AuthRequired(), d.BorrowHandler.GetAllPendings)
	borrowPending.Post("/:id/review", middleware.AuthRequired(), d.BorrowHandler.ReviewPending)

	// ── Events ───────────────────────────────────────────────────────────
	events := api.Group("/events")
	events.Get("/", d.EventHandler.GetAll)
	events.Post("/", middleware.AuthRequired(), d.EventHandler.Create)
	events.Put("/:id", middleware.AuthRequired(), d.EventHandler.Update)
	events.Delete("/:id", middleware.AuthRequired(), d.EventHandler.Delete)
	events.Post("/sync-holidays", middleware.AuthRequired(), d.EventHandler.SyncHolidays)

	// ── Pengadaan Asset Pending ──────────────────────────────────────────
	pengadaan := api.Group("/pengadaan_asset_pending")
	pengadaan.Post("/", middleware.AuthRequired(), d.PengadaanPendingHandler.Create)
	pengadaan.Get("/", middleware.AuthRequired(), d.PengadaanPendingHandler.GetAll)
	pengadaan.Get("/:id", middleware.AuthRequired(), d.PengadaanPendingHandler.GetByID)
	pengadaan.Put("/:id", middleware.AuthRequired(), d.PengadaanPendingHandler.Update)
	pengadaan.Delete("/:id", middleware.AuthRequired(), d.PengadaanPendingHandler.Delete)
	pengadaan.Post("/:id/approve", middleware.AuthRequired(), d.PengadaanPendingHandler.Approve)
	pengadaan.Post("/:id/reject", middleware.AuthRequired(), d.PengadaanPendingHandler.Reject)
}
