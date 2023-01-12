package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// NewRequestID returns the 'Request ID' middleware.
func NewRequestID(f *fiber.Ctx) error {
	u := uuid.New().String()

	// set X-Request-ID Header
	f.Set(fiber.HeaderXRequestID, u)

	// add UUID to context
	f.Context().SetUserValue(fiber.HeaderXRequestID, u)

	return f.Next()
}
