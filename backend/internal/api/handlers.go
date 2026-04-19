package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/medialogg/backend/internal/middleware"
)

// SetupRoutes registers all API routes
func SetupRoutes(
	api fiber.Router,
	authHandler *AuthHandler,
	mediaHandler *MediaHandler,
	logsHandler *LogsHandler,
	socialHandler *SocialHandler,
	usersHandler *UsersHandler,
	tmdbHandler *TMDBHandler,
	jwtSecret string,
) {
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.Refresh)
	auth.Get("/me", middleware.AuthMiddleware(jwtSecret), authHandler.Me)

	mediaHandler.RegisterRoutes(api)
	logsHandler.RegisterRoutes(api)
	usersHandler.RegisterRoutes(api)
	socialHandler.RegisterRoutes(api.Group("", middleware.AuthMiddleware(jwtSecret)))

	if tmdbHandler != nil {
		tmdbHandler.RegisterRoutes(api)
	}
}
