package main

import (
	"log"

	"smap-api/internal/config"
	"smap-api/internal/handler"
	"smap-api/internal/middleware"
	"smap-api/internal/repository"
	"smap-api/internal/router"
	"smap-api/internal/service"
	"smap-api/internal/ws"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("[Main] No .env file found, reading from system environment")
	}

	// Load config
	config.Load()

	// Init DB and Redis
	config.InitDB()

	// ─── Repositories ───────────────────────────────────────────────────────────
	userRepo := repository.NewUserRepository(config.DB)
	roleRepo := repository.NewRoleRepository(config.DB)
	assetRepo := repository.NewAssetRepository(config.DB)
	assetPendingRepo := repository.NewAssetPendingRepository(config.DB)
	messageRepo := repository.NewMessageRepository(config.DB)
	settingRepo := repository.NewSettingRepository(config.DB)
	maintenanceRepo := repository.NewAssetMaintenanceRepository(config.DB)
	borrowRepo := repository.NewAssetBorrowRepository(config.DB)
	eventRepo := repository.NewEventRepository(config.DB)
	pengadaanPendingRepo := repository.NewPengadaanAssetPendingRepository(config.DB)

	// ─── Services ───────────────────────────────────────────────────────────────
	// hub := ws.GlobalHub
	userSvc := service.NewUserService(userRepo)
	roleSvc := service.NewRoleService(roleRepo)
	assetSvc := service.NewAssetService(assetRepo)
	assetPendingSvc := service.NewAssetPendingService(assetPendingRepo, assetRepo)
	messageSvc := service.NewMessageService(messageRepo)
	settingSvc := service.NewSettingService(settingRepo)
	maintenanceSvc := service.NewAssetMaintenanceService(maintenanceRepo, assetRepo, settingSvc)
	borrowSvc := service.NewAssetBorrowService(borrowRepo, assetRepo)
	eventSvc := service.NewEventService(eventRepo)
	pengadaanPendingSvc := service.NewPengadaanAssetPendingService(pengadaanPendingRepo)

	// Setup Hub Persistence
	ws.GlobalHub.MessageSvc = messageSvc

	// ─── Handlers ───────────────────────────────────────────────────────────────
	authH := handler.NewAuthHandler(userSvc)
	roleH := handler.NewRoleHandler(roleSvc)
	assetH := handler.NewAssetHandler(assetSvc)
	assetPendingH := handler.NewAssetPendingHandler(assetPendingSvc)
	chatH := handler.NewChatHandler(messageSvc)
	maintenanceH := handler.NewAssetMaintenanceHandler(maintenanceSvc)
	borrowH := handler.NewAssetBorrowHandler(borrowSvc)
	eventH := handler.NewEventHandler(eventSvc)
	pengadaanPendingH := handler.NewPengadaanAssetPendingHandler(pengadaanPendingSvc)
	dashboardH := handler.NewDashboardHandler(assetRepo)

	// ─── Fiber App ───────────────────────────────────────────────────────────────
	app := fiber.New(fiber.Config{
		AppName: "SMAP API v1.0",
	})

	// Global middleware
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))
	app.Use(middleware.RateLimiter())

	// Register all routes
	router.Setup(app, router.Deps{
		AuthHandler:             authH,
		RoleHandler:             roleH,
		AssetHandler:            assetH,
		AssetPendingHandler:     assetPendingH,
		ChatHandler:             chatH,
		MaintenanceHandler:      maintenanceH,
		BorrowHandler:           borrowH,
		EventHandler:            eventH,
		PengadaanPendingHandler: pengadaanPendingH,
		DashboardHandler:        dashboardH,
	})

	log.Printf("[Main] Server starting on :%s", config.App.AppPort)
	if err := app.Listen(":" + config.App.AppPort); err != nil {
		log.Fatalf("[Main] Server failed: %v", err)
	}
}
