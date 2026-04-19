package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jackc/pgx/v5"
	"github.com/medialogg/backend/internal/api"
	"github.com/medialogg/backend/internal/config"
	"github.com/medialogg/backend/internal/db"
	"github.com/medialogg/backend/internal/tmdb"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()

	var queries *db.Queries
	if cfg.DatabaseURL != "" {
		conn, err := pgx.Connect(ctx, cfg.DatabaseURL)
		if err != nil {
			log.Printf("database connection failed: %v", err)
			log.Println("continuing without database connection")
		} else {
			defer conn.Close(ctx)
			queries = db.New(conn)
			log.Println("database connected successfully")
		}
	} else {
		log.Println("DATABASE_URL not set, running without database")
	}

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3001,http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods:  "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowCredentials: true,
	}))

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	if queries != nil {
		authHandler := api.NewAuthHandler(queries, cfg.JWTSecret)
		mediaHandler := api.NewMediaHandler(queries)
		logsHandler := api.NewLogsHandler(queries)
		socialHandler := api.NewSocialHandler(queries)
		usersHandler := api.NewUsersHandler(queries)

		// Initialize TMDB client if API key is available
		var tmdbHandler *api.TMDBHandler
		if cfg.TMDBAPIKey != "" {
			tmdbClient := tmdb.NewClient(cfg.TMDBAPIKey)
			tmdbHandler = api.NewTMDBHandler(tmdbClient, queries)
			log.Println("TMDB client initialized")
		}

		apiGroup := app.Group("/api")
		api.SetupRoutes(apiGroup, authHandler, mediaHandler, logsHandler, socialHandler, usersHandler, tmdbHandler, cfg.JWTSecret)
	}

	go func() {
		log.Printf("server starting on :%s", cfg.ServerPort)
		if err := app.Listen(":" + cfg.ServerPort); err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(shutdownCtx); err != nil {
		log.Printf("server forced to shutdown: %v", err)
	}

	log.Println("server exited")
}
