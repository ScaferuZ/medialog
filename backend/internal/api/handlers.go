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
	authRequired := middleware.AuthMiddleware(jwtSecret)

	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Post("/refresh", authHandler.Refresh)
	auth.Get("/me", authRequired, authHandler.Me)

	mediaHandler.RegisterRoutes(api)
	usersHandler.RegisterRoutes(api)
	logsHandler.RegisterProtectedRoutes(api.Group("/logs", authRequired))
	socialHandler.RegisterPublicUserRoutes(api.Group("/users"))
	socialHandler.RegisterProtectedUserRoutes(api.Group("/users", authRequired))
	socialHandler.RegisterProtectedLogRoutes(api.Group("/logs", authRequired))
	socialHandler.RegisterProtectedReviewRoutes(api.Group("/reviews", authRequired))

	if tmdbHandler != nil {
		tmdbHandler.RegisterRoutes(api)
	}
}
