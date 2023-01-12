package apiv1

import (
	"github.com/fabmation-gmbh/briefkasten-go/internal/config"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

// AddApiV1 will add all v1 API handlers to the Mux router.
func AddApiV1(r fiber.Router) {
	// ========================================================
	// Unauthenticated routes

	r.Get("/oauth2/login/:provider", AuthLogin)
	r.Get("/oauth2/callback/:provider", OAuthCallback)

	// TODO: Implement logout

	// ========================================================
	// Authenticated routes

	r.Use(jwtware.New(jwtware.Config{
		SigningMethod: config.C.General.JWT.SigningMethod,
		SigningKey:    []byte(config.C.General.JWT.SigningKey),
	}))

	r.Get("/users/:id/tags", GetTags)
	r.Delete("/users/:id/tags/:tag_id", DeleteTag)
	r.Put("/users/:id/tags/:tag_id", UpdateTag)
	r.Post("/users/:id/tags", CreateTag)
}
