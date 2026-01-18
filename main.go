package main

import (
	"fmt"
	"log"
	"strings"
	"time"
	"user_service/db"
	"user_service/handlers"

	"github.com/thedanisaur/jfl_platform/auth"
	"github.com/thedanisaur/jfl_platform/config"
	"github.com/thedanisaur/jfl_platform/security"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func main() {
	log.Println("Starting JFL Flight Logging Service...")
	config, err := config.LoadConfig("./config.json")
	if err != nil {
		log.Printf("Error opening config, cannot continue: %s\n", err.Error())
		return
	}
	app := fiber.New()
	database, err := db.GetInstance()
	if err != nil {
		log.Print(err.Error())
	} else {
		defer database.Close()
	}

	// ==========================================
	// Load Auth Keys
	// ==========================================
	public_key, err := security.LoadPublicKey("secrets/public_signing_key.pem")
	if err != nil {
		log.Printf("Error opening public signing key, cannot continue: %s\n", err.Error())
		return
	}

	// ==========================================
	// Start Workers here
	// ==========================================
	// go db.DeleteExpiredUserSessions(time.Duration(config.App.LoginExpirationMs) * time.Millisecond)

	// ==========================================
	// Add CORS
	// ==========================================
	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(config.App.Cors.AllowOrigins, ","),
		AllowHeaders:     strings.Join(config.App.Cors.AllowHeaders, ","),
		AllowCredentials: config.App.Cors.AllowCredentials,
	}))

	// ==========================================
	// Add Rate Limiter
	// ==========================================
	var middleware limiter.LimiterHandler
	if config.App.Limiter.LimiterSlidingMiddleware {
		middleware = limiter.SlidingWindow{}
	} else {
		middleware = limiter.FixedWindow{}
	}
	app.Use(limiter.New(limiter.Config{
		Max:                    config.App.Limiter.Max,
		Expiration:             time.Duration(config.App.Limiter.Expiration),
		LimiterMiddleware:      middleware,
		SkipSuccessfulRequests: config.App.Limiter.SkipSuccessfulRequests,
	}))

	// ==========================================
	// Basic Authentication
	// ==========================================

	// intentionally left empty

	// ==========================================
	// JWT Authentication
	// ==========================================
	app.Get("/flight-logs/:user_id", auth.AuthorizationMiddleware(config, public_key), handlers.GetFlightlogs(config))
	app.Get("/flight-logs/:user_id/:id", auth.AuthorizationMiddleware(config, public_key), handlers.GetFlightlog(config))

	app.Post("/flight-logs/:id", auth.AuthorizationMiddleware(config, public_key), handlers.CreateFlightlog(config))

	app.Put("/flight-logs/:id", auth.AuthorizationMiddleware(config, public_key), handlers.UpdateFlightlog(config))

	// ==========================================
	// Start Service
	// ==========================================
	port := fmt.Sprintf(":%d", config.App.Host.Port)
	if config.App.Host.UseTLS {
		err = app.ListenTLS(port, config.App.Host.CertificatePath, config.App.Host.KeyPath)
	} else {
		log.Println("Warning - not using TLS")
		err = app.Listen(port)
	}
	if err != nil {
		log.Fatal(err.Error())
	}
}
