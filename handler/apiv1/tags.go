package apiv1

import (
	"time"

	"github.com/fabmation-gmbh/briefkasten-go/handler/ftracer"
	"github.com/fabmation-gmbh/briefkasten-go/handler/rerr"
	"github.com/fabmation-gmbh/briefkasten-go/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// GetTags returns all tags the user has access to.
func GetTags(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := claims["user_id"].(uuid.UUID)
	ctx := ftracer.FromCtx(c)

	// TODO: Implement pagination
	var startTime time.Time

	tags, err := models.GetTagsByUserID(ctx, id, uuid.Nil, startTime, 100)
	if err != nil {
		return rerr.InternalServerError.With(err).WithLogMsg("unable to retrieve tags")
	}

	return c.JSON(tags)
}

type TagDeleteRequest struct {
	ID     uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}

// DeleteTag deletes a tag.
func DeleteTag(c *fiber.Ctx) error {
	// TODO: Use Params

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := claims["user_id"].(uuid.UUID)
	ctx := ftracer.FromCtx(c)

	var req TagDeleteRequest

	if err := c.BodyParser(&req); err != nil {
		return rerr.RequestMalformed.With(err).WithLogMsg("unable to parse body")
	}

	if id != req.UserID {
		return rerr.RequestMalformed.WithLogMsg("missmatching user IDs")
	}

	if err := models.DeleteTag(ctx, id, req.UserID); err != nil {
		return rerr.InternalServerError.With(err).WithLogMsg("unable to delete tag")
	}

	return c.JSON(fiber.Map{"message": "Deleted"})
}

// UpdateTag updates a tag.
func UpdateTag(c *fiber.Ctx) error {
	// TODO: Use Params (user)
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := claims["user_id"].(uuid.UUID)
	ctx := ftracer.FromCtx(c)

	tagID, err := parseUUIDParam(c, "tag_id")
	if err != nil {
		return rerr.RequestMalformed.With(err).WithLogMsg("invalid tag ID")
	}

	var req models.Tag

	if err := c.BodyParser(&req); err != nil {
		return rerr.RequestMalformed.With(err).WithLogMsg("unable to parse tag body")
	}

	req.ID = tagID
	req.UserID = id

	if id != req.UserID {
		return rerr.RequestMalformed.WithLogMsg("missmatching user IDs")
	}

	if err := req.Update(ctx); err != nil {
		return rerr.InternalServerError.With(err).WithLogMsg("unable to update tag")
	}

	return c.JSON(req)
}

// CreateTag creates a new tag.
func CreateTag(c *fiber.Ctx) error {
	// TODO: Use Params
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	id := claims["user_id"].(uuid.UUID)
	ctx := ftracer.FromCtx(c)

	var tag models.Tag

	if err := c.BodyParser(&tag); err != nil {
		return rerr.RequestMalformed.With(err).WithLogMsg("unable to parse tag body")
	}

	if id != tag.UserID {
		return rerr.RequestMalformed.WithLogMsg("missmatching user IDs")
	}

	tag.UserID = id
	tag.CreatedAt = time.Now()

	if err := tag.Create(ctx); err != nil {
		return rerr.InternalServerError.With(err).WithLogMsg("unable to create tag")
	}

	return c.JSON(tag)
}

func parseUUIDParam(c *fiber.Ctx, param string) (uuid.UUID, error) {
	str := c.Params(param)
	if str == "" {
		return uuid.Nil, rerr.RequestMalformed.WithLogMsg("missing ID parameter")
	}

	id, err := uuid.Parse(str)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
