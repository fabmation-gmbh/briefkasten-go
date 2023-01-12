package apiv1

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/fabmation-gmbh/briefkasten-go/handler/ftracer"
	"github.com/fabmation-gmbh/briefkasten-go/handler/rerr"
	"github.com/fabmation-gmbh/briefkasten-go/internal/config"
	"github.com/fabmation-gmbh/briefkasten-go/internal/redis"
	"github.com/fabmation-gmbh/briefkasten-go/models"
	"github.com/fabmation-gmbh/briefkasten-go/pkg/helper"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/markbates/goth"
)

// AuthLogin is the authentication endpoint.
func AuthLogin(c *fiber.Ctx) error {
	provider := strings.ToLower(c.Params("provider"))

	prov, err := goth.GetProvider(provider)
	if err != nil {
		return rerr.RequestMalformed.With(err).WithLogMsg("unable to resolve requested provider")
	}

	state := helper.RandString(64)
	storeOAuthCookie(c, "oauth_state", state)

	ses, err := prov.BeginAuth(state)
	if err != nil {
		return rerr.InternalServerError.With(err).WithLogMsg("unable to begin authentication flow")
	}

	url, err := ses.GetAuthURL()
	if err != nil {
		return rerr.InternalServerError.With(err).WithLogMsg("unable to retrieve auth URL")
	}

	ctx := ftracer.FromCtx(c)

	userSes := redis.UserSession{
		Session:  ses,
		Provider: provider,
	}

	if err := redis.StoreUserSession(ctx, state, userSes); err != nil {
		return rerr.InternalServerError.With(err).WithLogMsg("unable to store session")
	}

	return c.Redirect(url, http.StatusFound)
}

// OAuthCallback is the oauth2 callback handler.
func OAuthCallback(c *fiber.Ctx) error {
	provider := strings.ToLower(c.Params("provider"))

	prov, err := goth.GetProvider(provider)
	if err != nil {
		return rerr.RequestMalformed.With(err).WithLogMsg("unable to resolve requested provider")
	}

	ctx := ftracer.FromCtx(c)

	cookieSession := c.Cookies("oauth_state")
	if cookieSession == "" {
		return rerr.RequestMalformed.With(err).WithLogMsg("missing state cookie")
	}

	ses, ok := redis.GetUserSession(ctx, cookieSession, prov)
	if !ok {
		return rerr.RequestMalformed.WithLogMsg("user session is not known")
	}

	if ses.Provider != provider {
		return rerr.RequestMalformed.WithLogMsg("invalid provider")
	}

	if cookieSession == "" || c.Query("state") != cookieSession {
		return rerr.RequestMalformed.WithLogMsg("state value missmatch")
	}

	oUser, err := prov.FetchUser(ses.Session)
	if err != nil {
		// TODO:sess.Authorize(provider, params) and retry FetchUser

		return rerr.RequestMalformed.WithLogMsg("unable to retrieve user")
	}

	acc := models.UserAccount{
		Email: oUser.Email,
		Name:  oUser.Name,
	}

	u, err := models.GetOrCreateUser(ctx, acc)
	if err != nil {
		return rerr.RequestMalformed.WithLogMsg("unable to retrieve or create user")
	}

	// TODO: Store the access token longer?

	claims := jwt.MapClaims{
		"user_id": u.ID.String(),
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString(String2bytes(config.C.General.JWT.SigningKey))
	if err != nil {
		return rerr.InternalServerError.WithLogMsg("unable to sign JWT token")
	}

	tokenData, err := json.Marshal(t)
	if err != nil {
		return rerr.InternalServerError.WithLogMsg("unable to marhsal JWT token")
	}

	c.Cookie(&fiber.Cookie{
		Name:     "briefkasten_jwt",
		Value:    byteSlice2String(tokenData),
		MaxAge:   int((time.Hour * 72).Seconds()),
		Secure:   config.C.General.SecureCookie,
		HTTPOnly: true,
	})

	return c.Redirect("/", http.StatusFound)
}
