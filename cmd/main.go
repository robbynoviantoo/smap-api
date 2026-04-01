package main

import (
	"log"

	"smap-api/internal/config"
	"smap-api/internal/handler"
	"smap-api/internal/middleware"
	"smap-api/internal/repository"
	"smap-api/internal/router"
	"smap-api/internal/service"
	// "booking-bioskop/internal/ws"

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

	// ─── Services ───────────────────────────────────────────────────────────────
	// hub := ws.GlobalHub
	userSvc := service.NewUserService(userRepo)
	roleSvc := service.NewRoleService(roleRepo)
	assetSvc := service.NewAssetService(assetRepo)

	// ─── Handlers ───────────────────────────────────────────────────────────────
	authH := handler.NewAuthHandler(userSvc)
	roleH := handler.NewRoleHandler(roleSvc)
	assetH := handler.NewAssetHandler(assetSvc)

	// ─── Fiber App ───────────────────────────────────────────────────────────────
	app := fiber.New(fiber.Config{
		AppName: "Booking Bioskop API v1.0",
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
		AuthHandler:    authH,
		RoleHandler:    roleH,
		AssetHandler:   assetH,
	})

	log.Printf("[Main] Server starting on :%s", config.App.AppPort)
	if err := app.Listen(":" + config.App.AppPort); err != nil {
		log.Fatalf("[Main] Server failed: %v", err)
	}
}
